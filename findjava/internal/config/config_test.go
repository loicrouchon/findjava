package config

import (
	"findjava/internal/jvm"
	"findjava/internal/utils"
	"findjava/test"
	"fmt"
	"testing"
)

func TestLoadInvalidConfig(t *testing.T) {
	data := map[string][]string{
		"test-resources/invalid.conf": {
			"invalid configuration entry in file test-resources/invalid.conf for key 'non.existing.property' and value 'value'",
			"unknown key 'non.existing.property'",
		},
		"test-resources/invalid-min-java-version.conf": {
			"invalid configuration entry in file test-resources/invalid-min-java-version.conf for key 'java.specification.version.min' and value '-1'",
			"JVM version '-1' cannot be parsed as an unsigned int",
		},
		"test-resources/invalid-max-java-version.conf": {
			"invalid configuration entry in file test-resources/invalid-max-java-version.conf for key 'java.specification.version.max' and value 'this is obviously invalid'",
			"JVM version 'this is obviously invalid' cannot be parsed as an unsigned int",
		},
	}
	for path, expected := range data {
		_, err := loadConfig(path, defaultKey, "", "")
		description := fmt.Sprintf("loadConfig(\"%s\", \"%s\")", path, defaultKey)
		for _, e := range expected {
			test.AssertErrorContains(t, description, e, err)
		}
	}
}

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
	defaultJvmVersionRange := &jvm.VersionRange{
		Min: 0,
		Max: 0,
	}
	data := map[string]configEntry{
		"test-resources/missing.conf": {
			JvmLookupPaths:  defaultJvmLookupPath,
			JvmVersionRange: defaultJvmVersionRange,
		},
		"test-resources/empty.conf": {
			JvmLookupPaths:  defaultJvmLookupPath,
			JvmVersionRange: defaultJvmVersionRange,
		},
		"test-resources/path-lookup.conf": {
			JvmLookupPaths: utils.ResolvePaths([]string{
				"/usr/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			}),
			JvmVersionRange: defaultJvmVersionRange,
		},
		"test-resources/min-jvm-version.conf": {
			JvmLookupPaths: defaultJvmLookupPath,
			JvmVersionRange: &jvm.VersionRange{
				Min: 8,
			},
		},
		"test-resources/full.conf": {
			JvmLookupPaths: utils.ResolvePaths([]string{
				"/usr/bin/java",
				"/usr/lib/jvm",
				"~/.sdkman/candidates/java",
			}),
			JvmVersionRange: &jvm.VersionRange{
				Min: 8,
				Max: 17,
			},
		},
	}
	for path, expected := range data {
		actual, err := loadConfig(path, defaultKey, "", "")
		description := fmt.Sprintf("loadConfig(\"%s\", \"%s\")", path, defaultKey)
		test.AssertNoError(t, description, err)
		test.AssertEquals(t, description+".jvmLookupPaths()", expected.JvmLookupPaths, actual.JvmsLookupPaths)
		test.AssertEquals(t, description+".JvmVersionRange()", *expected.JvmVersionRange, actual.JvmVersionRange)
	}
}

func TestLoadConfigWithOverrides(t *testing.T) {
	data := map[string]configEntry{
		"abc": {
			JvmLookupPaths: utils.ResolvePaths([]string{
				"~/.sdkman/candidates/java",
			}),
			JvmVersionRange: &jvm.VersionRange{
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
			JvmVersionRange: &jvm.VersionRange{
				Max: 11,
			},
		},
	}
	for key, expected := range data {
		path := "test-resources/full.conf"
		actual, err := loadConfig(path, key, "", "")
		description := fmt.Sprintf("loadConfig(\"%s\", \"%s\")", path, key)
		test.AssertNoError(t, description, err)
		test.AssertEquals(t, description+".jvmLookupPaths()", expected.JvmLookupPaths, actual.JvmsLookupPaths)
		test.AssertEquals(t, description+".JvmVersionRange()", *expected.JvmVersionRange, actual.JvmVersionRange)
	}
}
