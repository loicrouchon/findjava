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

func loadJvmsInfos(path string, javaPaths *JavaExecutables) (JvmsInfos, error) {
	jvmInfos := loadJvmsInfosFromCache(path)
	for javaPath, modTime := range javaPaths.javaPaths {
		if err := jvmInfos.Fetch(javaPath, modTime); err != nil {
			return JvmsInfos{}, err
		}
	}
	jvmInfos.Save()
	return jvmInfos, nil
}

func loadJvmsInfosFromCache(path string) JvmsInfos {
	jvmsInfos := JvmsInfos{
		path:       path,
		dirtyCache: false,
		fetched:    make(map[string]bool),
		Jvms:       make(map[string]*Jvm),
	}
	// Failures to load will from cache will result in an empty JvmsInfos
	// which will cause every discovered JVM to be fetched
	if _, err := os.Stat(path); err == nil {
		logDebug("Loading cache from %s", path)
		if file, err := os.Open(path); err == nil {
			defer closeFile(file)
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&jvmsInfos); err == nil {
				for javaPath, jvm := range jvmsInfos.Jvms {
					jvm.javaPath = javaPath
					jvm.rebuild()
				}
				//logDebug("JVMs rebuilt loaded from cache: %#v", jvmsInfos)
			} else {
				logErr(wrapErr(err, "cannot read config file %s:", path))
			}
		} else {
			logErr(wrapErr(err, "cannot read config file %s:", path))
		}
	}
	return jvmsInfos
}

func (jvms *JvmsInfos) Fetch(javaPath string, modTime time.Time) error {
	jvms.fetched[javaPath] = true
	if info, found := jvms.Jvms[javaPath]; !found {
		logInfo("[CACHE MISS] %s", javaPath)
		return jvms.doFetch(javaPath)
	} else if modTime.After(info.FetchedAt) {
		logInfo("[CACHE OUTDATED] %s", javaPath)
		return jvms.doFetch(javaPath)
	} else {
		return nil
	}
}

func (jvms *JvmsInfos) doFetch(javaPath string) error {
	jvm, err := fetchJvmInfo(javaPath)
	if err != nil {
		return err
	}
	logDebug("%s:\n%s", javaPath, jvm)
	jvms.Jvms[javaPath] = jvm
	jvms.dirtyCache = true
	return nil
}

func (jvms *JvmsInfos) Save() error {
	for javaPath, jvmInfo := range jvms.Jvms {
		if value, found := jvms.fetched[javaPath]; !found || !value {
			if fileInfo, err := os.Stat(javaPath); err == nil {
				if fileInfo.ModTime().After(jvmInfo.FetchedAt) {
					if err := jvms.doFetch(javaPath); err != nil {
						return err
					}
				}
			} else {
				delete(jvms.Jvms, javaPath)
				jvms.dirtyCache = true
			}
		}
	}
	if jvms.dirtyCache {
		return writeToJson(jvms)
	}
	return nil
}

func writeToJson(jvmInfos *JvmsInfos) error {
	logDebug("Writing JVMs infos cache to %s", jvmInfos.path)
	file, err := json.MarshalIndent(jvmInfos, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(jvmInfos.path, file, 0644); err != nil {
		return wrapErr(err, "unable to write to file %s", jvmInfos.path)
	}
	return nil
}
