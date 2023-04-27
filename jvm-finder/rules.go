package main

import (
	"fmt"
)

type JvmSelectionRules struct {
	versionRange *VersionRange
}

func (rules *JvmSelectionRules) Matches(jvmInfo *JvmInfo) bool {
	return rules.versionRange.Matches(jvmInfo.javaSpecificationVersion)
}

func (rules JvmSelectionRules) String() string {
	return fmt.Sprintf("%v", rules.versionRange)
}

func jvmSelectionRules(minJavaVersion uint, maxJavaVersion uint, config *Config) *JvmSelectionRules {
	var rules JvmSelectionRules
	if minJavaVersion != allVersions || maxJavaVersion != allVersions {
		logDebug("Version range argument: [%d..%d], config: %v", minJavaVersion, maxJavaVersion, config.jvmVersionRange())
		rules = JvmSelectionRules{
			versionRange: &VersionRange{
				Min: minJavaVersion,
				Max: maxJavaVersion,
			},
		}
	} else {
		rules = JvmSelectionRules{
			versionRange: config.jvmVersionRange(),
		}
	}
	logDebug("Resolved matching rules %v", rules)
	return &rules
}
