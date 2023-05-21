package jvm

import (
	"fmt"
	"jvm-finder/internal/log"
	"os/exec"
	"strings"
	"time"
)

type Jvm struct {
	javaPath                 string
	JavaHome                 string
	JavaSpecificationVersion uint
	JavaVendor               string
	FetchedAt                time.Time
	SystemProperties         map[string]string
}

func (jvm *Jvm) rebuild() error {
	jvm.JavaHome = jvm.SystemProperties["java.home"]
	jvm.JavaVendor = jvm.SystemProperties["java.vendor"]
	if specVersion, err := parseJavaSpecificationVersion(jvm.SystemProperties["java.specification.version"]); err != nil {
		return err
	} else {
		jvm.JavaSpecificationVersion = specVersion
	}
	return nil
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
		jvm.JavaHome,
		jvm.JavaSpecificationVersion)
}

func fetchJvmInfo(javaPath string) (*Jvm, error) {
	cmd := exec.Command(javaPath, "-cp", "build/classes", "JvmMetadataExtractor")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, log.WrapErr(err, "fail to call %s ", javaPath)
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
	if err := jvmInfo.rebuild(); err != nil {
		return nil, err
	}
	return &jvmInfo, nil
}
