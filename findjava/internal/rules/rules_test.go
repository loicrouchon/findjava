package rules

import (
	"findjava/internal/config"
	"findjava/internal/jvm"
	"reflect"
	"testing"
)

func TestJvmSelectionRules(t *testing.T) {
	type TestData struct {
		minJavaVersion, maxJavaVersion uint
	}
	config := config.Config{
		JvmVersionRange: jvm.VersionRange{Min: 11, Max: jvm.AllVersions},
	}
	preferredRules := &JvmSelectionRules{VersionRange: &config.JvmVersionRange}
	versionRangesToSelectionRules := map[TestData]JvmSelectionRules{
		{minJavaVersion: 8, maxJavaVersion: 8}: {
			VersionRange:   &jvm.VersionRange{Min: 8, Max: 8},
			PreferredRules: preferredRules,
		},
		{minJavaVersion: 17, maxJavaVersion: jvm.AllVersions}: {
			VersionRange:   &jvm.VersionRange{Min: 17, Max: jvm.AllVersions},
			PreferredRules: preferredRules,
		},
		{minJavaVersion: jvm.AllVersions, maxJavaVersion: 11}: {
			VersionRange:   &jvm.VersionRange{Min: jvm.AllVersions, Max: 11},
			PreferredRules: preferredRules,
		},
		{minJavaVersion: 9, maxJavaVersion: 14}: {
			VersionRange:   &jvm.VersionRange{Min: 9, Max: 14},
			PreferredRules: preferredRules,
		},
		{minJavaVersion: jvm.AllVersions, maxJavaVersion: jvm.AllVersions}: {
			VersionRange:   &jvm.VersionRange{Min: jvm.AllVersions, Max: jvm.AllVersions},
			PreferredRules: preferredRules,
		},
	}
	for versionRange, expectedRules := range versionRangesToSelectionRules {
		rules := SelectionRules(&config, versionRange.minJavaVersion, versionRange.maxJavaVersion, nil, nil)
		if !reflect.DeepEqual(rules, &expectedRules) {
			t.Fatalf(`Expecting SelectionRules("%v") == %v but was %v`,
				versionRange, &expectedRules, rules)
		}
	}
}

func TestJvmSelectionRulesMatches(t *testing.T) {
	type TestData struct {
		rules       JvmSelectionRules
		jvmInfo     jvm.Jvm
		shouldMatch bool
	}
	testData := []TestData{
		// Exact version match
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 8, Max: 8}},
			jvmInfo:     jvmWithVersion(7),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 8, Max: 8}},
			jvmInfo:     jvmWithVersion(8),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 8, Max: 8}},
			jvmInfo:     jvmWithVersion(9),
			shouldMatch: false,
		},
		// Exact or next versions match
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 17, Max: 0}},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 17, Max: 0}},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 17, Max: 0}},
			jvmInfo:     jvmWithVersion(18),
			shouldMatch: true,
		},
		// Exact or previous versions match
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 0, Max: 17}},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 0, Max: 17}},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 0, Max: 17}},
			jvmInfo:     jvmWithVersion(18),
			shouldMatch: false,
		},
		// Full range match
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(10),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(11),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &jvm.VersionRange{Min: 11, Max: 17}},
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

func jvmWithVersion(version uint) jvm.Jvm {
	return jvm.Jvm{
		JavaHome:                 "/jvm",
		JavaSpecificationVersion: version,
	}
}
