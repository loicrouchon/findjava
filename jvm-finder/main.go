package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	args := parseArgs()
	config := loadConfig("/etc/jvm-finder/config.json", "default")
	var rules *JvmSelectionRules
	if len(args) > 1 {
		Usage()
	} else if len(args) == 1 {
		rules = jvmSelectionRules(&args[0], config)
		if rules == nil {
			Usage()
		}
	} else {
		rules = jvmSelectionRules(nil, config)
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

func parseArgs() []string {
	var logLevel string
	flag.StringVar(&logLevel, "loglevel", "error", "Log level: debug, info, error")
	flag.Parse()
	if err := setLogLevel(logLevel); err != nil {
		dierr(err)
	}
	return flag.Args()
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
