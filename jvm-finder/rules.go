package main

import (
	"fmt"
	"regexp"
)

type JvmSelectionRules struct {
	minJvmVersion int
	maxJvmVersion int
}

func (rules JvmSelectionRules) Matches(jvmInfo JvmInfo) bool {
	if rules.minJvmVersion > 0 && rules.minJvmVersion > jvmInfo.javaSpecificationVersion {
		return false
	}
	if rules.maxJvmVersion > 0 && rules.maxJvmVersion < jvmInfo.javaSpecificationVersion {
		return false
	}
	return true
}

func (rules JvmSelectionRules) String() string {
	return fmt.Sprintf("[%d..%d]}", rules.minJvmVersion, rules.maxJvmVersion)
}

func jvmSelectionRules(args []string) *JvmSelectionRules {
	var rules *JvmSelectionRules
	if len(args) == 1 {
		jvmVersionRange := args[0]
		logDebug("%s", jvmVersionRange)
		jvmVersionRegex := `[\d]+(?:\.[\d]+)*`
		r := regexp.MustCompile(fmt.Sprintf(`^(?:`+
			`(?P<exact>%[1]s)`+
			`|(?:(?P<min>%[1]s)\.\.)`+
			`|(?:\.\.(?P<max>%[1]s))`+
			`|(?:(?P<min>%[1]s)\.\.(?P<max>%[1]s))`+
			`)$`,
			jvmVersionRegex))
		groupNames := r.SubexpNames()
		match := r.FindStringSubmatch(jvmVersionRange)
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
			minJvmVersion: minJvmVersion,
			maxJvmVersion: maxJvmVersion,
		}
	} else {
		rules = &JvmSelectionRules{}
	}
	logDebug("%s", rules)
	return rules
}
