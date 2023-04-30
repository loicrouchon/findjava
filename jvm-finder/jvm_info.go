package main

import (
	"encoding/json"
	"os"
	"time"
)

type JvmsInfos struct {
	path       string
	dirtyCache bool
	fetched    map[string]bool
	Jvms       map[string]*Jvm
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
		fetched:    make(map[string]bool),
		Jvms:       make(map[string]*Jvm),
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

func (jvms *JvmsInfos) Fetch(javaPath string, modTime time.Time) {
	var jvmInfo *Jvm
	if info, found := jvms.Jvms[javaPath]; !found {
		logInfo("[CACHE MISS] %s", javaPath)
		jvmInfo = jvms.doFetch(javaPath)
	} else if modTime.After(info.FetchedAt) {
		logInfo("[CACHE OUTDATED] %s", javaPath)
		jvmInfo = jvms.doFetch(javaPath)
	} else {
		jvmInfo = info
	}
	jvms.fetched[javaPath] = true
	jvms.Jvms[javaPath] = jvmInfo
}

func (jvms *JvmsInfos) doFetch(javaPath string) *Jvm {
	jvmInfo := fetchJvmInfo(javaPath)
	jvms.dirtyCache = true
	logDebug("%s:\n%s", javaPath, jvmInfo)
	return jvmInfo
}

func (jvms *JvmsInfos) Save() {
	for javaPath, jvmInfo := range jvms.Jvms {
		if value, found := jvms.fetched[javaPath]; !found || !value {
			if fileInfo, err := os.Stat(javaPath); err == nil {
				if fileInfo.ModTime().After(jvmInfo.FetchedAt) {
					jvms.doFetch(javaPath)
				}
			} else {
				delete(jvms.Jvms, javaPath)
				jvms.dirtyCache = true
			}
		}
	}
	if jvms.dirtyCache {
		writeToJson(jvms)
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
