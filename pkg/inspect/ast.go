package inspect

import (
	"fmt"
	"go/ast"
	"path/filepath"

	"github.com/pkg/errors"
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
