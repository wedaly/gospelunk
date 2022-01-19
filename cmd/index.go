package cmd

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"

	"github.com/wedaly/gospelunk/db"
	pb "github.com/wedaly/gospelunk/db/protobuf"
	"github.com/wedaly/gospelunk/log"
	"github.com/wedaly/gospelunk/pkgmeta"
)

// Index adds definitions from Go packages to the search index.
func Index(dbPath string, pkgPatterns []string, includeImports bool, transitive bool) error {
	log.Info("Searching for packages matching %v\n", pkgPatterns)
	pkgMetas, err := pkgmeta.Lookup(pkgPatterns)
	if err != nil {
		return errors.Wrapf(err, "pkgmeta.Lookup")
	}

	if includeImports {
		log.Info("Searching for imported packages\n")
		importedPkgMetas, err := lookupImportedPackages(pkgMetas, transitive)
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
		hash, err := pkgmeta.HashFileInfo(pkg)
		if err != nil {
			log.Warn("could not hash file info for pkg %s (%s)\n", pkg.Name, err)
		}

		cachedPkg, err := db.ReadPackage(pkg.Dir)
		if err != nil {
			return errors.Wrapf(err, "db.ReadPackage")
		}

		hashString := hash.String()
		if cachedPkg != nil && !hash.Empty() && cachedPkg.Hash == hashString {
			log.Info("Skipping pkg %s because it hasn't changed\n", cachedPkg.ImportPath)
			continue
		}

		log.Info("Indexing %s\n", pkg.ImportPath)
		pbPkg := protobufPackage(pkg)
		pbPkg.Hash = hashString
		for _, filename := range pkg.AllGoFiles() {
			path := filepath.Join(pkg.Dir, filename)
			defs, err := loadDefsFromGoFile(path)
			if err != nil {
				log.Warn("could not index %s (%s)\n", path, err)
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

func lookupImportedPackages(pkgMetas []pkgmeta.Package, transitive bool) ([]pkgmeta.Package, error) {
	importedSet := make(map[string]pkgmeta.Package, 0)
	for len(pkgMetas) > 0 {
		importedPkgMetas, err := pkgmeta.Lookup(uniqueSortedImportPkgNames(pkgMetas))
		if err != nil {
			return nil, errors.Wrapf(err, "pkgmeta.Lookup")
		}

		pkgMetas = nil
		for _, pkg := range importedPkgMetas {
			if _, ok := importedSet[pkg.Dir]; !ok {
				importedSet[pkg.Dir] = pkg
				pkgMetas = append(pkgMetas, pkg)
			}
		}

		if !transitive {
			break
		}
	}

	result := make([]pkgmeta.Package, 0, len(importedSet))
	for _, pkg := range importedSet {
		result = append(result, pkg)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result, nil
}

func uniqueSortedImportPkgNames(pkgMetas []pkgmeta.Package) []string {
	nameSet := make(map[string]struct{}, 0)
	for _, pkg := range pkgMetas {
		for _, importPkgName := range pkg.AllImports() {
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
	return &pb.Package{
		Name:       pkg.Name,
		Dir:        pkg.Dir,
		ImportPath: pkg.ImportPath,
		GoFiles:    make([]*pb.GoFile, 0, pkg.NumGoFiles()),
		Imports:    pkg.AllImports(),
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
