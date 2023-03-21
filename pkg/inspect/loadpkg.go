package inspect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/wedaly/gospelunk/pkg/file"
)

func isGoTestFile(path string) bool {
	return strings.HasSuffix(filepath.Base(path), "_test.go")
}

func loadGoPackageForFileLoc(loc file.Loc) (*packages.Package, error) {
	absPath, err := filepath.Abs(loc.Path)
	if err != nil {
		return nil, fmt.Errorf("filepath.Abs: %w", err)
	}

	cfg := &packages.Config{
		Mode: (packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedTypesInfo),
		Dir:       filepath.Dir(absPath),
		ParseFile: selectivelyParseFileFunc(absPath, loc.Line),
		Tests:     isGoTestFile(loc.Path),
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, fmt.Errorf("packages.Load: %w", err)
	}

	// If tests are included, pkgs will include both test and non-test packages.
	// The test packages have both test files *and* non-test Go files.
	// For non-test files, this loop will choose whichever package comes first
	// (when using the default build tool backend, that will be the non-test pkg).
	for _, pkg := range pkgs {
		for _, goFilePath := range pkg.GoFiles {
			if absPath == goFilePath {
				return pkg, nil
			}
		}
	}

	return nil, fmt.Errorf("Could not find Go package for path %q", loc.Path)
}

func loadGoPackagesMatchingPredicate(searchDir string, mode packages.LoadMode, includeTests bool, f func(skeletonPkg) bool) ([]*packages.Package, error) {
	// Find possible Go modules in the search directory (recursively).
	// This always includes the search directory itself, which may or may not be a Go module.
	possibleGoModDirs, err := findPossibleGoModDirsInSearchDir(searchDir)
	if err != nil {
		return nil, err
	}

	if includeTests {
		// Needed to deduplicate test/non-test pkgs.
		mode |= packages.NeedName
		mode |= packages.NeedFiles
	}

	var resultPkgs []*packages.Package
	for _, dir := range possibleGoModDirs {
		// Load minimal metadata for all packages in each possible Go module,
		// so we can quickly find packages that equal or import the target package.
		candidatePkgs, err := goListSkeletonPkgs(dir) // Returns an empty slice if dir isn't in a Go module.
		if err != nil {
			return nil, err
		}

		// Filter for pkgs that match the predicate.
		pkgPaths := make([]string, 0, len(candidatePkgs))
		for _, pkg := range candidatePkgs {
			if f(pkg) {
				pkgPaths = append(pkgPaths, pkg.ImportPath)
			}
		}

		if len(pkgPaths) == 0 {
			continue
		}

		// Parse and typecheck packages that either equal or import the target package.
		cfg := &packages.Config{
			Mode:  mode,
			Dir:   dir,
			Tests: includeTests,
		}

		pkgs, err := packages.Load(cfg, pkgPaths...)
		if err != nil {
			return nil, fmt.Errorf("packages.Load: %w", err)
		}

		// If tests are included, pkgs will include both test and non-test packages.
		// The test packages have both test files *and* non-test Go files.
		// Deduplicate these by choosing the test package over the non-test package.
		if includeTests {
			pkgs = deduplicateTestPkgs(pkgs)
		}

		resultPkgs = append(resultPkgs, pkgs...)
	}

	return resultPkgs, nil
}

func deduplicateTestPkgs(pkgs []*packages.Package) []*packages.Package {
	pkgSet := make(map[string]*packages.Package, len(pkgs))
	for _, pkg := range pkgs {
		if _, ok := pkgSet[pkg.PkgPath]; !ok {
			// Haven't seen this pkg yet, so choose it.
			pkgSet[pkg.PkgPath] = pkg
			continue
		}

		for _, goFilePath := range pkg.GoFiles {
			if isGoTestFile(goFilePath) {
				// Prioritize test pkg over non-test pkg.
				pkgSet[pkg.PkgPath] = pkg
				break
			}
		}
	}

	dedupedPkgs := make([]*packages.Package, 0, len(pkgSet))
	for _, pkg := range pkgSet {
		dedupedPkgs = append(dedupedPkgs, pkg)
	}
	return dedupedPkgs
}

func findPossibleGoModDirsInSearchDir(searchDir string) ([]string, error) {
	candidateSet := make(map[string]struct{}, 1)

	// Always include the search directory, even if it isn't in a Go module.
	candidateSet[filepath.Clean(searchDir)] = struct{}{}

	// Find subdirectories containing a "go.mod" file.
	err := filepath.WalkDir(searchDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && d.Name() == "go.mod" {
			candidateSet[filepath.Clean(filepath.Dir(path))] = struct{}{}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("filepath.WalkDir: %w", err)
	}

	// Convert relative paths to absolute paths in sorted order.
	result := make([]string, 0, len(candidateSet))
	for path := range candidateSet {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("filepath.Abs: %w", err)
		}
		result = append(result, absPath)
	}

	sort.Strings(result)
	return result, nil
}

// skeletonPkg contains minimal metadata for a package.
// Field names match the JSON output for the Package struct
// output by the `go list` cmd (see `go help list`).
type skeletonPkg struct {
	ImportPath string // Equivalent to the PkgPath field in packages.Package
	Imports    []string
}

// ImportsPkg checks whether the skeleton pkg imports a given package.
func (skel skeletonPkg) ImportsPkg(targetPkgPath string) bool {
	for _, importPkgPath := range skel.Imports {
		if importPkgPath == targetPkgPath {
			return true
		}
	}
	return false
}

// goListSkeletonPkgs returns skeleton pkgs for every package in a Go module.
// If goModDir isn't in a Go module, this returns an empty slice (no error).
func goListSkeletonPkgs(goModDir string) ([]skeletonPkg, error) {
	// We use the `go list` command directly instead of packages.Load
	// because we need the Dir field, which isn't exposed by packages.Load.
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("go", "list", "-json=ImportPath,Imports", "./...")
	cmd.Dir = goModDir
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Run(); err != nil {
		if strings.HasPrefix("go: go.mod file not found", stderrBuf.String()) ||
			strings.Contains(stderrBuf.String(), "does not contain main module or its selected dependencies") {
			// It's okay if we're not in a Go module.
			return nil, nil
		}
		return nil, fmt.Errorf("cmd.Run: %w", err)
	}

	// Stdout is a sequence of JSON-encoded package dictionaries.
	// This is NOT a JSON-encoded array, just one dict after another.
	var result []skeletonPkg
	for dec := json.NewDecoder(&stdoutBuf); dec.More(); {
		var skel skeletonPkg
		if err := dec.Decode(&skel); err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
		result = append(result, skel)
	}

	return result, nil
}
