package main

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	defaultJvmLookupPath := resolvePaths([]string{
		"$JAVA_HOME/bin/java",
		"$GRAALVM_HOME/bin/java",
		"/bin/java",
		"/usr/bin/java",
		"/usr/local/bin/java",
		"/usr/lib/jvm",
		"~/.sdkman/candidates/java",
		"$HOMEBREW_CELLAR/openjdk",
	})
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
			JvmLookupPaths: resolvePaths([]string{
				"/usr/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			}),
			JvmVersionRange: defaultJvmVersionRange,
		},
		"test-resources/min-jvm-version-config.json": {
			JvmLookupPaths: defaultJvmLookupPath,
			JvmVersionRange: &VersionRange{
				Min: 8,
			},
		},
		"test-resources/full-config.json": {
			JvmLookupPaths: resolvePaths([]string{
				"/usr/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			}),
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
	for path, expected := range data {
		actual, err := loadConfig(path, defaultKey)
		description := fmt.Sprintf("loadConfig(\"%s\", \"%s\")", path, defaultKey)
		assertNoError(t, description, err)
		assertEquals(t, description+".jvmLookupPaths()", expected.JvmLookupPaths, actual.jvmsLookupPaths)
		assertEquals(t, description+".jvmVersionRange()", *expected.JvmVersionRange, actual.jvmVersionRange)
	}
}

func TestLoadConfigWithOverrides(t *testing.T) {
	data := map[string]ConfigEntry{
		"abc": {
			JvmLookupPaths: resolvePaths([]string{
				"~/.sdkman/candidates/java",
			}),
			JvmVersionRange: &VersionRange{
				Min: 8,
				Max: 17,
			},
		},
		"xyz": {
			JvmLookupPaths: resolvePaths([]string{
				"/usr/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			}),
			JvmVersionRange: &VersionRange{
				Min: 11,
			},
		},
	}
	for key, expected := range data {
		path := "test-resources/full-config.json"
		actual, err := loadConfig(path, key)
		description := fmt.Sprintf("loadConfig(\"%s\", \"%s\")", path, key)
		assertNoError(t, description, err)
		assertEquals(t, description+".jvmLookupPaths()", expected.JvmLookupPaths, actual.jvmsLookupPaths)
		assertEquals(t, description+".jvmVersionRange()", *expected.JvmVersionRange, actual.jvmVersionRange)
	}
}
