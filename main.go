package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		die(err)
	}
	config, err := loadConfig("/etc/jvm-finder/config.json", args.configKey)
	if err != nil {
		die(err)
	}
	javaExecutables, err := findAllJavaExecutables(&config.jvmsLookupPaths)
	if err != nil {
		die(err)
	}
	jvmInfos, err := loadJvmsInfos(config.jvmsMetadataCachePath, &javaExecutables)
	if err != nil {
		die(err)
	}
	rules := jvmSelectionRules(config, args.minJavaVersion, args.maxJavaVersion, args.vendors)
	if jvms := jvmInfos.Select(rules); len(jvms) > 0 {
		jvm := jvms[0]
		logJvmList("[SELECTED]", jvms[0:1])
		printf("%s\n", filepath.Join(jvm.javaHome, "bin", "java"))
	} else {
		die(fmt.Errorf("unable to find a JVM matching requirements %s", rules))
	}
}
