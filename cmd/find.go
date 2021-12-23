package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/pkg/errors"

	"github.com/wedaly/gospelunk/db"
	pb "github.com/wedaly/gospelunk/db/protobuf"
	"github.com/wedaly/gospelunk/pkgmeta"
)

// TemplateInput provides data to the format template.
type TemplateInput struct {
	Path    string
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

	return walkDefs(dbPath, pkgPatterns, includeImports, func(pkg *pb.Package, goFile *pb.GoFile, goDef *pb.GoDef) error {
		if !queryRegexp.MatchString(goDef.Name) {
			return nil
		}

		err := tpl.Execute(os.Stdout, TemplateInput{
			Path:    filepath.Join(pkg.Dir, goFile.Filename),
			LineNum: goDef.LineNum,
			Name:    goDef.Name,
		})
		if err != nil {
			return errors.Wrapf(err, "template.Execute")
		}

		fmt.Printf("\n")
		return nil
	})
}

func walkDefs(dbPath string, pkgPatterns []string, includeImports bool, f func(*pb.Package, *pb.GoFile, *pb.GoDef) error) error {
	db, err := db.OpenReadOnly(dbPath)
	if err != nil {
		return errors.Wrapf(err, "db.OpenReadOnly")
	}
	defer db.Close()

	pkgDirs, err := pkgmeta.ListDirs(pkgPatterns)
	if err != nil {
		return errors.Wrapf(err, "pkgmeta.ListDirs")
	}

	type stackItem struct {
		pkgDir        string
		isPkgImported bool
	}
	stack := make([]stackItem, 0, len(pkgDirs))
	for _, pkgDir := range pkgDirs {
		stack = append(stack, stackItem{pkgDir: pkgDir, isPkgImported: false})
	}

	visitedPkgDirs := make(map[string]struct{}, 0)
	var currentItem stackItem
	for len(stack) > 0 {
		currentItem, stack = stack[len(stack)-1], stack[0:len(stack)-1]
		if _, ok := visitedPkgDirs[currentItem.pkgDir]; ok {
			// Already visited this package.
			continue
		} else {
			// Mark this package as visited.
			visitedPkgDirs[currentItem.pkgDir] = struct{}{}
		}

		pkg, err := db.ReadPackage(currentItem.pkgDir)
		if err != nil {
			return errors.Wrapf(err, "db.ReadPackage")
		}

		if pkg == nil {
			// The package is not available in the search index, so skip it.
			continue
		}

		for _, goFile := range pkg.GoFiles {
			for _, def := range goFile.Defs {
				if currentItem.isPkgImported && !def.Exported {
					// Skip private definitions from imported packages.
					continue
				}
				if err := f(pkg, goFile, def); err != nil {
					return err
				}
			}
		}

		// Optionally include directly imported packages.
		if includeImports && !currentItem.isPkgImported && len(pkg.Imports) > 0 {
			importedPkgDirs, err := pkgmeta.ListDirs(pkg.Imports)
			if err != nil {
				fmt.Fprintf(os.Stderr, "WARN: could not find imported packages for %s\n", pkg.Name)
			} else {
				for _, pkgDir := range importedPkgDirs {
					stack = append(stack, stackItem{
						pkgDir:        pkgDir,
						isPkgImported: true,
					})
				}
			}
		}
	}

	return nil
}

func RelPathTplFunc() func(string) string {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARN: could not find current working directory (%s)\n", err)
		cwd = ""
	}

	return func(path string) string {
		relpath, err := filepath.Rel(cwd, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARN: could not determine relative path (%s)\n", err)
			relpath = path
		}
		return relpath
	}
}
