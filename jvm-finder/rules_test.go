package main

import (
	"reflect"
	"testing"
)

func TestJvmSelectionRules(t *testing.T) {
	type TestData struct {
		minJavaVersion, maxJavaVersion uint
	}
	config := Config{
		jvmVersionRange: VersionRange{Min: 11, Max: allVersions},
	}
	versionRangesToSelectionRules := map[TestData]JvmSelectionRules{
		{minJavaVersion: 8, maxJavaVersion: 8}:                     {versionRange: &VersionRange{Min: 8, Max: 8}},
		{minJavaVersion: 17, maxJavaVersion: allVersions}:          {versionRange: &VersionRange{Min: 17, Max: allVersions}},
		{minJavaVersion: allVersions, maxJavaVersion: 11}:          {versionRange: &VersionRange{Min: allVersions, Max: 11}},
		{minJavaVersion: 9, maxJavaVersion: 14}:                    {versionRange: &VersionRange{Min: 9, Max: 14}},
		{minJavaVersion: allVersions, maxJavaVersion: allVersions}: {versionRange: &VersionRange{Min: 11, Max: allVersions}},
	}
	for versionRange, expectedRules := range versionRangesToSelectionRules {
		rules := jvmSelectionRules(versionRange.minJavaVersion, versionRange.maxJavaVersion, &config)
		if !reflect.DeepEqual(rules, &expectedRules) {
			t.Fatalf(`Expecting jvmSelectionRules("%v") == %v but was %v`,
				versionRange, expectedRules, rules)
		}
	}
}

func TestJvmSelectionRulesMatches(t *testing.T) {
	type TestData struct {
		rules       JvmSelectionRules
		jvmInfo     JvmInfo
		shouldMatch bool
	}
	testData := []TestData{
		// Exact version match
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 8, Max: 8}},
			jvmInfo:     jvmWithVersion(7),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 8, Max: 8}},
			jvmInfo:     jvmWithVersion(8),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 8, Max: 8}},
			jvmInfo:     jvmWithVersion(9),
			shouldMatch: false,
		},
		// Exact or next versions match
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 17, Max: 0}},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 17, Max: 0}},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 17, Max: 0}},
			jvmInfo:     jvmWithVersion(18),
			shouldMatch: true,
		},
		// Exact or previous versions match
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 0, Max: 17}},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 0, Max: 17}},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 0, Max: 17}},
			jvmInfo:     jvmWithVersion(18),
			shouldMatch: false,
		},
		// Full range match
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(10),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(11),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{versionRange: &VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(18),
			shouldMatch: false,
		},
	}
	for _, data := range testData {
		matches := data.rules.Matches(&data.jvmInfo)
		if matches != data.shouldMatch {
			t.Fatalf(`Expecting rules(%v).Matches("%v") == %t but was %t`,
				data.rules, data.jvmInfo, data.shouldMatch, matches)
		}
	}
}

func jvmWithVersion(version uint) JvmInfo {
	return JvmInfo{
		javaHome:                 "/jvm",
		javaSpecificationVersion: version,
	}
}
