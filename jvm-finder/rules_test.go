package main

import (
	"testing"
)

func TestJvmSelectionRules(t *testing.T) {
	versionRangesToSelectionRules := map[string]JvmSelectionRules{
		"8":     {minJvmVersion: 8, maxJvmVersion: 8},
		"17..":  {minJvmVersion: 17, maxJvmVersion: 0},
		"..11":  {minJvmVersion: 0, maxJvmVersion: 11},
		"9..14": {minJvmVersion: 9, maxJvmVersion: 14},
	}
	for versionRange, expectedRules := range versionRangesToSelectionRules {
		rules := jvmSelectionRules(&versionRange)
		if rules.minJvmVersion != expectedRules.minJvmVersion {
			t.Fatalf(`Expecting jvmSelectionRules("%s") == %v but was %v`,
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
			rules:       JvmSelectionRules{minJvmVersion: 8, maxJvmVersion: 8},
			jvmInfo:     jvmWithVersion(7),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{minJvmVersion: 8, maxJvmVersion: 8},
			jvmInfo:     jvmWithVersion(8),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{minJvmVersion: 8, maxJvmVersion: 8},
			jvmInfo:     jvmWithVersion(9),
			shouldMatch: false,
		},
		// Exact or next versions match
		{
			rules:       JvmSelectionRules{minJvmVersion: 17, maxJvmVersion: 0},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{minJvmVersion: 17, maxJvmVersion: 0},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{minJvmVersion: 17, maxJvmVersion: 0},
			jvmInfo:     jvmWithVersion(18),
			shouldMatch: true,
		},
		// Exact or previous versions match
		{
			rules:       JvmSelectionRules{minJvmVersion: 0, maxJvmVersion: 17},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{minJvmVersion: 0, maxJvmVersion: 17},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{minJvmVersion: 0, maxJvmVersion: 17},
			jvmInfo:     jvmWithVersion(18),
			shouldMatch: false,
		},
		// Full range match
		{
			rules:       JvmSelectionRules{minJvmVersion: 11, maxJvmVersion: 17},
			jvmInfo:     jvmWithVersion(10),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{minJvmVersion: 11, maxJvmVersion: 17},
			jvmInfo:     jvmWithVersion(11),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{minJvmVersion: 11, maxJvmVersion: 17},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{minJvmVersion: 11, maxJvmVersion: 17},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{minJvmVersion: 11, maxJvmVersion: 17},
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

func jvmWithVersion(version int) JvmInfo {
	return JvmInfo{
		javaPath:                 "/jvm/bin/java",
		javaHome:                 "/jvm",
		javaSpecificationVersion: version,
	}
}
