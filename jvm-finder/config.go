package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const defaultKey = ""

var defaultConfigEntry = ConfigEntry{
	path:                  "<DEFAULT>",
	JvmsMetadataCachePath: "./build/jvm-finder.json",
	JvmLookupPaths: []string{
		"$JAVA_HOME/bin/java",
		"$GRAALVM_HOME/bin/java",
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
	jvmsMetadataCachePath string
	jvmsLookupPaths       []string
	jvmVersionRange       VersionRange
}

func (cfg *Config) String() string {
	return fmt.Sprintf(`{
	JvmsMetadataCachePath: %s
	JvmLookupPaths: %v
	JvmVersionRange: %s
}`, cfg.jvmsMetadataCachePath, cfg.jvmsLookupPaths, &cfg.jvmVersionRange)
}

type ConfigEntry struct {
	path                  string
	JvmsMetadataCachePath string
	JvmLookupPaths        []string
	JvmVersionRange       *VersionRange
}

func (cfg ConfigEntry) String() string {
	return fmt.Sprintf(`{
	path: %s
	JvmsMetadataCachePath: %s
	JvmLookupPaths: %v
	JvmVersionRange: %s
}`, cfg.path, cfg.JvmsMetadataCachePath, cfg.JvmLookupPaths, cfg.JvmVersionRange)
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

func loadConfig(defaultConfigPath string, name string) *Config {
	var configs []ConfigEntry
	configPaths := configPaths(name, defaultConfigPath)
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
			configs = append(configs, configEntry)
		} else {
			logDebug("Config file %s not found: %v", path, err)
		}
	}
	configs = append(configs, defaultConfigEntry)
	logDebug("Config entries: %v", configs)
	config := Config{
		jvmsMetadataCachePath: jvmsMetadataCachePath(configs),
		jvmsLookupPaths:       jvmsLookupPaths(configs),
		jvmVersionRange:       jvmVersionRange(configs),
	}
	logDebug("Resolved config: %s", &config)
	return &config
}

func configPaths(name string, defaultConfigPath string) []string {
	if name != defaultKey {
		specificConfigPath := strings.TrimSuffix(defaultConfigPath, ".json") + "." + name + ".json"
		return []string{specificConfigPath, defaultConfigPath}
	} else {
		return []string{defaultConfigPath}
	}
}

func jvmsMetadataCachePath(configs []ConfigEntry) string {
	for _, cfg := range configs {
		if cfg.JvmsMetadataCachePath != "" {
			return resolvePath(cfg.JvmsMetadataCachePath)
		}
	}
	die("no JVMs metadata cache path defined in configuration files %v", paths(configs))
	panic("unreachable")
}

func jvmsLookupPaths(configs []ConfigEntry) []string {
	for _, cfg := range configs {
		if len(cfg.JvmLookupPaths) > 0 {
			resolvedPaths := resolvePaths(cfg.JvmLookupPaths)
			if len(resolvedPaths) > 0 {
				return resolvedPaths
			}
		}
	}
	die("no JVMs lookup path defined in configuration files %v", paths(configs))
	panic("unreachable")
}

func jvmVersionRange(configs []ConfigEntry) VersionRange {
	for _, cfg := range configs {
		if cfg.JvmVersionRange != nil {
			return *cfg.JvmVersionRange
		}
	}
	die("no version range defined in configuration files %v", paths(configs))
	panic("unreachable")
}

func paths(configs []ConfigEntry) []string {
	var paths []string
	for _, cfg := range configs {
		paths = append(paths, cfg.path)
	}
	return paths
}
