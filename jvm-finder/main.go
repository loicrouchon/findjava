package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	args := parseArgs()
	config := loadConfig("/etc/jvm-finder/config.json", args.configKey)
	var rules = jvmSelectionRules(args.minJavaVersion, args.maxJavaVersion, config)
	if rules == nil {
		Usage()
	}
	javaExecutables := findAllJavaExecutables(config.jvmLookupPaths())
	jvmInfos := loadJvmInfos("./build/jvm-finder.properties", &javaExecutables)
	if jvm := jvmInfos.Select(rules); jvm != nil {
		logInfo("[SELECTED]  %s (%d)", jvm.javaHome, jvm.javaSpecificationVersion)
		fmt.Printf("%s\n", filepath.Join(jvm.javaHome, "bin", "java"))
	} else {
		die("Unable to find a JVM matching requirements %s", rules)
	}
}

type Args struct {
	logLevel       string
	configKey      string
	minJavaVersion uint
	maxJavaVersion uint
}

func parseArgs() *Args {
	args := Args{}
	flag.StringVar(&args.logLevel, "log-level", "error", "Log level: debug, info, error")
	flag.StringVar(&args.configKey, "config-key", "default", "The configuration to load")
	flag.UintVar(&args.minJavaVersion, "min-java-version", allVersions, "The minimum (inclusive) Java Language Specification version to search for.")
	flag.UintVar(&args.maxJavaVersion, "max-java-version", allVersions, "The maximum (inclusive) Java Language Specification version to search for.")
	flag.Parse()
	if err := setLogLevel(args.logLevel); err != nil {
		dierr(err)
	}
	if len(flag.Args()) > 0 {
		Usage()
	}
	return &args
}

func Usage() {
	logError("Usage jvm-finder [VERSION]")
	logError("  VERSION: A JVM version range:")
	logError("      - 17        exact version)")
	logError("      - 17..      17 or above)")
	logError("      - ..17      up to 17")
	logError("      - 11..17    From 11 to 17")
	os.Exit(1)
}
