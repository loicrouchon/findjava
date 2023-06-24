package rules

import (
	"findjava/internal/config"
	. "findjava/internal/jvm"
	"findjava/internal/log"
	"findjava/internal/utils"
	"fmt"
	"os"
	"path/filepath"
)

type JvmSelectionRules struct {
	VersionRange   *VersionRange
	Vendors        utils.List
	Programs       utils.List
	PreferredRules *JvmSelectionRules
}

func (rules *JvmSelectionRules) String() string {
	return fmt.Sprintf(`
    VersionRange: %v
    Vendors: %v
    Programs: %v
    PreferredRules: %v`, rules.VersionRange, rules.Vendors, rules.Programs, rules.PreferredRules)
}

func (rules *JvmSelectionRules) Matches(jvm *Jvm) bool {
	if !rules.VersionRange.Matches(jvm.JavaSpecificationVersion) {
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
	if len(rules.Vendors) > 0 {
		for _, vendor := range rules.Vendors {
			if jvm.JavaVendor == vendor {
				return true
			}
		}
		return false
	}
	return true
}

func (rules *JvmSelectionRules) matchPrograms(jvm *Jvm) bool {
	for _, program := range rules.Programs {
		if program != "java" {
			programPath := filepath.Join(jvm.JavaHome, "bin", program)
			if fileInfo, err := os.Stat(programPath); err == nil {
				if fileInfo.Mode()&0111 == 0 {
					log.Debug("Program %s is not executable", programPath)
					return false
				}
			} else {
				log.Debug("Program %s not found", programPath)
				return false
			}
		}
	}
	return true
}

func SelectionRules(config *config.Config, minJavaVersion uint, maxJavaVersion uint, vendors utils.List, programs utils.List) *JvmSelectionRules {
	rules := &JvmSelectionRules{}
	rules.VersionRange = &VersionRange{
		Min: minJavaVersion,
		Max: maxJavaVersion,
	}
	rules.Vendors = vendors
	rules.Programs = programs
	rules.PreferredRules = &JvmSelectionRules{
		VersionRange: &config.JvmVersionRange,
	}
	//log.Debug("Requested version range: %v, preferred one: %v", rules.VersionRange, rules.preferredVersionRange)
	log.Debug("Resolved matching rules %v", rules)
	return rules
}
