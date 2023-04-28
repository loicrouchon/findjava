package main

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	defaultJvmLookupPath := []string{
		"/bin/java",
		"/usr/bin/java",
		"/usr/local/bin/java",
		"/usr/lib/jvm",
		"~/.sdkman/candidates/java",
	}
	defaultJvmVersionRange := &VersionRange{
		Min: 0,
		Max: 0,
	}
	data := map[string]ConfigEntry{
		"test-resources/missing-config.json": {
			JvmLookupPaths:  defaultJvmLookupPath,
			JvmVersionRange: defaultJvmVersionRange,
		},
		"test-resources/empty-config.json": {
			JvmLookupPaths:  defaultJvmLookupPath,
			JvmVersionRange: defaultJvmVersionRange,
		},
		"test-resources/path-lookup-config.json": {
			JvmLookupPaths: []string{
				"/usr/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			},
			JvmVersionRange: defaultJvmVersionRange,
		},
		"test-resources/min-jvm-version-config.json": {
			JvmLookupPaths: defaultJvmLookupPath,
			JvmVersionRange: &VersionRange{
				Min: 8,
			},
		},
		"test-resources/full-config.json": {
			JvmLookupPaths: []string{
				"/usr/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			},
			JvmVersionRange: &VersionRange{
				Min: 8,
				Max: 17,
			},
		},
		"test-resources/invalid-config.json": {
			JvmLookupPaths:  defaultJvmLookupPath,
			JvmVersionRange: defaultJvmVersionRange,
		},
	}
	for path, expectedConfigEntry := range data {
		actualConfig := loadConfig(path, defaultKey)
		actualJvmLookupPath := *actualConfig.jvmsLookupPaths()
		if !reflect.DeepEqual(actualJvmLookupPath, expectedConfigEntry.JvmLookupPaths) {
			t.Fatalf(`Expecting loadConfig("%s", "%s").jvmLookupPaths() == %v but was %v`,
				path, defaultKey, expectedConfigEntry.JvmLookupPaths, actualJvmLookupPath)
		}
		actualJvmVersionRange := *actualConfig.jvmVersionRange()
		if !reflect.DeepEqual(actualJvmVersionRange, *(expectedConfigEntry.JvmVersionRange)) {
			t.Fatalf(`Expecting loadConfig("%s", "%s").jvmVersionRange() == %v but was %v`,
				path, defaultKey, *expectedConfigEntry.JvmVersionRange, actualJvmVersionRange)
		}
	}
}

func TestLoadConfigWithOverrides(t *testing.T) {
	data := map[string]ConfigEntry{
		"abc": {
			JvmLookupPaths: []string{
				"~/.sdkman/candidates/java",
			},
			JvmVersionRange: &VersionRange{
				Min: 8,
				Max: 17,
			},
		},
		"xyz": {
			JvmLookupPaths: []string{
				"/usr/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			},
			JvmVersionRange: &VersionRange{
				Min: 11,
			},
		},
	}
	for key, expectedConfigEntry := range data {
		path := "test-resources/full-config.json"
		actualConfig := loadConfig(path, key)
		actualJvmLookupPath := *actualConfig.jvmsLookupPaths()
		if !reflect.DeepEqual(actualJvmLookupPath, expectedConfigEntry.JvmLookupPaths) {
			t.Fatalf(`Expecting loadConfig("%s", "%s").jvmLookupPaths() == %v but was %v`,
				path, defaultKey, expectedConfigEntry.JvmLookupPaths, actualJvmLookupPath)
		}
		actualJvmVersionRange := *actualConfig.jvmVersionRange()
		if !reflect.DeepEqual(actualJvmVersionRange, *(expectedConfigEntry.JvmVersionRange)) {
			t.Fatalf(`Expecting loadConfig("%s", "%s").jvmVersionRange() == %v but was %v`,
				path, defaultKey, *expectedConfigEntry.JvmVersionRange, actualJvmVersionRange)
		}
	}
}
