package inspect

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

type FileLoc struct {
	Path   string
	Line   int
	Column int
}

type Definition struct {
	FileLoc
	Pkg  string
	Name string
}

type Result struct {
	Name string
	Type string
	Def  Definition
}

func Inspect(loc FileLoc) (*Result, error) {
	pkg, err := loadGoPackage(loc.Path, loc.Line)
	if err != nil {
		return nil, err
	}

	astFile, err := astFileForPath(pkg, loc.Path)
	if err != nil {
		return nil, err
	}

	astIdent := astIdentForLineAndCol(astFile, pkg.Fset, loc.Line, loc.Column)
	if astIdent == nil {
		return nil, nil
	}

	result := resultForAstIdent(pkg, astIdent)
	return result, nil
}

func loadGoPackage(path string, line int) (*packages.Package, error) {
	absPath, err := filepath.Abs(path)
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
		ParseFile: selectivelyParseFileFunc(line),
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

	return nil, fmt.Errorf("Could not find Go package for path %q", path)
}

type ParseFileFunc func(fset *token.FileSet, filename string, src []byte) (*ast.File, error)

func selectivelyParseFileFunc(line int) ParseFileFunc {
	return func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
		astFile, err := parser.ParseFile(fset, filename, src, 0)
		if err != nil {
			return nil, err
		}

		// Delete function bodies that don't contain the target line.
		// This reduces the amount of code we need to typecheck later.
		ast.Inspect(astFile, func(node ast.Node) bool {
			if node == nil {
				return false
			}

			funcDecl, ok := node.(*ast.FuncDecl)
			if !ok || funcDecl.Body == nil {
				return true
			}

			start := fset.Position(funcDecl.Body.Lbrace)
			end := fset.Position(funcDecl.Body.Rbrace)
			if line < start.Line || end.Line < line {
				funcDecl.Body = nil
				return false
			} else {
				return true
			}
		})

		return astFile, nil
	}
}

func astFileForPath(pkg *packages.Package, path string) (*ast.File, error) {
	targetPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrapf(err, "filepath.Abs")
	}

	for i, path := range pkg.GoFiles {
		if path == targetPath {
			return pkg.Syntax[i], nil
		}
	}

	return nil, fmt.Errorf("Could not find ast.File for %q", targetPath)
}

func astIdentForLineAndCol(astFile *ast.File, fset *token.FileSet, line int, col int) *ast.Ident {
	var foundIdent *ast.Ident
	ast.Inspect(astFile, func(node ast.Node) bool {
		if node == nil {
			return false
		}
		start := fset.Position(node.Pos())
		end := fset.Position(node.End())
		if line < start.Line || end.Line < line {
			return false
		}
		if ident, ok := node.(*ast.Ident); ok {
			if col < start.Column || end.Column < col {
				return false
			}
			foundIdent = ident
			return false
		}
		return true
	})
	return foundIdent
}

func resultForAstIdent(pkg *packages.Package, ident *ast.Ident) *Result {
	obj, ok := pkg.TypesInfo.Defs[ident]
	if !ok {
		obj = pkg.TypesInfo.Uses[ident]
	}

	if obj == nil || !obj.Pos().IsValid() {
		return nil
	}

	position := pkg.Fset.Position(obj.Pos())

	var pkgName string
	if obj.Pkg() != nil {
		pkgName = obj.Pkg().Name()
	}

	var typeName string
	if obj.Type() != types.Typ[types.Invalid] {
		typeName = obj.Type().String()
	}

	return &Result{
		Name: ident.Name,
		Type: typeName,
		Def: Definition{
			Pkg:  pkgName,
			Name: obj.Name(),
			FileLoc: FileLoc{
				Path:   position.Filename,
				Line:   position.Line,
				Column: position.Column,
			},
		},
	}
}
