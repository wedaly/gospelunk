package main

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strconv"

	"golang.org/x/tools/go/packages"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s FILE LINE COL\n", os.Args[0])
		os.Exit(1)
	}

	pathArg := os.Args[1]
	lineArg, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid line number %q\n", os.Args[2])
		os.Exit(1)
	}

	colArg, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid column number %q\n", os.Args[3])
		os.Exit(1)
	}

	err = lookupGoDef(pathArg, lineArg, colArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%w\n", err)
		os.Exit(1)
	}
}

func lookupGoDef(path string, line int, col int) error {
	// Step 1: load the Go package.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// TODO strip this down...
	loadMode := (packages.NeedName |
		packages.NeedFiles |
		packages.NeedSyntax |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedTypesInfo)

	cfg := &packages.Config{
		Mode: loadMode,
		Dir:  filepath.Dir(absPath),
	}

	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return err
	} else if len(pkgs) == 0 {
		return fmt.Errorf("No packages loaded")
	}

	pkg := pkgs[0] // TODO: explain this...

	// Step 2: find the AST for the target path.
	var astFile *ast.File
	for _, f := range pkg.Syntax {
		if pkg.Fset.Position(f.Pos()).Filename == absPath {
			astFile = f
			break
		}
	}
	if astFile == nil {
		return fmt.Errorf("Could not find AST file for %q", absPath)
	}

	// Step 3: find the AST identifier at the line/column
	var astIdent *ast.Ident
	ast.Inspect(astFile, func(node ast.Node) bool {
		if node == nil || astIdent != nil {
			return false
		}
		start, end := pkg.Fset.Position(node.Pos()), pkg.Fset.Position(node.End())
		if line < start.Line || line > end.Line || (line == start.Line && col < start.Column) || (line == end.Line && col > end.Column) {
			return false
		}
		if node, ok := node.(*ast.Ident); ok {
			astIdent = node
			return false
		}
		return true
	})

	// Step 4: lookup the definition.
	obj, ok := pkg.TypesInfo.Uses[astIdent]
	if !ok {
		obj = pkg.TypesInfo.Defs[astIdent]
	}
	if obj == nil {
		return fmt.Errorf("Could not find type object for ident %q at %s:%d:%d", astIdent.Name, path, line, col)
	} else if !obj.Pos().IsValid() {
		return fmt.Errorf("Invalid position for type object for %q at %s:%d:%d", astIdent.Name, path, line, col)
	}

	// Step 5: print the result.
	defPosition := pkg.Fset.Position(obj.Pos())
	fmt.Printf("%q is defined at %s:%d:%d\n", defPosition.Filename, defPosition.Line, defPosition.Column)

	return nil
}
