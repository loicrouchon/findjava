package main

import (
	"fmt"
	"io"
	"os"
)

type Console struct {
	stdout io.Writer
	stderr io.Writer
}

var console = Console{
	stdout: os.Stdout,
	stderr: os.Stderr,
}

const logLevelError = 0
const logLevelWarning = logLevelError + 1
const logLevelInfo = logLevelWarning + 1
const logLevelDebug = logLevelInfo + 1

var currentLogLevel uint

func wrapErr(err error, message string, v ...any) error {
	return fmt.Errorf("%s\n\t%s", fmt.Sprintf(message, v...), err)
}

func printf(message string, v ...any) {
	fmt.Fprintf(console.stdout, message, v...)
}

func setLogLevel(level string) error {
	switch level {
	case "debug":
		currentLogLevel = logLevelDebug
	case "info":
		currentLogLevel = logLevelInfo
	case "warn":
		currentLogLevel = logLevelWarning
	case "error":
		currentLogLevel = logLevelError
	default:
		currentLogLevel = logLevelError
		return fmt.Errorf("invalid log level: '%s'. Available levels are: debug, info, warn, error", level)
	}
	return nil
}

func logDebug(message string, v ...any) {
	if currentLogLevel >= logLevelDebug {
		fmt.Fprintf(console.stdout, "[DEBUG] %s\n", fmt.Sprintf(message, v...))
	}
}
func logInfo(message string, v ...any) {
	if currentLogLevel >= logLevelInfo {
		fmt.Fprintf(console.stdout, "[INFO] %s\n", fmt.Sprintf(message, v...))
	}
}

func logWarn(err error) {
	if currentLogLevel >= logLevelWarning {
		fmt.Fprintf(console.stderr, "[WARNING] %s\n", err)
	}
}

func logErr(err error) {
	if currentLogLevel >= logLevelError {
		fmt.Fprintf(console.stderr, "[ERROR] %s\n", err)
	}
}

func die(err error) {
	logErr(err)
	os.Exit(1)
}
