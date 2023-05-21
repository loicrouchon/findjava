package rules

import (
	"jvm-finder/internal/config"
	. "jvm-finder/internal/jvm"
	"reflect"
	"testing"
)

func TestJvmSelectionRules(t *testing.T) {
	type TestData struct {
		minJavaVersion, maxJavaVersion uint
	}
	config := config.Config{
		JvmVersionRange: VersionRange{Min: 11, Max: AllVersions},
	}
	preferredRules := &JvmSelectionRules{VersionRange: &config.JvmVersionRange}
	versionRangesToSelectionRules := map[TestData]JvmSelectionRules{
		{minJavaVersion: 8, maxJavaVersion: 8}: {
			VersionRange:   &VersionRange{Min: 8, Max: 8},
			PreferredRules: preferredRules,
		},
		{minJavaVersion: 17, maxJavaVersion: AllVersions}: {
			VersionRange:   &VersionRange{Min: 17, Max: AllVersions},
			PreferredRules: preferredRules,
		},
		{minJavaVersion: AllVersions, maxJavaVersion: 11}: {
			VersionRange:   &VersionRange{Min: AllVersions, Max: 11},
			PreferredRules: preferredRules,
		},
		{minJavaVersion: 9, maxJavaVersion: 14}: {
			VersionRange:   &VersionRange{Min: 9, Max: 14},
			PreferredRules: preferredRules,
		},
		{minJavaVersion: AllVersions, maxJavaVersion: AllVersions}: {
			VersionRange:   &VersionRange{Min: AllVersions, Max: AllVersions},
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
		jvmInfo     Jvm
		shouldMatch bool
	}
	testData := []TestData{
		// Exact version match
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 8, Max: 8}},
			jvmInfo:     jvmWithVersion(7),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 8, Max: 8}},
			jvmInfo:     jvmWithVersion(8),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 8, Max: 8}},
			jvmInfo:     jvmWithVersion(9),
			shouldMatch: false,
		},
		// Exact or next versions match
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 17, Max: 0}},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 17, Max: 0}},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 17, Max: 0}},
			jvmInfo:     jvmWithVersion(18),
			shouldMatch: true,
		},
		// Exact or previous versions match
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 0, Max: 17}},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 0, Max: 17}},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 0, Max: 17}},
			jvmInfo:     jvmWithVersion(18),
			shouldMatch: false,
		},
		// Full range match
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(10),
			shouldMatch: false,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(11),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(15),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 11, Max: 17}},
			jvmInfo:     jvmWithVersion(17),
			shouldMatch: true,
		},
		{
			rules:       JvmSelectionRules{VersionRange: &VersionRange{Min: 11, Max: 17}},
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

func jvmWithVersion(version uint) Jvm {
	return Jvm{
		JavaHome:                 "/jvm",
		JavaSpecificationVersion: version,
	}
}
