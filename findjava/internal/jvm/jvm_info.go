package jvm

import (
	"encoding/json"
	"findjava/internal/discovery"
	"findjava/internal/log"
	"findjava/internal/utils"
	"os"
	"path/filepath"
	"time"
)

// JvmsInfos is a container type for all detected JVMs
type JvmsInfos struct {
	// The absolute path of the cache file
	cachePath string
	// True if the cache is dirty and needs to be updated
	dirtyCache bool
	// A map where keys are java executable absolute paths and values a boolean indicating whether their metadata have
	// been fetched or not.
	fetched map[string]bool
	// A map where keys are java executable absolute paths and values the associated [Jvm] metadata.
	Jvms map[string]*Jvm
}

// LoadJvmsInfos returns a [JvmsInfos] object providing information of the different JVMs denoted by the
// `javaExecutables` variable. The information of a given JVM are extracted thanks to the [MetadataReader] and
// cached on disk (cache location: `cachePath`).
//
// Cache entries are automatically updated in case the following cases:
//   - A new JVM not present in the cache is discovered -> metadata will be extracted and cached.
//   - A known (to the cache) Java executable has a file modification time more recent than the cache entry's timestamp
//     -> metadata will be extracted again and cached.
//   - A JVM entry in the cache does not exist anymore on disk -> the cache entry is deleted (cleanup).
//     Note that it might happen findjava is called with a configuration which only considers a sub-set of the cached
//     JVMs. In this case, those entries won't be evicted from the cache, unless the file on disk has been deleted.
//     This ensures cache entries are not aggressively removed from the cache when alternating calls to findjava with
//     configurations referring to different `jvm.lookup.paths`.
func LoadJvmsInfos(metadataReader *MetadataReader, cachePath string, javaExecutables *discovery.JavaExecutables) (JvmsInfos, error) {
	jvmInfos := loadJvmsInfosFromCache(cachePath)
	for javaPath, modTime := range javaExecutables.JavaPaths {
		if err := jvmInfos.Fetch(metadataReader, javaPath, modTime); err != nil {
			return JvmsInfos{}, err
		}
	}
	jvmInfos.evictOutdatedEntriesFromCache()
	if err := jvmInfos.Save(); err != nil {
		log.Warn(err)
	}
	return jvmInfos, nil
}

func loadJvmsInfosFromCache(path string) JvmsInfos {
	jvmsInfos := JvmsInfos{
		cachePath:  path,
		dirtyCache: false,
		fetched:    make(map[string]bool),
		Jvms:       make(map[string]*Jvm),
	}
	// Failures to load will from cache will result in an empty JvmsInfos
	// which will cause every discovered JVM to be fetched
	if _, err := os.Stat(path); err == nil {
		log.Debug("Loading cache from %s", path)
		if file, err := os.Open(path); err == nil {
			defer utils.CloseFile(file)
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&jvmsInfos); err == nil {
				for javaPath, jvm := range jvmsInfos.Jvms {
					jvm.javaPath = javaPath
					if err := jvm.rebuild(); err != nil {
						log.Warn(log.WrapErr(err, "cannot parse java specification version for JVM %s:", path))
						delete(jvmsInfos.Jvms, javaPath)
					}
				}
				//log.Debug("JVMs rebuilt loaded from cache: %#v", jvmsInfos)
			} else {
				log.Warn(log.WrapErr(err, "cannot parse cache file %s:", path))
				jvmsInfos.dirtyCache = true
			}
		} else {
			log.Warn(log.WrapErr(err, "cannot open cache file %s:", path))
			jvmsInfos.dirtyCache = true
		}
	}
	return jvmsInfos
}

func (jvms *JvmsInfos) Fetch(metadataReader *MetadataReader, javaPath string, modTime time.Time) error {
	jvms.fetched[javaPath] = true
	if info, found := jvms.Jvms[javaPath]; !found {
		log.Info("[CACHE MISS] %s", javaPath)
		return jvms.doFetch(metadataReader, javaPath)
	} else if modTime.After(info.FetchedAt) {
		log.Info("[CACHE OUTDATED] %s", javaPath)
		return jvms.doFetch(metadataReader, javaPath)
	} else {
		return nil
	}
}

func (jvms *JvmsInfos) doFetch(metadataReader *MetadataReader, javaPath string) error {
	jvm, err := metadataReader.fetchJvmInfo(javaPath)
	if err != nil {
		return err
	}
	log.Debug("%s:\n%s", javaPath, jvm)
	jvms.Jvms[javaPath] = jvm
	jvms.dirtyCache = true
	return nil
}

func (jvms *JvmsInfos) evictOutdatedEntriesFromCache() {
	for javaPath, jvmInfo := range jvms.Jvms {
		if value, found := jvms.fetched[javaPath]; !found || !value {
			if fileInfo, err := os.Stat(javaPath); err != nil || fileInfo.ModTime().After(jvmInfo.FetchedAt) {
				delete(jvms.Jvms, javaPath)
				log.Debug("evicting cache entry for JVM %s", javaPath)
				jvms.dirtyCache = true
			}
		}
	}
}

func (jvms *JvmsInfos) Save() error {
	if jvms.dirtyCache {
		return writeToJson(jvms)
	}
	return nil
}

func writeToJson(jvmInfos *JvmsInfos) error {
	log.Debug("Writing JVMs infos cache to %s", jvmInfos.cachePath)
	file, err := json.MarshalIndent(jvmInfos, "", "  ")
	if err != nil {
		return err
	}
	if err := utils.CreateDirectory(filepath.Dir(jvmInfos.cachePath)); err != nil {
		return log.WrapErr(err, "unable to create directory to host cache %s", jvmInfos.cachePath)
	}
	if err := utils.WriteFile(jvmInfos.cachePath, file, 0644); err != nil {
		return log.WrapErr(err, "unable to write to file %s", jvmInfos.cachePath)
	}
	return nil
}
