package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type JvmsInfos struct {
	path       string
	Jvms       map[string]*JvmInfo
	dirtyCache bool
}

type JvmInfo struct {
	javaPath                 string
	javaHome                 string
	javaSpecificationVersion uint
	fetched                  bool
	FetchedAt                time.Time
	SystemProperties         map[string]string
}

func (jvm *JvmInfo) rebuild() {
	jvm.javaHome = jvm.SystemProperties["java.home"]
	jvm.javaSpecificationVersion = parseVersion(jvm.SystemProperties["java.specification.version"])
}

func (jvm *JvmInfo) String() string {
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

func loadJvmsInfos(path string, javaPaths *JavaExecutables) JvmsInfos {
	jvmInfos := loadJvmsInfosFromCache(path)
	for javaPath, modTime := range javaPaths.javaPaths {
		jvmInfos.Fetch(javaPath, modTime)
	}
	jvmInfos.Save()
	return jvmInfos
}
func loadJvmsInfosFromCache(path string) JvmsInfos {
	jvmsInfos := JvmsInfos{
		path:       path,
		dirtyCache: false,
		Jvms:       make(map[string]*JvmInfo),
	}
	if _, err := os.Stat(path); err == nil {
		logDebug("Loading cache from %s", path)
		file, _ := os.Open(path)
		defer closeFile(file)
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&jvmsInfos)
		if err != nil {
			dierr(err)
		}
		for javaPath, jvm := range jvmsInfos.Jvms {
			jvm.javaPath = javaPath
			jvm.rebuild()
		}
		// logDebug("JVMs rebuilt loaded from cache: %#v", jvmsInfos)
	}
	return jvmsInfos
}

func (jvmInfos *JvmsInfos) Fetch(javaPath string, modTime time.Time) {
	var jvmInfo *JvmInfo
	if info, found := jvmInfos.Jvms[javaPath]; !found {
		logInfo("[CACHE MISS] %s", javaPath)
		jvmInfo = jvmInfos.doFetch(javaPath)
	} else if modTime.After(info.FetchedAt) {
		logInfo("[CACHE OUTDATED] %s", javaPath)
		jvmInfo = jvmInfos.doFetch(javaPath)
	} else {
		jvmInfo = info
	}
	jvmInfo.fetched = true
	jvmInfos.Jvms[javaPath] = jvmInfo
}

func (jvmInfos *JvmsInfos) doFetch(javaPath string) *JvmInfo {
	jvmInfo := fetchJvmInfo(javaPath)
	jvmInfos.dirtyCache = true
	logDebug("%s:\n%s", javaPath, jvmInfo)
	return jvmInfo
}

func fetchJvmInfo(javaPath string) *JvmInfo {
	cmd := exec.Command(javaPath, "-cp", "build/classes", "JvmInfo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		die("Fail to call %s %s", javaPath, err)
	}
	lines := strings.Split(string(output), "\n")
	systemProperties := make(map[string]string)
	for _, line := range lines {
		split := strings.SplitN(line, "=", 2)
		if split[0] == "java.home" || split[0] == "java.specification.version" {
			systemProperties[split[0]] = strings.TrimSpace(split[1])
		}
	}
	jvmInfo := JvmInfo{
		javaPath:         javaPath,
		fetched:          true,
		FetchedAt:        time.Now(),
		SystemProperties: systemProperties,
	}
	jvmInfo.rebuild()
	return &jvmInfo
}

func (jvmsInfos *JvmsInfos) Save() {
	for javaPath, jvmInfo := range jvmsInfos.Jvms {
		if !jvmInfo.fetched {
			if fileInfo, err := os.Stat(javaPath); err == nil {
				if fileInfo.ModTime().After(jvmInfo.FetchedAt) {
					jvmsInfos.doFetch(javaPath)
				}
			} else {
				delete(jvmsInfos.Jvms, javaPath)
				jvmsInfos.dirtyCache = true
			}
		}
	}
	if jvmsInfos.dirtyCache {
		writeToJson(jvmsInfos)
	}
}

func writeToJson(jvmInfos *JvmsInfos) {
	logDebug("Writing JVMs infos cache to %s", jvmInfos.path)
	file, _ := json.MarshalIndent(jvmInfos, "", "  ")
	err := os.WriteFile(jvmInfos.path, file, 0644)
	if err != nil {
		die("Unable to write to file %s, %s", jvmInfos.path, err)
	}
}
