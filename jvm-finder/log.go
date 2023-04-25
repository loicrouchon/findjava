package main

import (
	"fmt"
	"os"
)

const logLevelDebug = 2
const logLevelInfo = 1
const logLevelError = 0

var currentLogLevel int

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
		return fmt.Errorf("Invalid log level: '%s'. Available levels are: debug, info, error", level)
	}
	return nil
}

func logDebug(message string, v ...any) {
	if currentLogLevel >= logLevelDebug {
		fmt.Fprintf(os.Stdout, "[DEBUG] %s\n", fmt.Sprintf(message, v...))
	}
}
func logInfo(message string, v ...any) {
	if currentLogLevel >= logLevelInfo {
		fmt.Fprintf(os.Stdout, "[INFO] %s\n", fmt.Sprintf(message, v...))
	}
}

func logError(message string, v ...any) {
	if currentLogLevel >= logLevelError {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", fmt.Sprintf(message, v...))
	}
}

func logErr(err error) {
	if currentLogLevel >= logLevelError {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
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
