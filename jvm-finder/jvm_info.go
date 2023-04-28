package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type JvmInfos struct {
	path       string
	timestamp  time.Time
	jvmInfos   map[string]JvmInfo
	dirtyCache bool
}

type JvmInfo struct {
	javaPath                 string
	javaHome                 string
	javaSpecificationVersion uint
	fetched                  bool
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

func loadJvmsInfos(path string, javaPaths *JavaExecutables) JvmInfos {
	jvmInfos := loadJvmsInfosFromCache(path)
	for javaPath, modTime := range javaPaths.javaPaths {
		jvmInfos.Fetch(javaPath, modTime)
	}
	jvmInfos.Save()
	return jvmInfos
}
func loadJvmsInfosFromCache(path string) JvmInfos {
	var timestamp time.Time
	infos := make(map[string]JvmInfo)
	if fileinfo, err := os.Stat(path); err == nil {
		timestamp = fileinfo.ModTime()
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
				registerIfComplete(jvmInfo, infos)
				jvmInfo = JvmInfo{
					javaPath: strings.Trim(value, "[]"),
					fetched:  false,
				}
			} else if value, ok := strings.CutPrefix(line, "java.home="); ok {
				jvmInfo.javaHome = value
			} else if value, ok := strings.CutPrefix(line, "java.specification.version="); ok {
				if version, err := strconv.Atoi(value); err == nil {
					jvmInfo.javaSpecificationVersion = uint(version)
				}
			}
		}
		registerIfComplete(jvmInfo, infos)
		if err := scanner.Err(); err != nil {
			logErr(err)
		}
	}
	return JvmInfos{
		path:       path,
		timestamp:  timestamp,
		jvmInfos:   infos,
		dirtyCache: false,
	}
}

func registerIfComplete(jvmInfo JvmInfo, infos map[string]JvmInfo) {
	if len(jvmInfo.javaHome) > 0 && jvmInfo.javaSpecificationVersion > 0 {
		infos[jvmInfo.javaPath] = jvmInfo
	}
}

func (jvmInfos *JvmInfos) Fetch(javaPath string, modTime time.Time) {
	var jvmInfo JvmInfo
	if info, found := jvmInfos.jvmInfos[javaPath]; !found {
		logInfo("[CACHE MISS] %s", javaPath)
		jvmInfo = jvmInfos.doFetch(javaPath)
	} else if modTime.After(jvmInfos.timestamp) {
		logInfo("[CACHE OUTDATED] %s", javaPath)
		jvmInfo = jvmInfos.doFetch(javaPath)
	} else {
		jvmInfo = info
	}
	jvmInfo.fetched = true
	jvmInfos.jvmInfos[javaPath] = jvmInfo
}

func (jvmInfos *JvmInfos) doFetch(javaPath string) JvmInfo {
	jvmInfo := fetchJvmInfo(javaPath)
	jvmInfos.dirtyCache = true
	logDebug("%s: %s", javaPath, jvmInfo)
	return jvmInfo
}

func fetchJvmInfo(javaPath string) JvmInfo {
	cmd := exec.Command(javaPath, "-cp", "build/classes", "JvmInfo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		die("Fail to call %s %s", javaPath, err)
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
		fetched:                  true,
	}
}

func (cache *JvmInfos) Save() {
	for _, jvmInfo := range cache.jvmInfos {
		if !jvmInfo.fetched {
			cache.dirtyCache = true
			break
		}
	}
	if cache.dirtyCache {
		output := ""
		for _, jvmInfo := range cache.jvmInfos {
			jvmInfoAsStr := fmt.Sprintf(`
[%s]
java.home=%s
java.specification.version=%d
`, jvmInfo.javaPath, jvmInfo.javaHome, jvmInfo.javaSpecificationVersion)
			if jvmInfo.fetched {
				output += jvmInfoAsStr
			} else if _, err := os.Stat(jvmInfo.javaPath); err == nil {
				output += jvmInfoAsStr
			} else {
				logInfo("[ORPHAN JVM] %s", jvmInfo.javaPath)
			}
		}
		logDebug(output)
		if err := os.WriteFile(cache.path, []byte(output), 0666); err != nil {
			die("Unable to write to file %s, %s", cache.path, err)
		}
	}
}
