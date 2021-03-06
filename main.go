package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adrg/xdg"
	"github.com/pkg/errors"

	"github.com/wedaly/gospelunk/cmd"
	"github.com/wedaly/gospelunk/log"
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
	fs := flag.NewFlagSet("index", flag.ExitOnError)
	fs.Usage = func() {
		printUsage("Update the search index", IndexUsageMsg)
		fmt.Fprintf(os.Stderr, "Flags:\n\n")
		fs.PrintDefaults()
	}

	quietFlagArg := fs.Bool("q", false, "quiet")
	includeImportsArg := fs.Bool("i", false, "include imported packages")
	transitiveArg := fs.Bool("t", false, "include transitive dependencies")

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	log.Init(*quietFlagArg)

	posArgs := fs.Args()
	if len(posArgs) < 1 {
		printUsage("Please specify the packages to index", IndexUsageMsg)
		os.Exit(1)
	}

	if checkHelpFlag(posArgs) {
		fs.Usage()
		return
	}

	dbPath, err := lookupDatabasePath()
	if err != nil {
		exitWithError(err)
	}

	packages := posArgs
	if err := cmd.Index(dbPath, packages, *includeImportsArg, *transitiveArg); err != nil {
		exitWithError(err)
	}
}

func handleFindCmd(args []string) {
	fs := flag.NewFlagSet("find", flag.ExitOnError)
	fs.Usage = func() {
		printUsage("Find definitions from the search index", FindUsageMsg)
		fmt.Fprintf(os.Stderr, "Flags:\n\n")
		fs.PrintDefaults()
	}

	quietFlagArg := fs.Bool("q", false, "quiet")
	includeImportsArg := fs.Bool("i", false, "include imported packages")
	formatTplArg := fs.String("f", "{{ .Path | RelPath }}:{{ .LineNum }} {{ .Kind }} {{ .Name }}", "format the output using Go template syntax")

	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	log.Init(*quietFlagArg)

	posArgs := fs.Args()
	if len(posArgs) < 1 {
		printUsage("Please provide a search query", FindUsageMsg)
		os.Exit(1)
	}

	if checkHelpFlag(posArgs) {
		fs.Usage()
		return
	}

	query, packages := posArgs[0], posArgs[1:len(posArgs)]
	if len(packages) == 0 {
		packages = []string{"."}
	}

	dbPath, err := lookupDatabasePath()
	if err != nil {
		exitWithError(err)
	}

	if err := cmd.Find(dbPath, query, packages, *includeImportsArg, *formatTplArg); err != nil {
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
	log.Error("%s\n", err)
	os.Exit(1)
}
