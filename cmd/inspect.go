package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/wedaly/gospelunk/pkg/inspect"
)

var (
	FileArg     string
	LineArg     int
	ColumnArg   int
	TemplateArg string
)

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "inspect Go code",
	Long:  "inspect type information and definition location for an identifier in a Go source file",
	RunE: func(cmd *cobra.Command, args []string) error {
		tmpl, err := template.New("template").
			Funcs(template.FuncMap{"RelPath": RelPathTplFunc()}).
			Parse(TemplateArg)
		if err != nil {
			return errors.Wrapf(err, "template.Parse")
		}

		loc := inspect.FileLoc{Path: FileArg, Line: LineArg, Column: ColumnArg}
		result, err := inspect.Inspect(loc)
		if err != nil {
			return err
		}

		if result == nil {
			fmt.Fprintf(os.Stderr, "No results found for %s:%d:%d\n", FileArg, LineArg, ColumnArg)
			return nil
		}

		err = tmpl.Execute(os.Stdout, result)
		if err != nil {
			return errors.Wrapf(err, "template.Execute")
		}

		return nil
	},
}

func RelPathTplFunc() func(string) string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	return func(path string) string {
		relpath, err := filepath.Rel(cwd, path)
		if err != nil {
			relpath = path
		}
		return relpath
	}
}

func init() {
	inspectCmd.Flags().StringVarP(&FileArg, "file", "f", "", "Go source file")
	inspectCmd.MarkFlagRequired("file")

	inspectCmd.Flags().IntVarP(&LineArg, "line", "l", 1, "Line number in Go source file")
	inspectCmd.MarkFlagRequired("line")

	inspectCmd.Flags().IntVarP(&ColumnArg, "column", "c", 1, "Column number in Go source file")
	inspectCmd.MarkFlagRequired("column")

	inspectCmd.Flags().StringVarP(&TemplateArg, "template", "t", "{{.Def.Name}} {{.Def.Path|RelPath}}:{{.Def.Line}}:{{.Def.Column}}\n", "Go template for formatting result output")

	rootCmd.AddCommand(inspectCmd)
}
