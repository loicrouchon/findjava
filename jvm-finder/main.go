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
	rules := jvmSelectionRules(args.minJavaVersion, args.maxJavaVersion, config)
	javaExecutables := findAllJavaExecutables(config.jvmsLookupPaths())
	jvmInfos := loadJvmInfos(config.jvmsMetadataCachePath(), &javaExecutables)
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
	flag.StringVar(&args.logLevel, "log-level", "error",
		"Sets the log level to one of: debug, info, error")
	flag.StringVar(&args.configKey, "config-key", "",
		"If specified, will look for an optional config.<KEY>.json to load before loading the default configuration")
	flag.UintVar(&args.minJavaVersion, "min-java-version", allVersions,
		"The minimum (inclusive) Java Language Specification version the found JVMs should provide")
	flag.UintVar(&args.maxJavaVersion, "max-java-version", allVersions,
		"The maximum (inclusive) Java Language Specification version the found JVMs should provide")
	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(0)
	}
	if err := setLogLevel(args.logLevel); err != nil {
		dierr(err)
	}
	return &args
}
