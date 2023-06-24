package jvm

import (
	"findjvm/internal/log"
	"os/exec"
	"strings"
	"time"
)

type MetadataReader struct {
	Classpath string
}

func (f *MetadataReader) fetchJvmInfo(javaPath string) (*Jvm, error) {
	cmd := exec.Command(javaPath, "-cp", f.Classpath, "JvmMetadataExtractor")
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
