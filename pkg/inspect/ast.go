package inspect

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"golang.org/x/tools/go/packages"

	"github.com/wedaly/gospelunk/pkg/file"
)

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
		return nil, fmt.Errorf("filepath.Abs: %w", err)
	}

	for _, f := range pkg.Syntax {
		path := pkg.Fset.Position(f.Pos()).Filename
		if path == targetPath {
			return f, nil
		}
	}

	return nil, fmt.Errorf("Could not find ast.File for %q", targetPath)
}

// selectivelyParseFileFunc removes function bodies that do not contain the target line.
func selectivelyParseFileFunc(targetFilename string, targetLine int) func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
	return func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
		isTargetFile := filename == targetFilename
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
			if !isTargetFile || targetLine < start.Line || targetLine > end.Line {
				funcDecl.Body = nil
				return false
			} else {
				return true
			}
		})

		return astFile, nil
	}
}
