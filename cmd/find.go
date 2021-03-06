package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/wedaly/gospelunk/db"
	pb "github.com/wedaly/gospelunk/db/protobuf"
	"github.com/wedaly/gospelunk/log"
	"github.com/wedaly/gospelunk/pkgmeta"
)

// TemplateInput provides data to the format template.
type TemplateInput struct {
	Path    string
	Kind    string
	LineNum int64
	Name    string
}

// Find searches Go definitions in the search index.
func Find(dbPath string, query string, pkgPatterns []string, includeImports bool, formatTpl string) error {
	queryRegexp, err := regexp.Compile(query)
	if err != nil {
		return errors.Wrapf(err, "regexp.Compile")
	}

	tpl, err := template.New("gospelunk").
		Funcs(template.FuncMap{
			"RelPath":  RelPathTplFunc(),
			"BasePath": filepath.Base,
		}).
		Parse(formatTpl)
	if err != nil {
		return errors.Wrapf(err, "template.Parse")
	}

	matchedPkgs, err := getMatchingPackagesFromDb(dbPath, pkgPatterns, includeImports)
	if err != nil {
		return err
	}

	for _, mp := range matchedPkgs {
		for _, goFile := range mp.pkg.GoFiles {
			for _, goDef := range goFile.Defs {
				if mp.fromImport && (!goDef.Exported || isTestFile(goFile)) {
					// Skip private definitions and tests from imported packages.
					continue
				}

				if !queryRegexp.MatchString(goDef.Name) {
					continue
				}

				err := tpl.Execute(os.Stdout, TemplateInput{
					Path:    filepath.Join(mp.pkg.Dir, goFile.Filename),
					Kind:    pkgmeta.GoDefKind(goDef.Kind).String(),
					LineNum: goDef.LineNum,
					Name:    goDef.Name,
				})
				if err != nil {
					return errors.Wrapf(err, "template.Execute")
				}

				fmt.Printf("\n")
			}
		}
	}

	return nil
}

func isTestFile(goFile *pb.GoFile) bool {
	return strings.HasSuffix(goFile.Filename, "_test.go")
}

// matchedPackage represents a package from the database.
type matchedPackage struct {
	pkg        *pb.Package
	fromImport bool
}

func getMatchingPackagesFromDb(dbPath string, pkgPatterns []string, includeImports bool) ([]matchedPackage, error) {
	db, err := db.OpenReadOnly(dbPath)
	if err != nil {
		return nil, errors.Wrapf(err, "db.OpenReadOnly")
	}
	defer db.Close()

	pkgDirs, err := pkgmeta.ListDirs(pkgPatterns)
	if err != nil {
		return nil, errors.Wrapf(err, "pkgmeta.ListDirs")
	}

	packages := make([]matchedPackage, 0, len(pkgDirs))
	seenPkgDirSet := make(map[string]struct{}, 0)
	importSet := make(map[string]struct{}, 0)
	err = iterPackages(db, pkgDirs, func(pkg *pb.Package) {
		packages = append(packages, matchedPackage{
			pkg:        pkg,
			fromImport: false,
		})
		seenPkgDirSet[pkg.Dir] = struct{}{}
		if includeImports {
			for _, importPkgName := range pkg.Imports {
				importSet[importPkgName] = struct{}{}
			}
		}
	})
	if err != nil {
		return nil, err
	}

	if includeImports {
		importPkgNames := make([]string, 0, len(importSet))
		for pkgName := range importSet {
			importPkgNames = append(importPkgNames, pkgName)
		}
		sort.Strings(importPkgNames)

		if len(importPkgNames) > 0 {
			pkgDirs, err := pkgmeta.ListDirs(importPkgNames)
			if err != nil {
				return nil, errors.Wrapf(err, "pkgmeta.ListDirs")
			}

			err = iterPackages(db, pkgDirs, func(pkg *pb.Package) {
				if _, ok := seenPkgDirSet[pkg.Dir]; ok {
					return
				}
				packages = append(packages, matchedPackage{
					pkg:        pkg,
					fromImport: true,
				})
				seenPkgDirSet[pkg.Dir] = struct{}{}
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return packages, nil
}

func iterPackages(db *db.DB, pkgDirs []string, f func(*pb.Package)) error {
	for _, pkgDir := range pkgDirs {
		pkg, err := db.ReadPackage(pkgDir)
		if err != nil {
			return errors.Wrapf(err, "db.ReadPackage")
		}
		if pkg != nil {
			f(pkg)
		}
	}
	return nil
}

func RelPathTplFunc() func(string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Warn("could not find current working directory (%s)\n", err)
		cwd = ""
	}

	return func(path string) string {
		relpath, err := filepath.Rel(cwd, path)
		if err != nil {
			log.Warn("could not determine relative path (%s)\n", err)
			relpath = path
		}
		return relpath
	}
}
