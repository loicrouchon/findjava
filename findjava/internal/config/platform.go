package config

import (
	"findjava/internal/log"
	"findjava/internal/utils"
	"fmt"
	"os"
	"path/filepath"
)

// Platform contains information that are platform dependent.
type Platform struct {
	// SelfPath The path to the findjava binary.
	SelfPath string
	// ConfigDir The path to the directory holding the configuration of findjava.
	ConfigDir string
	// CacheDir The path to the directory in which JVMs metadata will be cached.
	CacheDir string
	// MetadataExtractorDir The path to the directory in which the JvmMetadataExtractor is located.
	MetadataExtractorDir string
}

func (p *Platform) String() string {
	return fmt.Sprintf(`platform:
	program:                        %s
	config directory:               %s
	cache directory:                %s
	metadata extractor directory:   %s`, p.SelfPath, p.ConfigDir, p.CacheDir, p.MetadataExtractorDir)
}

// LoadConfig loads the configuration for the given [Platform] into a [Config] object.
//
// Configurations will be loaded from three different sources:
//
//   - specific configuration file: `p.ConfigDir + "/config." + key + ".conf"`
//   - system configuration file:   `p.ConfigDir + "/config.conf"`
//   - default configuration:       constant [defaultConfigEntry]
//
// Those configurations sources are merged together into a single [Config] object.
// The merge algorithm works at property level. It resolves each property individually looking for it first in the
// specific configuration file (if it exists). If the property is not found, it then tries in the system configuration
// file. If still not found, the default value from the [defaultConfigEntry] will be used.
// This mechanism allows most specific configuration sources to override values from lower priority configuration
// sources without requiring to redefine the whole configuration.
func (p *Platform) LoadConfig(key string) (*Config, error) {
	err := p.Resolve()
	if err != nil {
		return nil, err
	}
	return loadConfig(filepath.Join(p.ConfigDir, "config.conf"), key, p.CacheDir, p.MetadataExtractorDir)
}

func (p *Platform) Resolve() error {
	self, err := os.Executable()
	if err != nil {
		return log.WrapErr(err, "unable to resolve findjava self location")
	}
	self, err = filepath.EvalSymlinks(self)
	if err != nil {
		return log.WrapErr(err, "unable to resolve findjava self location")
	}
	self, err = filepath.Abs(self)
	if err != nil {
		return log.WrapErr(err, "unable to resolve findjava self location as an absolute path")
	}
	p.SelfPath = self
	selfDir := filepath.Dir(p.SelfPath)
	p.ConfigDir, err = toAbsolutePath(selfDir, p.ConfigDir)
	if err != nil {
		return err
	}
	p.CacheDir, err = toAbsolutePath(selfDir, p.CacheDir)
	if err != nil {
		return err
	}
	p.MetadataExtractorDir, err = toAbsolutePath(selfDir, p.MetadataExtractorDir)
	if err != nil {
		return err
	}
	log.Debug("%v", p)
	return nil
}

func toAbsolutePath(self string, path string) (string, error) {
	path, err := utils.ResolvePath(path)
	if err != nil {
		return "", nil
	}
	if filepath.IsAbs(path) {
		return path, nil
	}
	return filepath.Join(self, path), nil
}
