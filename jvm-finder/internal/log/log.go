package log

import (
	"fmt"
	"jvm-finder/internal/console"
	"os"
)

const logLevelError = 0
const logLevelWarning = logLevelError + 1
const logLevelInfo = logLevelWarning + 1
const logLevelDebug = logLevelInfo + 1

var currentLogLevel uint = logLevelError

func SetLogLevel(level string) error {
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
		return fmt.Errorf("invalid log level: \"%s\". Available levels are: debug, info, warn, error", level)
	}
	return nil
}

func Debug(message string, v ...interface{}) {
	if currentLogLevel >= logLevelDebug {
		console.Writer.Printf("[DEBUG] %s\n", fmt.Sprintf(message, v...))
	}
}
func Info(message string, v ...interface{}) {
	if currentLogLevel >= logLevelInfo {
		console.Writer.Printf("[INFO] %s\n", fmt.Sprintf(message, v...))
	}
}

func Warn(err error) {
	if currentLogLevel >= logLevelWarning {
		console.Writer.Eprintf("[WARNING] %s\n", err)
	}
}

func Err(err error) {
	if currentLogLevel >= logLevelError {
		console.Writer.Eprintf("[ERROR] %s\n", err)
	}
}

func Die(err error) {
	Err(err)
	os.Exit(1)
}

func WrapErr(err error, message string, v ...interface{}) error {
	return fmt.Errorf("%s\n\t%s", fmt.Sprintf(message, v...), err)
}
