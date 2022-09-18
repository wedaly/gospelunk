package cmd

import (
	"bytes"
	"os"

	"github.com/pkg/errors"
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

func ExecuteInTest(args []string, cwd string) (stdout string, stderr string, err error) {
	// Set the current working directory and restore on exit.
	oldWd, err := os.Getwd()
	if err != nil {
		return "", "", errors.Wrapf(err, "os.Getwd")
	}
	defer os.Chdir(oldWd)

	err = os.Chdir(cwd)
	if err != nil {
		return "", "", errors.Wrapf(err, "os.Chdir")
	}

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
