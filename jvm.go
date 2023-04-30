package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Jvm struct {
	javaPath                 string
	javaHome                 string
	javaSpecificationVersion uint
	FetchedAt                time.Time
	SystemProperties         map[string]string
}

func (jvm *Jvm) rebuild() {
	jvm.javaHome = jvm.SystemProperties["java.home"]
	jvm.javaSpecificationVersion = parseVersion(jvm.SystemProperties["java.specification.version"])
}

func (jvm *Jvm) String() string {
	return fmt.Sprintf(
		`[%v]
timestamp: %s
java.home: %s
java.specification.version: %d
`,
		jvm.javaPath,
		jvm.FetchedAt,
		jvm.javaHome,
		jvm.javaSpecificationVersion)
}

func fetchJvmInfo(javaPath string) *Jvm {
	cmd := exec.Command(javaPath, "-cp", "build/classes", "JvmInfo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		die("Fail to call %s %s", javaPath, err)
	}
	lines := strings.Split(string(output), "\n")
	systemProperties := make(map[string]string)
	for _, line := range lines {
		split := strings.SplitN(line, "=", 2)
		if len(split) == 2 {
			systemProperties[split[0]] = strings.TrimSpace(split[1])
		}
	}
	jvmInfo := Jvm{
		javaPath:         javaPath,
		FetchedAt:        time.Now(),
		SystemProperties: systemProperties,
	}
	jvmInfo.rebuild()
	return &jvmInfo
}
