package log

import (
	"fmt"
	"os"
)

var verboseFlag bool

func Init(verbose bool) {
	verboseFlag = verbose
}

func Info(msg string, args ...interface{}) {
	if verboseFlag {
		fmt.Fprintf(os.Stderr, msg, args...)
	}
}

func Warn(msg string, args ...interface{}) {
	if verboseFlag {
		fmt.Fprintf(os.Stderr, "[WARN] ")
		fmt.Fprintf(os.Stderr, msg, args...)
	}
}

func Error(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ERROR] ")
	fmt.Fprintf(os.Stderr, msg, args...)
}
