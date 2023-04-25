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
