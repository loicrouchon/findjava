package config

import (
	"fmt"
	. "jvm-finder/internal/jvm"
	"jvm-finder/internal/utils"
	"jvm-finder/test"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	defaultJvmLookupPath := utils.ResolvePaths([]string{
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
			JvmLookupPaths: utils.ResolvePaths([]string{
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
			JvmLookupPaths: utils.ResolvePaths([]string{
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
		actual, err := LoadConfig(path, defaultKey)
		description := fmt.Sprintf("LoadConfig(\"%s\", \"%s\")", path, defaultKey)
		test.AssertNoError(t, description, err)
		test.AssertEquals(t, description+".jvmLookupPaths()", expected.JvmLookupPaths, actual.JvmsLookupPaths)
		test.AssertEquals(t, description+".JvmVersionRange()", *expected.JvmVersionRange, actual.JvmVersionRange)
	}
}

func TestLoadConfigWithOverrides(t *testing.T) {
	data := map[string]ConfigEntry{
		"abc": {
			JvmLookupPaths: utils.ResolvePaths([]string{
				"~/.sdkman/candidates/java",
			}),
			JvmVersionRange: &VersionRange{
				Min: 8,
				Max: 17,
			},
		},
		"xyz": {
			JvmLookupPaths: utils.ResolvePaths([]string{
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
		actual, err := LoadConfig(path, key)
		description := fmt.Sprintf("LoadConfig(\"%s\", \"%s\")", path, key)
		test.AssertNoError(t, description, err)
		test.AssertEquals(t, description+".jvmLookupPaths()", expected.JvmLookupPaths, actual.JvmsLookupPaths)
		test.AssertEquals(t, description+".JvmVersionRange()", *expected.JvmVersionRange, actual.JvmVersionRange)
	}
}
