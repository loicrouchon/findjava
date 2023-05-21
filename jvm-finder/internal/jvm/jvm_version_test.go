package jvm

import (
	"fmt"
	"jvm-finder/test"
	"strconv"
	"testing"
)

func TestParseVersion(t *testing.T) {
	versions := map[string]uint{
		"1.0": 1,
		"1.1": 1,
		"1.2": 2,
		"1.3": 3,
		"1.4": 4,
		"1.5": 5,
		"1.6": 6,
		"1.7": 7,
		"1.8": 8,
	}
	for i := 1; i < 25; i++ {
		versions[strconv.Itoa(i)] = uint(i)
	}
	for versionToParse, expected := range versions {
		actual, err := parseJavaSpecificationVersion(versionToParse)
		description := fmt.Sprintf("parseJavaSpecificationVersion(%s)", versionToParse)
		test.AssertNoError(t, description, err)
		test.AssertEquals(t, description, expected, actual)
	}
}

func TestParseVersionError(t *testing.T) {
	versions := map[string]string{
		"":    "JVM version '' cannot be parsed as an unsigned int",
		"-1":  "JVM version '-1' cannot be parsed as an unsigned int",
		"one": "JVM version 'one' cannot be parsed as an unsigned int",
	}
	for versionToParse, expected := range versions {
		_, err := parseJavaSpecificationVersion(versionToParse)
		description := fmt.Sprintf("parseJavaSpecificationVersion(%s)", versionToParse)
		test.AssertErrorContains(t, description, expected, err)
	}
}
