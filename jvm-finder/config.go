package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	JvmLookupPaths         []string
	DefaultJvmVersionRange VersionRange
}

type VersionRange struct {
	Min int
	Max int
}

func loadConfig(path string) *Config {
	var config *Config
	if _, err := os.Stat(path); err == nil {
		file, _ := os.Open(path)
		defer file.Close()
		decoder := json.NewDecoder(file)
		config = &Config{}
		err := decoder.Decode(config)
		if err != nil {
			dierr(err)
		}
	} else {
		config = &Config{
			JvmLookupPaths: []string{
				"/bin/java",
				"/usr/bin/java",
				"/usr/local/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			},
		}
	}
	return config
}
