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

const logLevelDebug = 2
const logLevelInfo = 1
const logLevelError = 0

var currentLogLevel uint

func printf(message string, v ...any) {
	fmt.Fprintf(console.stdout, message, v...)
}

func setLogLevel(level string) error {
	switch level {
	case "debug":
		currentLogLevel = logLevelDebug
	case "info":
		currentLogLevel = logLevelInfo
	case "error":
		currentLogLevel = logLevelError
	default:
		currentLogLevel = logLevelError
		return fmt.Errorf("invalid log level: '%s'. Available levels are: debug, info, error", level)
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

func logError(message string, v ...any) {
	if currentLogLevel >= logLevelError {
		fmt.Fprintf(console.stderr, "[ERROR] %s\n", fmt.Sprintf(message, v...))
	}
}

func logErr(err error) {
	if currentLogLevel >= logLevelError {
		fmt.Fprintf(console.stderr, "[ERROR] %s\n", err)
	}
}

func die(message string, v ...any) {
	logError(message, v...)
	os.Exit(1)
}

func dierr(err error) {
	logErr(err)
	os.Exit(1)
}
