package list

import (
	"fmt"
	"go/ast"
	"sort"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	"github.com/wedaly/gospelunk/pkg/file"
)

type Options struct {
	IncludeStructFields     bool
	IncludeInterfaceMethods bool
	IncludePrivate          bool
	IncludeTests            bool
	IncludeImports          bool
}

type Result struct {
	Defs []Definition
}

type Package struct {
	Name string
	ID   string
}

type Definition struct {
	file.Loc
	Name string
	Pkg  Package
}

func List(patterns []string, opts Options) (Result, error) {
	var result Result

	pkgs, err := loadGoPackages(patterns, opts)
	if err != nil {
		return result, err
	}

	seenFiles := make(map[string]struct{})
	for _, pkg := range pkgs {
		goPaths := make(map[string]struct{}, len(pkg.GoFiles))
		for _, p := range pkg.GoFiles {
			goPaths[p] = struct{}{}
		}

		for _, astFile := range pkg.Syntax {
			path := pkg.Fset.Position(astFile.Pos()).Filename

			if _, ok := goPaths[path]; !ok {
				// Likely a compiled file from cgo. Ignore it.
				continue
			}

			if _, ok := seenFiles[path]; ok {
				// When opts.IncludeTests is true, the pkgs list will contain both the original pkg
				// as well as the pkg compiled for tests. Deduplicate the file paths to avoid duplicating
				// non-test definitions.
				continue
			} else {
				seenFiles[path] = struct{}{}
			}

			ast.Inspect(astFile, func(node ast.Node) bool {
				switch x := node.(type) {
				case *ast.ValueSpec:
					loadDefsFromValueSpec(pkg, opts, x, &result.Defs)
					return false

				case *ast.TypeSpec:
					loadDefsFromTypeSpec(pkg, opts, x, &result.Defs)
					return false

				case *ast.FuncDecl:
					loadDefsFromFuncDecl(pkg, opts, x, &result.Defs)
					return false

				default:
					return true
				}
			})
		}
	}

	sort.Slice(result.Defs, func(i, j int) bool {
		a, b := result.Defs[i], result.Defs[j]
		if a.Path != b.Path {
			return a.Path < b.Path
		} else if a.Line != b.Line {
			return a.Line < b.Line
		} else if a.Column != b.Column {
			return a.Column < b.Column
		} else {
			return a.Name < b.Name
		}
	})

	return result, nil
}

func loadGoPackages(patterns []string, opts Options) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes,
		Tests: opts.IncludeTests,
	}

	if opts.IncludeImports {
		cfg.Mode |= (packages.NeedImports | packages.NeedDeps)
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, errors.Wrapf(err, "packages.Load")
	}

	if opts.IncludeImports {
		pkgs = pkgsAndDirectImports(pkgs)
	}

	return pkgs, nil
}

func pkgsAndDirectImports(pkgs []*packages.Package) []*packages.Package {
	uniquePkgs := make(map[string]*packages.Package, len(pkgs))
	for _, pkg := range pkgs {
		uniquePkgs[pkg.ID] = pkg
		for _, importedPkg := range pkg.Imports {
			if _, ok := uniquePkgs[importedPkg.ID]; !ok {
				uniquePkgs[importedPkg.ID] = importedPkg
			}
		}
	}

	result := make([]*packages.Package, 0, len(uniquePkgs))
	for _, pkg := range uniquePkgs {
		result = append(result, pkg)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

func loadDefsFromValueSpec(pkg *packages.Package, opts Options, valueSpec *ast.ValueSpec, defs *[]Definition) {
	position := pkg.Fset.Position(valueSpec.Pos())
	for _, nameIdent := range valueSpec.Names {
		if nameIdent != nil && (opts.IncludePrivate || nameIdent.IsExported()) {
			valueName := nameIdent.Name
			*defs = append(*defs, Definition{
				Name: valueName,
				Pkg: Package{
					ID:   pkg.ID,
					Name: pkg.Name,
				},
				Loc: file.Loc{
					Path:   position.Filename,
					Line:   position.Line,
					Column: position.Column,
				},
			})
		}
	}
}

func loadDefsFromTypeSpec(pkg *packages.Package, opts Options, typeSpec *ast.TypeSpec, defs *[]Definition) {
	if typeSpec.Name == nil || (!opts.IncludePrivate && !typeSpec.Name.IsExported()) {
		return
	}

	position := pkg.Fset.Position(typeSpec.Pos())
	typeName := typeSpec.Name.Name

	switch x := typeSpec.Type.(type) {
	case *ast.StructType:
		loadDefsFromStructType(pkg, opts, typeName, x, defs)
	case *ast.InterfaceType:
		loadDefsFromInterfaceType(pkg, opts, typeName, x, defs)
	}

	*defs = append(*defs, Definition{
		Name: typeName,
		Pkg: Package{
			ID:   pkg.ID,
			Name: pkg.Name,
		},

		Loc: file.Loc{
			Path:   position.Filename,
			Line:   position.Line,
			Column: position.Column,
		},
	})
}

func loadDefsFromStructType(pkg *packages.Package, opts Options, typeName string, structType *ast.StructType, defs *[]Definition) {
	if !opts.IncludeStructFields {
		return
	}

	for _, field := range structType.Fields.List {
		position := pkg.Fset.Position(field.Pos())
		for _, nameIdent := range field.Names {
			if nameIdent != nil && (opts.IncludePrivate || nameIdent.IsExported()) {
				fieldName := nameIdent.Name
				*defs = append(*defs, Definition{
					Name: fmt.Sprintf("%s.%s", typeName, fieldName),
					Pkg: Package{
						ID:   pkg.ID,
						Name: pkg.Name,
					},
					Loc: file.Loc{
						Path:   position.Filename,
						Line:   position.Line,
						Column: position.Column,
					},
				})
			}
		}
	}
}

func loadDefsFromInterfaceType(pkg *packages.Package, opts Options, typeName string, interfaceType *ast.InterfaceType, defs *[]Definition) {
	if !opts.IncludeInterfaceMethods {
		return
	}

	for _, method := range interfaceType.Methods.List {
		position := pkg.Fset.Position(method.Pos())
		for _, nameIdent := range method.Names {
			if nameIdent != nil && (opts.IncludePrivate || nameIdent.IsExported()) {
				methodName := nameIdent.Name
				*defs = append(*defs, Definition{
					Name: fmt.Sprintf("%s.%s", typeName, methodName),
					Pkg: Package{
						ID:   pkg.ID,
						Name: pkg.Name,
					},
					Loc: file.Loc{
						Path:   position.Filename,
						Line:   position.Line,
						Column: position.Column,
					},
				})
			}
		}
	}
}

func loadDefsFromFuncDecl(pkg *packages.Package, opts Options, funcDecl *ast.FuncDecl, defs *[]Definition) {
	if funcDecl.Name == nil || (!opts.IncludePrivate && !funcDecl.Name.IsExported()) {
		return
	}
	position := pkg.Fset.Position(funcDecl.Pos())
	name := funcDecl.Name.Name
	if funcDecl.Recv != nil {
		name = fmt.Sprintf("%s.%s", findFuncRecvName(funcDecl), name)
	}
	*defs = append(*defs, Definition{
		Name: name,
		Pkg: Package{
			ID:   pkg.ID,
			Name: pkg.Name,
		},
		Loc: file.Loc{
			Path:   position.Filename,
			Line:   position.Line,
			Column: position.Column,
		},
	})
}

func findFuncRecvName(funcDecl *ast.FuncDecl) string {
	var typeName string
	for _, field := range funcDecl.Recv.List {
		ast.Inspect(field.Type, func(node ast.Node) bool {
			if ident, ok := node.(*ast.Ident); ok {
				typeName = ident.Name
				return false
			}
			return true
		})
	}
	return typeName
}
