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

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	"github.com/wedaly/gospelunk/pkg/file"
)

func loadGoPackageForFileLoc(loc file.Loc) (*packages.Package, error) {
	absPath, err := filepath.Abs(loc.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "filepath.Abs")
	}

	cfg := &packages.Config{
		Mode: (packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedTypesInfo),
		Dir:       filepath.Dir(absPath),
		ParseFile: selectivelyParseFileFunc(loc.Line),
		Tests:     true,
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, errors.Wrapf(err, "packages.Load")
	}

	for _, pkg := range pkgs {
		for _, goFilePath := range pkg.GoFiles {
			if absPath == goFilePath {
				return pkg, nil
			}
		}
	}

	return nil, fmt.Errorf("Could not find Go package for path %q", loc.Path)
}

func loadGoPackagesMatchingPredicate(searchDir string, mode packages.LoadMode, f func(skeletonPkg) bool) ([]*packages.Package, error) {
	// Find possible Go modules in the search directory (recursively).
	// This always includes the search directory itself, which may or may not be a Go module.
	possibleGoModDirs, err := findPossibleGoModDirsInSearchDir(searchDir)
	if err != nil {
		return nil, err
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
			Mode: mode,
			Dir:  dir,
		}

		pkgs, err := packages.Load(cfg, pkgPaths...)
		if err != nil {
			return nil, errors.Wrapf(err, "packages.Load")
		}

		resultPkgs = append(resultPkgs, pkgs...)
	}

	return resultPkgs, nil
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
		return nil, errors.Wrapf(err, "filepath.WalkDir")
	}

	// Convert relative paths to absolute paths in sorted order.
	result := make([]string, 0, len(candidateSet))
	for path := range candidateSet {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, errors.Wrapf(err, "filepath.Abs")
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
	ImportPath string // Equivalent to the ID field in packages.Package
	Imports    []string
}

// ImportsPkg checks whether the skeleton pkg imports a given package.
func (skel skeletonPkg) ImportsPkg(targetPkgId string) bool {
	for _, importPkgId := range skel.Imports {
		if importPkgId == targetPkgId {
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
		return nil, errors.Wrapf(err, "cmd.Run")
	}

	// Stdout is a sequence of JSON-encoded package dictionaries.
	// This is NOT a JSON-encoded array, just one dict after another.
	var result []skeletonPkg
	for dec := json.NewDecoder(&stdoutBuf); dec.More(); {
		var skel skeletonPkg
		if err := dec.Decode(&skel); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal")
		}
		result = append(result, skel)
	}

	return result, nil
}
