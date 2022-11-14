package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/wedaly/gospelunk/pkg/file"
	"github.com/wedaly/gospelunk/pkg/inspect"
	"github.com/wedaly/gospelunk/pkg/output"
)

var (
	InspectFileArg     string
	InspectLineArg     int
	InspectColumnArg   int
	InspectTemplateArg string
)

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "inspect Go code",
	Long:  "inspect type information and definition location for an identifier in a Go source file",
	RunE: func(cmd *cobra.Command, args []string) error {
		tmpl, err := output.Template(InspectTemplateArg)
		if err != nil {
			return err
		}

		loc := file.Loc{
			Path:   InspectFileArg,
			Line:   InspectLineArg,
			Column: InspectColumnArg,
		}
		result, err := inspect.Inspect(loc)
		if err != nil {
			return err
		}

		if result == nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "No results found for %s", loc)
			return nil
		}

		err = tmpl.Execute(cmd.OutOrStdout(), result)
		if err != nil {
			return errors.Wrapf(err, "template.Execute")
		}

		return nil
	},
}

func init() {
	inspectCmd.Flags().StringVarP(&InspectFileArg, "file", "f", "", "Go source file")
	inspectCmd.MarkFlagRequired("file")

	inspectCmd.Flags().IntVarP(&InspectLineArg, "line", "l", 1, "Line number in Go source file")
	inspectCmd.MarkFlagRequired("line")

	inspectCmd.Flags().IntVarP(&InspectColumnArg, "column", "c", 1, "Column number in Go source file")
	inspectCmd.MarkFlagRequired("column")

	defaultTpl := "{{range .Relations}}{{.Name}} {{.Path|RelPath}}:{{.Line}}:{{.Column}}{{end}}\n"
	inspectCmd.Flags().StringVarP(&InspectTemplateArg, "template", "t", defaultTpl, "Go template for formatting result output")

	rootCmd.AddCommand(inspectCmd)
}
