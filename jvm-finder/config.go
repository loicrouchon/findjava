package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const defaultKey = "default"

// TODO add lookup in JAVA_HOME env var too
// TODO instead of a harcoded list, could also look for
//
//	all $DIR/java where $DIR are the individual $PATH entries
var defaultConfigEntry = ConfigEntry{
	path: "<DEFAULT>",
	JvmLookupPaths: []string{
		"/bin/java",
		"/usr/bin/java",
		"/usr/local/bin/java",
		"/usr/lib/jvm",
		"~/.sdkman/candidates/java",
	},
	JvmVersionRange: &VersionRange{
		Min: 0,
		Max: 0,
	},
}

type Config struct {
	configs []ConfigEntry
}

type ConfigEntry struct {
	path            string
	JvmLookupPaths  []string
	JvmVersionRange *VersionRange
}

const allVersions = 0

type VersionRange struct {
	Min uint
	Max uint
}

func (versionRange *VersionRange) String() string {
	return fmt.Sprintf("[%d..%d]}", versionRange.Min, versionRange.Max)
}

func (versionRange *VersionRange) Matches(version uint) bool {
	if versionRange.Min != allVersions && versionRange.Min > version {
		return false
	}
	if versionRange.Max != allVersions && versionRange.Max < version {
		return false
	}
	return true
}

func (config *Config) paths() []string {
	var paths []string
	for _, cfg := range config.configs {
		paths = append(paths, cfg.path)
	}
	return paths
}

func (config *Config) jvmLookupPaths() *[]string {
	for _, cfg := range config.configs {
		if len(cfg.JvmLookupPaths) > 0 {
			return &cfg.JvmLookupPaths
		}
	}
	die("no JVM lookup path defined in configuration files %v", config.paths())
	panic("unreachable")
}

func (config *Config) jvmVersionRange() *VersionRange {
	for _, cfg := range config.configs {
		if cfg.JvmVersionRange != nil {
			return cfg.JvmVersionRange
		}
	}
	die("no JVM Version range defined in configuration files %v", config.paths())
	panic("unreachable")
}

func loadConfig(defaultConfigPath string, name string) *Config {
	config := &Config{}
	var configPaths []string
	if name != defaultKey {
		specificConfigPath := strings.TrimSuffix(defaultConfigPath, ".json") + "." + name + ".json"
		configPaths = []string{specificConfigPath, defaultConfigPath}
	} else {
		configPaths = []string{defaultConfigPath}
	}
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			logDebug("Loading config from %s", path)
			configEntry := ConfigEntry{
				path: path,
			}
			file, _ := os.Open(path)
			defer file.Close()
			decoder := json.NewDecoder(file)
			err := decoder.Decode(&configEntry)
			if err != nil {
				dierr(err)
			}
			config.configs = append(config.configs, configEntry)
		} else {
			logDebug("Config file %s not found: %v", path, err)
		}
	}
	config.configs = append(config.configs, defaultConfigEntry)
	return config
}
