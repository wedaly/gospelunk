package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"

	"github.com/wedaly/gospelunk/db"
	pb "github.com/wedaly/gospelunk/db/protobuf"
	"github.com/wedaly/gospelunk/pkgmeta"
)

// Index adds definitions from Go packages to the search index.
func Index(dbPath string, pkgPatterns []string, includeImports bool) error {
	fmt.Printf("Searching for packages matching %v\n", pkgPatterns)
	pkgMetas, err := pkgmeta.Lookup(pkgPatterns)
	if err != nil {
		return errors.Wrapf(err, "pkgmeta.Lookup")
	}

	if includeImports {
		fmt.Printf("Searching for imported packages\n")
		importedPkgMetas, err := pkgmeta.Lookup(uniqueSortedImportPkgNames(pkgMetas))
		if err != nil {
			return errors.Wrapf(err, "lookupImportedPackages")
		}
		pkgMetas = append(pkgMetas, importedPkgMetas...)
	}

	db, err := db.OpenReadWrite(dbPath)
	if err != nil {
		return errors.Wrapf(err, "db.OpenReadWrite")
	}
	defer db.Close()

	for _, pkg := range pkgMetas {
		fmt.Printf("Indexing %s\n", pkg.ImportPath)
		pbPkg := protobufPackage(pkg)
		for _, filename := range pkg.GoFiles {
			path := filepath.Join(pkg.Dir, filename)
			defs, err := loadDefsFromGoFile(path)
			if err != nil {
				fmt.Printf("WARN: could not index %s (%s)\n", path, err)
				continue
			}
			pbPkg.GoFiles = append(pbPkg.GoFiles, protobufGoFile(filename, defs))
		}

		if err := db.WritePackage(pbPkg); err != nil {
			return err
		}
	}

	return nil
}

func uniqueSortedImportPkgNames(pkgMetas []pkgmeta.Package) []string {
	nameSet := make(map[string]struct{}, 0)
	for _, pkg := range pkgMetas {
		for _, importPkgName := range pkg.Imports {
			nameSet[importPkgName] = struct{}{}
		}
	}

	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
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
	imports := make([]string, len(pkg.Imports))
	copy(imports, pkg.Imports)
	return &pb.Package{
		Name:       pkg.Name,
		Dir:        pkg.Dir,
		ImportPath: pkg.ImportPath,
		GoFiles:    make([]*pb.GoFile, 0, len(pkg.GoFiles)),
		Imports:    imports,
	}
}

func protobufGoFile(filename string, defs []pkgmeta.GoDef) *pb.GoFile {
	pbDefs := make([]*pb.GoDef, len(defs))
	for i, def := range defs {
		pbDefs[i] = &pb.GoDef{
			Name:     def.Name,
			Kind:     int32(def.Kind),
			LineNum:  int64(def.LineNum),
			Exported: def.Exported,
		}
	}
	return &pb.GoFile{
		Filename: filename,
		Defs:     pbDefs,
	}
}
