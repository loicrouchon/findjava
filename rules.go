package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type JvmSelectionRules struct {
	versionRange   *VersionRange
	vendors        list
	programs       list
	preferredRules *JvmSelectionRules
}

func (rules *JvmSelectionRules) String() string {
	return fmt.Sprintf(`{
    versionRange: %v
    vendors: %v
    programs: %v
    preferredRules: %v
}`, rules.versionRange, rules.vendors, rules.programs, rules.preferredRules)
}

func (rules *JvmSelectionRules) Matches(jvm *Jvm) bool {
	if !rules.versionRange.Matches(jvm.javaSpecificationVersion) {
		return false
	}
	if !rules.matchVendor(jvm) {
		return false
	}
	if !rules.matchPrograms(jvm) {
		return false
	}
	return true
}

func (rules *JvmSelectionRules) matchVendor(jvm *Jvm) bool {
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

func (rules *JvmSelectionRules) matchPrograms(jvm *Jvm) bool {
	for _, program := range rules.programs {
		if program != "java" {
			programPath := filepath.Join(jvm.javaHome, "bin", program)
			if fileInfo, err := os.Stat(programPath); err == nil {
				if fileInfo.Mode()&0111 == 0 {
					logDebug("Program %s is not executable", programPath)
					return false
				}
			} else {
				logDebug("Program %s not found", programPath)
				return false
			}
		}
	}
	return true
}

func jvmSelectionRules(config *Config, minJavaVersion uint, maxJavaVersion uint, vendors list, programs list) *JvmSelectionRules {
	rules := &JvmSelectionRules{}
	rules.versionRange = &VersionRange{
		Min: minJavaVersion,
		Max: maxJavaVersion,
	}
	rules.vendors = vendors
	rules.programs = programs
	rules.preferredRules = &JvmSelectionRules{
		versionRange: &config.jvmVersionRange,
	}
	//logDebug("Requested version range: %v, preferred one: %v", rules.versionRange, rules.preferredVersionRange)
	logDebug("Resolved matching rules %v", rules)
	return rules
}
