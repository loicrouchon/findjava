package config

import (
	"encoding/json"
	"fmt"
	. "jvm-finder/internal/jvm"
	"jvm-finder/internal/log"
	"jvm-finder/internal/utils"
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
		"$HOMEBREW_CELLAR/openjdk",
	},
	JvmVersionRange: &VersionRange{
		Min: 0,
		Max: 0,
	},
}

type Config struct {
	JvmsMetadataCachePath string
	JvmsLookupPaths       []string
	JvmVersionRange       VersionRange
}

func (cfg *Config) String() string {
	return fmt.Sprintf(`{
	JvmsMetadataCachePath: %s
	JvmLookupPaths: %v
	JvmVersionRange: %s
}`, cfg.JvmsMetadataCachePath, cfg.JvmsLookupPaths, &cfg.JvmVersionRange)
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

func LoadConfig(defaultConfigPath string, name string) (*Config, error) {
	var configs []ConfigEntry
	configPaths := configPaths(name, defaultConfigPath)
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			if configEntry, err := loadConfigFromFile(path); err != nil {
				return nil, err
			} else {
				configs = append(configs, configEntry)
			}
		} else {
			log.Debug("Config file %s not found: %v", path, err)
		}
	}
	configs = append(configs, defaultConfigEntry)
	return parseConfig(configs)
}

func parseConfig(configs []ConfigEntry) (*Config, error) {
	log.Debug("Config entries: %v", configs)
	path, err := jvmsMetadataCachePath(configs)
	if err != nil {
		return nil, err
	}
	lookupPaths, err := jvmsLookupPaths(configs)
	if err != nil {
		return nil, err
	}
	versionRange, err := jvmVersionRange(configs)
	if err != nil {
		return nil, err
	}
	config := Config{
		JvmsMetadataCachePath: path,
		JvmsLookupPaths:       lookupPaths,
		JvmVersionRange:       versionRange,
	}
	log.Debug("Resolved config: %s", &config)
	return &config, nil
}

func loadConfigFromFile(path string) (ConfigEntry, error) {
	log.Debug("Loading config from %s", path)
	configEntry := ConfigEntry{
		path: path,
	}
	file, _ := os.Open(path)
	defer utils.CloseFile(file)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configEntry)
	return configEntry, err
}

func configPaths(name string, defaultConfigPath string) []string {
	if name != defaultKey {
		specificConfigPath := strings.TrimSuffix(defaultConfigPath, ".json") + "." + name + ".json"
		return []string{specificConfigPath, defaultConfigPath}
	} else {
		return []string{defaultConfigPath}
	}
}

func jvmsMetadataCachePath(configs []ConfigEntry) (string, error) {
	for _, cfg := range configs {
		if cfg.JvmsMetadataCachePath != "" {
			return utils.ResolvePath(cfg.JvmsMetadataCachePath)
		}
	}
	return "", fmt.Errorf("no JVMs metadata cache path defined in configuration files %v\n", paths(configs))
}

func jvmsLookupPaths(configs []ConfigEntry) ([]string, error) {
	for _, cfg := range configs {
		if len(cfg.JvmLookupPaths) > 0 {
			resolvedPaths := utils.ResolvePaths(cfg.JvmLookupPaths)
			if len(resolvedPaths) > 0 {
				return resolvedPaths, nil
			}
		}
	}
	return nil, fmt.Errorf("no JVMs lookup path defined in configuration files %v\n", paths(configs))
}

func jvmVersionRange(configs []ConfigEntry) (VersionRange, error) {
	for _, cfg := range configs {
		if cfg.JvmVersionRange != nil {
			return *cfg.JvmVersionRange, nil
		}
	}
	return VersionRange{}, fmt.Errorf("no version range defined in configuration files %v\n", paths(configs))
}

func paths(configs []ConfigEntry) []string {
	var paths []string
	for _, cfg := range configs {
		paths = append(paths, cfg.path)
	}
	return paths
}
