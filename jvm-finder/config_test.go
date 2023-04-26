package main

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	data := map[string]Config{
		"test-resources/missing-config.json": {
			JvmLookupPaths: []string{
				"/bin/java",
				"/usr/bin/java",
				"/usr/local/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			},
		},
		"test-resources/empty-config.json": {},
		"test-resources/path-lookup-config.json": {
			JvmLookupPaths: []string{
				"/usr/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			},
		},
		"test-resources/min-jvm-version-config.json": {
			DefaultJvmVersionRange: VersionRange{
				Min: 8,
			},
		},
		"test-resources/full-config.json": {
			JvmLookupPaths: []string{
				"/bin/java",
				"/usr/bin/java",
				"/usr/local/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			},
			DefaultJvmVersionRange: VersionRange{
				Min: 8,
				Max: 17,
			},
		},
	}
	for path, expectedConfig := range data {
		actualConfig := loadConfig(path)
		if !reflect.DeepEqual(*actualConfig, expectedConfig) {
			t.Fatalf(`Expecting loadConfig("%s") == %v but was %v`,
				path, expectedConfig, *actualConfig)
		}
	}
}
