package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/pkg/errors"

	"github.com/wedaly/gospelunk/db"
	"github.com/wedaly/gospelunk/pkgmeta"
)

// TemplateInput provides data to the format template.
type TemplateInput struct {
	Path    string
	LineNum int64
	Name    string
}

// Find searches Go definitions in the search index.
func Find(dbPath string, query string, pkgPatterns []string, formatTpl string) error {
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

	pkgDirs, err := pkgmeta.ListDirs(pkgPatterns)
	if err != nil {
		return errors.Wrapf(err, "pkgmeta.ListDirs")
	}

	db, err := db.OpenReadOnly(dbPath)
	if err != nil {
		return errors.Wrapf(err, "db.OpenReadOnly")
	}
	defer db.Close()

	for _, pkgDir := range pkgDirs {
		pbPkg, err := db.ReadPackage(pkgDir)
		if err != nil {
			return errors.Wrapf(err, "db.ReadPackage")
		}

		if pbPkg == nil {
			// The package is not available in the search index, so skip it.
			continue
		}

		for _, goFile := range pbPkg.GoFiles {
			for _, def := range goFile.Defs {
				if queryRegexp.MatchString(def.Name) {
					err := tpl.Execute(os.Stdout, TemplateInput{
						Path:    filepath.Join(pbPkg.Dir, goFile.Filename),
						LineNum: def.LineNum,
						Name:    def.Name,
					})
					if err != nil {
						return errors.Wrapf(err, "template.Execute")
					}
					fmt.Printf("\n")
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
