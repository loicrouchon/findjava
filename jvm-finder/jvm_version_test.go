package main

import (
	"strconv"
	"testing"
)

func TestParseVersion(t *testing.T) {
	versions := map[string]int{
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
		versions[strconv.Itoa(i)] = i
	}
	for versionToParse, expectedVersion := range versions {
		parsedVersion := parseVersion(versionToParse)
		if parsedVersion != expectedVersion {
			t.Fatalf(`Expecting parseVersion(%s) == %d but was %d`,
				versionToParse, expectedVersion, parsedVersion)
		}
	}
}
