package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := parseArgs()
	config := loadConfig("/etc/jvm-finder/config.json", args.configKey)
	javaExecutables := findAllJavaExecutables(&config.jvmsLookupPaths)
	jvmInfos := loadJvmsInfos(config.jvmsMetadataCachePath, &javaExecutables)
	rules := jvmSelectionRules(config, args.minJavaVersion, args.maxJavaVersion, args.vendors)
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
	vendors        list
}

type list []string

func (i *list) String() string {
	return "[" + strings.Join(*i, ", ") + "]"
}
func (i *list) Set(value string) error {
	*i = append(*i, value)
	return nil
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
	flag.Var(&args.vendors, "vendors",
		"The vendors to filter on. If empty, no vendor filtering will be done")
	flag.Parse()
	if len(flag.Args()) > 0 {
		flag.Usage()
		os.Exit(1)
	}
	if err := setLogLevel(args.logLevel); err != nil {
		dierr(err)
	}
	return &args
}
