package config

import (
	"findjvm/internal/log"
	"findjvm/internal/utils"
	"fmt"
	"os"
	"path/filepath"
)

type Platform struct {
	SelfPath             string
	ConfigDir            string
	CacheDir             string
	MetadataExtractorDir string
}

func (p *Platform) String() string {
	return fmt.Sprintf(`platform:
	program:                        %s
	config directory:               %s
	cache directory:                %s
	metadata extractor directory:   %s`, p.SelfPath, p.ConfigDir, p.CacheDir, p.MetadataExtractorDir)
}

func (p *Platform) LoadConfig(key string) (*Config, error) {
	err := p.setSelfPath()
	if err != nil {
		return nil, err
	}
	return loadConfig(filepath.Join(p.ConfigDir, "config.conf"), key, p.CacheDir, p.MetadataExtractorDir)
}

func (p *Platform) setSelfPath() error {
	self, err := os.Executable()
	if err != nil {
		return log.WrapErr(err, "unable to resolve findjvm self location")
	}
	self, err = filepath.EvalSymlinks(self)
	if err != nil {
		return log.WrapErr(err, "unable to resolve findjvm self location")
	}
	self, err = filepath.Abs(self)
	if err != nil {
		return log.WrapErr(err, "unable to resolve findjvm self location as an absolute path")
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
