package inspect

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"

	"github.com/wedaly/gospelunk/pkg/file"
)

func loadGoPackageForFileLoc(loc file.Loc) (*packages.Package, error) {
	absPath, err := filepath.Abs(loc.Path)
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
		ParseFile: selectivelyParseFileFunc(loc.Line),
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

	return nil, fmt.Errorf("Could not find Go package for path %q", loc.Path)
}

func selectivelyParseFileFunc(line int) func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
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
