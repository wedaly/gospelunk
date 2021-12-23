package main

import (
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/pkg/errors"

	"github.com/wedaly/gospelunk/cmd"
)

const UsageMsg = `Usage: gospelunk <command> [args]

Available commands:
	index    create or update the search index
	find     search for definitions in the index`

const IndexUsageMsg = `Usage: gospelunk index [flags] <packages>`

const FindUsageMsg = `Usage: gospelunk find [flags] <query> [packages]`

func main() {
	if len(os.Args) < 2 {
		handleInvalidCmd("")
	}

	switch cmd, args := os.Args[1], os.Args[2:]; cmd {
	case "index":
		handleIndexCmd(args)
	case "find":
		handleFindCmd(args)
	case "-h":
		printUsage("Search definitions in a Go project and its dependencies", UsageMsg)
	default:
		handleInvalidCmd(cmd)
	}
}

func handleIndexCmd(args []string) {
	if len(args) < 1 {
		printUsage("Please specify the packages to index", IndexUsageMsg)
		os.Exit(1)
	} else if checkHelpFlag(args) {
		printUsage("Update the search index", IndexUsageMsg)
		return
	}

	dbPath, err := lookupDatabasePath()
	if err != nil {
		exitWithError(err)
	}

	packages := args
	if err := cmd.Index(dbPath, packages); err != nil {
		exitWithError(err)
	}
}

func handleFindCmd(args []string) {
	if len(args) < 1 {
		printUsage("Please provide a search query", FindUsageMsg)
		os.Exit(1)
	} else if checkHelpFlag(args) {
		printUsage("Find definitions from the search index", FindUsageMsg)
		return
	}

	dbPath, err := lookupDatabasePath()
	if err != nil {
		exitWithError(err)
	}

	query, packages := args[0], args[1:len(args)]
	if len(packages) == 0 {
		packages = []string{"."}
	}

	if err := cmd.Find(dbPath, query, packages); err != nil {
		exitWithError(err)
	}
}

func handleInvalidCmd(cmd string) {
	if cmd == "" {
		fmt.Fprintf(os.Stderr, "No command provided\n\n")
	} else {
		fmt.Fprintf(os.Stderr, "Invalid command: %s\n\n", cmd)
	}
	fmt.Fprintf(os.Stderr, "%s\n\n", UsageMsg)
	os.Exit(1)
}

func lookupDatabasePath() (string, error) {
	path, err := xdg.DataFile("gospelunk/index.db")
	if err != nil {
		return "", errors.Wrapf(err, "xdg.DataFile")
	}
	return path, nil
}

func checkHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" {
			return true
		}
	}
	return false
}

func printUsage(preamble string, usageMsg string) {
	fmt.Fprintf(os.Stderr, "%s\n\n%s\n\n", preamble, usageMsg)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}
