package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	// Step 2: load the Go package.
	absPath, err := filepath.Abs(pathArg)
	if err != nil {
		return err
	}

	// Step 3: find the AST identifier at the line/column
	// TODO

	// Step 4: lookup the definition.
	// TODO

	// Step 5: print the result.
	// TODO
}
