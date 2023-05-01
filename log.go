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

func (console *Console) printf(message string, v ...any) {
	_, _ = fmt.Fprintf(console.stdout, message, v...)
}

func (console *Console) eprintf(message string, v ...any) {
	_, _ = fmt.Fprintf(console.stderr, message, v...)
}

const logLevelError = 0
const logLevelWarning = logLevelError + 1
const logLevelInfo = logLevelWarning + 1
const logLevelDebug = logLevelInfo + 1

var currentLogLevel uint

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
		console.printf("[DEBUG] %s\n", fmt.Sprintf(message, v...))
	}
}
func logInfo(message string, v ...any) {
	if currentLogLevel >= logLevelInfo {
		console.printf("[INFO] %s\n", fmt.Sprintf(message, v...))
	}
}

func logWarn(err error) {
	if currentLogLevel >= logLevelWarning {
		console.eprintf("[WARNING] %s\n", err)
	}
}

func logErr(err error) {
	if currentLogLevel >= logLevelError {
		console.eprintf("[ERROR] %s\n", err)
	}
}

func die(err error) {
	logErr(err)
	os.Exit(1)
}

func wrapErr(err error, message string, v ...any) error {
	return fmt.Errorf("%s\n\t%s", fmt.Sprintf(message, v...), err)
}
