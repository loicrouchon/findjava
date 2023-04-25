package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var javaLookUpPaths = []string{
	"/bin/java",
	"/usr/bin/java",
	"/usr/local/bin/java",
	"/usr/lib/jvm",
	"~/.sdkman/candidates/java",
}

func main() {
	args := parseArgs()
	if len(args) > 1 {
		logError("Usage jvm-finder [VERSION]")
		logError("  VERSION: A JVM version range:")
		logError("      - 17        exact version)")
		logError("      - 17..      17 or above)")
		logError("      - ..17      up to 17")
		logError("      - 11..17    From 11 to 17")
		os.Exit(1)
	}
	rules := jvmSelectionRules(args)
	javaExecutables := findAllJavaPaths(javaLookUpPaths)
	jvmInfos := loadJvmInfos("./build/jvm-finder.properties", &javaExecutables)
	if jvm, found := jvmInfos.Select(rules); found {
		logInfo("[SELECTED]  %s (%d)", jvm.javaHome, jvm.javaSpecificationVersion)
		fmt.Printf("%s\n", filepath.Join(jvm.javaHome, "bin", "java"))
	} else {
		logError("Unable to find a JVM matching requirements %s", rules)
		os.Exit(1)
	}
}

var logLevel string

func parseArgs() []string {
	flag.StringVar(&logLevel, "loglevel", "error", "Log level: debug, info, error")
	flag.Parse()
	if logLevel != "debug" && logLevel != "info" && logLevel != "error" {
		logError("Invalid log level: '%s'. Available levels are: debug, info, error", logLevel)
		os.Exit(1)
	}
	return flag.Args()
}

func logDebug(message string, v ...any) {
	if logLevel == "debug" {
		fmt.Fprintf(os.Stdout, "[DEBUG] %s\n", fmt.Sprintf(message, v...))
	}
}
func logInfo(message string, v ...any) {
	if logLevel == "debug" || logLevel == "info" {
		fmt.Fprintf(os.Stdout, "[INFO] %s\n", fmt.Sprintf(message, v...))
	}
}

func logError(message string, v ...any) {
	fmt.Fprintf(os.Stderr, "[ERROR] %s\n", fmt.Sprintf(message, v...))
}
