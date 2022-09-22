package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/wedaly/gospelunk/pkg/list"
	"github.com/wedaly/gospelunk/pkg/output"
)

var (
	ListTemplateArg                string
	ListIncludeStructFieldsArg     bool
	ListIncludeInterfaceMethodsArg bool
	ListIncludePrivateArg          bool
	ListIncludeTestsArg            bool
)

var listCmd = &cobra.Command{
	Use:   "list [flags] [packages]",
	Short: "list definitions in Go packages",
	Long:  "list definitions of functions, consts, structs, interfaces, and methods from Go packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		tmpl, err := output.Template(ListTemplateArg)
		if err != nil {
			return err
		}

		patterns := args // Passed to Go build system to locate packages.
		opts := list.Options{
			IncludeStructFields:     ListIncludeStructFieldsArg,
			IncludeInterfaceMethods: ListIncludeInterfaceMethodsArg,
			IncludePrivate:          ListIncludePrivateArg,
			IncludeTests:            ListIncludeTestsArg,
		}
		result, err := list.List(patterns, opts)
		if err != nil {
			return err
		}

		err = tmpl.Execute(cmd.OutOrStdout(), result)
		if err != nil {
			return errors.Wrapf(err, "template.Execute")
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&ListTemplateArg, "template", "t", "{{ range .Defs }}{{.Name}} {{.Path|RelPath}}:{{.Line}}:{{.Column}}\n{{end}}", "Go template for formatting result output")
	listCmd.Flags().BoolVar(&ListIncludeStructFieldsArg, "include-struct-fields", false, "Include struct fields")
	listCmd.Flags().BoolVar(&ListIncludeInterfaceMethodsArg, "include-interface-methods", false, "Include interface methods")
	listCmd.Flags().BoolVarP(&ListIncludePrivateArg, "include-private", "p", false, "Include private definitions")
	listCmd.Flags().BoolVar(&ListIncludeTestsArg, "include-tests", false, "Include definitions from tests")
	rootCmd.AddCommand(listCmd)
}
