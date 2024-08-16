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
	// The version range in which the `java.specification.version` should be contained.
	VersionRange *VersionRange
	// A list of JVM vendors to filter on (no filtering if empty).
	Vendors utils.List
	// A list of executable binaries the JVM must provide in its `${java.home}/bin` folder. Defaults to `java`.
	Programs utils.List
	// The preferred rules coming from the configuration to be applied if possible.
	PreferredRules *JvmSelectionRules
}

func (rules *JvmSelectionRules) String() string {
	return fmt.Sprintf(`
    VersionRange: %v
    Vendors: %v
    Programs: %v
    PreferredRules: %v`, rules.VersionRange, rules.Vendors, rules.Programs, rules.PreferredRules)
}

// Matches returns `true` if [JvmSelectionRules.VersionRange], [JvmSelectionRules.Vendors], and
// [JvmSelectionRules.Programs] rules are fulfilled.
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

// SelectionRules builds a [JvmSelectionRules] by:
//
//   - considering selection criteria given on the command line as strong constraints that must be fulfilled
//   - considering configuration related constraints as soft constraints that should be fulfilled if possible.
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
