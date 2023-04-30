package main

import (
	"fmt"
)

type JvmSelectionRules struct {
	versionRange *VersionRange
	vendors      list
}

func (rules *JvmSelectionRules) Matches(jvm *Jvm) bool {
	if !rules.versionRange.Matches(jvm.javaSpecificationVersion) {
		return false
	}
	if len(rules.vendors) > 0 {
		for _, vendor := range rules.vendors {
			if jvm.javaVendor == vendor {
				return true
			}
		}
		return false
	}
	return true
}

func (rules *JvmSelectionRules) String() string {
	return fmt.Sprintf(`{
    versionRange: %v
    vendors: %v
}`, rules.versionRange, rules.vendors)
}

func jvmSelectionRules(config *Config, minJavaVersion uint, maxJavaVersion uint, vendors list) *JvmSelectionRules {
	rules := &JvmSelectionRules{}
	if minJavaVersion != allVersions || maxJavaVersion != allVersions {
		logDebug("Version range argument: [%d..%d], config: %v", minJavaVersion, maxJavaVersion, config.jvmVersionRange)
		rules.versionRange = &VersionRange{
			Min: minJavaVersion,
			Max: maxJavaVersion,
		}
	} else {
		rules.versionRange = &config.jvmVersionRange
	}
	rules.vendors = vendors
	logDebug("Resolved matching rules %v", rules)
	return rules
}
