package main

import (
	"fmt"
	"regexp"
)

var jvmVersionRegex = `[\d]+(?:\.[\d]+)*`
var r = regexp.MustCompile(fmt.Sprintf(`^(?:`+
	`(?P<exact>%[1]s)`+
	`|(?:(?P<min>%[1]s)\.\.)`+
	`|(?:\.\.(?P<max>%[1]s))`+
	`|(?:(?P<min>%[1]s)\.\.(?P<max>%[1]s))`+
	`)$`,
	jvmVersionRegex))
var groupNames = r.SubexpNames()

type JvmSelectionRules struct {
	versionRange *VersionRange
}

func (rules *JvmSelectionRules) Matches(jvmInfo *JvmInfo) bool {
	return rules.versionRange.Matches(jvmInfo.javaSpecificationVersion)
}

func (rules JvmSelectionRules) String() string {
	return fmt.Sprintf("%v", rules.versionRange)
}

func jvmSelectionRules(jvmVersionRange *string, config *Config) *JvmSelectionRules {
	var rules *JvmSelectionRules
	if jvmVersionRange != nil && len(*jvmVersionRange) > 0 {
		logDebug("Version range argument: %s, config: %v", *jvmVersionRange, config.jvmVersionRange())
		match := r.FindStringSubmatch(*jvmVersionRange)
		if len(match) <= 0 {
			return nil
		}
		var minJvmVersion int
		var maxJvmVersion int
		for i, m := range match {
			if len(m) > 0 {
				switch groupNames[i] {
				case "exact":
					minJvmVersion = parseVersion(m)
					maxJvmVersion = parseVersion(m)
				case "min":
					minJvmVersion = parseVersion(m)
				case "max":
					maxJvmVersion = parseVersion(m)
				}
			}
		}
		rules = &JvmSelectionRules{
			versionRange: &VersionRange{
				Min: minJvmVersion,
				Max: maxJvmVersion,
			},
		}
	} else {
		rules = &JvmSelectionRules{
			versionRange: config.jvmVersionRange(),
		}
	}
	logDebug("Resolved matching rules %v", rules)
	return rules
}
