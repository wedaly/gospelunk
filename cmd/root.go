package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gospelunk",
	Short: "quickly find definitions in Go projects",
	Long:  "gospelunk is a command-line tool for quickly finding definitions in Go projects",
}

func Execute() error {
	return rootCmd.Execute()
}
