package cmd

import (
	"bytes"

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

func ExecuteInTest(args []string) (stdout string, stderr string, err error) {
	// Configure the cmd to capture stdout and stderr and use test-provided args.
	var stdoutBuf, stderrBuf bytes.Buffer
	rootCmd.SetArgs(args)
	rootCmd.SetOut(&stdoutBuf)
	rootCmd.SetErr(&stderrBuf)
	defer func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
	}()

	// Execute the cmd and return the output.
	err = rootCmd.Execute()
	return stdoutBuf.String(), stderrBuf.String(), err
}
