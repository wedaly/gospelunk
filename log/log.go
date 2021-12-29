package log

import (
	"fmt"
	"os"
)

var quietFlag bool

func Init(quiet bool) {
	quietFlag = quiet
}

func Info(msg string, args ...interface{}) {
	if quietFlag {
		return
	}
	fmt.Fprintf(os.Stderr, msg, args...)
}

func Warn(msg string, args ...interface{}) {
	if quietFlag {
		return
	}
	fmt.Fprintf(os.Stderr, "[WARN] ")
	fmt.Fprintf(os.Stderr, msg, args...)
}

func Error(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[ERROR] ")
	fmt.Fprintf(os.Stderr, msg, args...)
}
