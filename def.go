package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	// Step 1: parse command-line arguments.
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

	// Step 2: load the Go package.
	// TODO

	// Step 3: find the AST identifier at the line/column
	// TODO

	// Step 4: lookup the definition.
	// TODO

	// Step 5: print the result.
	// TODO
}
