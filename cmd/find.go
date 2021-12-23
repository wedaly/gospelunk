package cmd

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"

	"github.com/wedaly/gospelunk/db"
	pb "github.com/wedaly/gospelunk/db/protobuf"
	"github.com/wedaly/gospelunk/pkgmeta"
)

// Find searches Go definitions in the search index.
func Find(dbPath string, query string, pkgPatterns []string) error {
	queryRegexp, err := regexp.Compile(query)
	if err != nil {
		return errors.Wrapf(err, "regexp.Compile")
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
					outputGoDef(pbPkg, goFile, def)
				}
			}
		}
	}

	return nil
}

func outputGoDef(pkg *pb.Package, goFile *pb.GoFile, def *pb.GoDef) {
	// TODO: make this configurable with a Go template argument.
	fmt.Printf("+%d -N %s/%s %s\n", def.LineNum, pkg.Dir, goFile.Filename, def.Name)
}
