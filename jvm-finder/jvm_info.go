package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type JvmInfos struct {
	path       string
	jvmInfos   map[string]JvmInfo
	dirtyCache bool
}

type JvmInfo struct {
	javaPath                 string
	javaHome                 string
	javaSpecificationVersion int
}

func (jvmInfo JvmInfo) String() string {
	return fmt.Sprintf(
		`[%s]
java.home: %s
java.specification.version: %d
`,
		jvmInfo.javaPath,
		jvmInfo.javaHome,
		jvmInfo.javaSpecificationVersion)
}

func loadJvmInfos(path string, javaPaths []string) JvmInfos {
	infos := make(map[string]JvmInfo)
	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		if err != nil {
			logError("Unable to read file %s: %s", path, err)
		}
		defer file.Close()

		var jvmInfo JvmInfo
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if value, found := strings.CutPrefix(line, "["); found {
				if len(jvmInfo.javaHome) > 0 && jvmInfo.javaSpecificationVersion > 0 {
					infos[jvmInfo.javaPath] = jvmInfo
				}
				jvmInfo = JvmInfo{javaPath: strings.Trim(value, "[]")}
			} else if value, ok := strings.CutPrefix(line, "java.home="); ok {
				jvmInfo.javaHome = value
			} else if value, ok := strings.CutPrefix(line, "java.specification.version="); ok {
				if version, err := strconv.Atoi(value); err == nil {
					jvmInfo.javaSpecificationVersion = version
				}
			}
		}
		if len(jvmInfo.javaHome) > 0 && jvmInfo.javaSpecificationVersion > 0 {
			infos[jvmInfo.javaPath] = jvmInfo
		}

		if err := scanner.Err(); err != nil {
			logError("error", err)
		}
	}
	jvmInfos := JvmInfos{
		path:       path,
		jvmInfos:   infos,
		dirtyCache: false,
	}

	for _, javaPath := range javaPaths {
		jvmInfos.Fetch(javaPath)
	}
	jvmInfos.Save()
	return jvmInfos
}

func (cache *JvmInfos) Fetch(javaPath string) {
	if _, found := cache.jvmInfos[javaPath]; !found {
		logInfo("[CACHE MISS] %s", javaPath)
		jvmInfo := fetchJvmInfo(javaPath)
		cache.jvmInfos[javaPath] = jvmInfo
		cache.dirtyCache = true
		logDebug("%s: %s", javaPath, jvmInfo)
	}
}

func fetchJvmInfo(javaPath string) JvmInfo {
	cmd := exec.Command(javaPath, "-cp", "build/classes", "JvmInfo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logError("Fail to call %s %s", javaPath, err)
		os.Exit(1)
	}
	lines := strings.Split(string(output), "\n")
	var javaSpecificationVersion string
	var javaHome string
	for _, line := range lines {
		if value, found := strings.CutPrefix(line, "java.home="); found {
			javaHome = strings.TrimSpace(value)
		}
		if value, found := strings.CutPrefix(line, "java.specification.version="); found {
			javaSpecificationVersion = strings.TrimSpace(value)
		}
	}
	return JvmInfo{
		javaPath:                 javaPath,
		javaHome:                 javaHome,
		javaSpecificationVersion: parseVersion(javaSpecificationVersion),
	}
}

func (cache *JvmInfos) Save() {
	if cache.dirtyCache {
		output := ""
		for _, jvmInfo := range cache.jvmInfos {
			output = fmt.Sprintf(`%s
[%s]
java.home=%s
java.specification.version=%d
`, output, jvmInfo.javaPath, jvmInfo.javaHome, jvmInfo.javaSpecificationVersion)
		}
		logDebug(output)
		if err := os.WriteFile(cache.path, []byte(output), 0666); err != nil {
			logError("Unable to write to file %s, %s", cache.path, err)
			os.Exit(1)
		}
	}
}
