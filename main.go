package main

import (
	"os"
	"path/filepath"
)

func main() {
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		dierr(err)
	}
	config, err := loadConfig("/etc/jvm-finder/config.json", args.configKey)
	if err != nil {
		dierr(err)
	}
	javaExecutables, err := findAllJavaExecutables(&config.jvmsLookupPaths)
	if err != nil {
		dierr(err)
	}
	jvmInfos := loadJvmsInfos(config.jvmsMetadataCachePath, &javaExecutables)
	rules := jvmSelectionRules(config, args.minJavaVersion, args.maxJavaVersion, args.vendors)
	if jvm := jvmInfos.Select(rules); jvm != nil {
		logInfo("[SELECTED]  %s (%d)", jvm.javaHome, jvm.javaSpecificationVersion)
		printf("%s\n", filepath.Join(jvm.javaHome, "bin", "java"))
	} else {
		die("Unable to find a JVM matching requirements %s", rules)
	}
}
