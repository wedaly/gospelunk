package cmd

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/wedaly/gospelunk/db"
	pb "github.com/wedaly/gospelunk/db/protobuf"
	"github.com/wedaly/gospelunk/pkgmeta"
)

// Index adds definitions from Go packages to the search index.
func Index(dbPath string, pkgPatterns []string) error {
	pkgMetas, err := pkgmeta.Lookup(pkgPatterns)
	if err != nil {
		return errors.Wrapf(err, "pkgmeta.Lookup")
	}

	db, err := db.OpenReadWrite(dbPath)
	if err != nil {
		return errors.Wrapf(err, "db.OpenReadWrite")
	}
	defer db.Close()

	for _, pkg := range pkgMetas {
		pbPkg := protobufPackage(pkg)
		for _, filename := range pkg.GoFiles {
			path := filepath.Join(pkg.Dir, filename)
			defs, err := loadDefsFromGoFile(path)
			if err != nil {
				return err
			}
			pbPkg.GoFiles = append(pbPkg.GoFiles, protobufGoFile(filename, defs))
		}

		if err := db.WritePackage(pbPkg); err != nil {
			return err
		}
	}

	return nil
}

func loadDefsFromGoFile(path string) ([]pkgmeta.GoDef, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "os.Open")
	}
	defer f.Close()

	defs, err := pkgmeta.LoadGoDefs(f)
	if err != nil {
		return nil, errors.Wrapf(err, "pkgmeta.LoadGoDefs")
	}

	return defs, nil
}

func protobufPackage(pkg pkgmeta.Package) *pb.Package {
	return &pb.Package{
		Name:       pkg.Name,
		Dir:        pkg.Dir,
		ImportPath: pkg.ImportPath,
		GoFiles:    make([]*pb.GoFile, 0, len(pkg.GoFiles)),
	}
}

func protobufGoFile(filename string, defs []pkgmeta.GoDef) *pb.GoFile {
	pbDefs := make([]*pb.GoDef, len(defs))
	for i, def := range defs {
		pbDefs[i] = &pb.GoDef{
			Name:     def.Name,
			LineNum:  int64(def.LineNum),
			Exported: def.Exported,
		}
	}
	return &pb.GoFile{
		Filename: filename,
		Defs:     pbDefs,
	}
}
