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

	"github.com/wedaly/gospelunk/pkg/file"
)

type Result struct {
	Name      string
	Type      string
	Relations []Relation
}

func Inspect(loc file.Loc) (*Result, error) {
	pkg, err := loadGoPackage(loc.Path, loc.Line)
	if err != nil {
		return nil, err
	}

	astIdent, err := astNodeAtLoc[*ast.Ident](pkg, loc)
	if err != nil {
		return nil, err
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

// astNodeAtLoc locates the ast.Node of a particular type at a line and column in a file.
// If there are multiple such nodes, it returns the first one in depth-first traversal order.
// If there are no such nodes, it returns an error.
func astNodeAtLoc[T ast.Node](pkg *packages.Package, loc file.Loc) (T, error) {
	var (
		found     bool
		foundNode T
	)

	fset := pkg.Fset
	astFile, err := astFileForPath(pkg, loc.Path)
	if err != nil {
		return foundNode, err
	}

	ast.Inspect(astFile, func(node ast.Node) bool {
		if node == nil || found {
			return false
		}
		start := fset.Position(node.Pos())
		end := fset.Position(node.End())
		if loc.Line < start.Line || loc.Line > end.Line || (loc.Line == start.Line && loc.Column < start.Column) || (loc.Line == end.Line && loc.Column > end.Column) {
			return false
		}
		if node, ok := node.(T); ok {
			found = true
			foundNode = node
			return false
		}
		return true
	})

	if !found {
		return foundNode, fmt.Errorf("Could not find AST node of type %T at location %s", foundNode, loc)
	}

	return foundNode, nil
}

func astFileForPath(pkg *packages.Package, path string) (*ast.File, error) {
	targetPath, err := filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrapf(err, "filepath.Abs")
	}

	for _, f := range pkg.Syntax {
		path := pkg.Fset.Position(f.Pos()).Filename
		if path == targetPath {
			return f, nil
		}
	}

	return nil, fmt.Errorf("Could not find ast.File for %q", targetPath)
}

func resultForAstIdent(pkg *packages.Package, ident *ast.Ident) *Result {
	obj, ok := pkg.TypesInfo.Uses[ident]
	if !ok {
		obj = pkg.TypesInfo.Defs[ident]
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
		Relations: []Relation{
			{
				Kind: RelationKindDef,
				Pkg:  pkgName,
				Name: obj.Name(),
				Loc: file.Loc{
					Path:   position.Filename,
					Line:   position.Line,
					Column: position.Column,
				},
			},
		},
	}
}
