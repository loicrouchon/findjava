package main

import (
	"sort"
)

func (jvms *JvmsInfos) Select(rules *JvmSelectionRules) []Jvm {
	candidates, ignored := filter(rules, jvms)
	sort.Slice(ignored[:], func(i, j int) bool { return sortCandidates(ignored, i, j) })
	sort.Slice(candidates[:], func(i, j int) bool { return sortCandidates(candidates, i, j) })
	logJvmList("[IGNORED]", ignored)
	logJvmList("[CANDIDATE]", candidates)
	return candidates
}

func filter(rules *JvmSelectionRules, jvms *JvmsInfos) ([]Jvm, []Jvm) {
	var allJvms []Jvm
	for _, jvm := range jvms.Jvms {
		allJvms = append(allJvms, *jvm)
	}
	return filterJvmList(rules, allJvms)
}

func filterJvmList(rules *JvmSelectionRules, allJvms []Jvm) ([]Jvm, []Jvm) {
	var candidates []Jvm
	var ignored []Jvm
	for _, jvm := range allJvms {
		if rules.Matches(&jvm) {
			candidates = append(candidates, jvm)
		} else {
			ignored = append(ignored, jvm)
		}
	}
	if len(candidates) > 0 && rules.preferredRules != nil {
		preferredCandidates, preferredIgnored := filterJvmList(rules.preferredRules, candidates)
		if len(preferredCandidates) > 0 {
			return preferredCandidates, preferredIgnored
		} else if !rules.versionRange.isBounded() {
			return nil, preferredIgnored
		} else {
			logInfo("Unable to satisfy preferred selection rules %v, ignoring them", rules.preferredRules)
		}
	}
	return candidates, ignored
}

func logJvmList(displayType string, jvms []Jvm) {
	for i := len(jvms) - 1; i >= 0; i = i - 1 {
		jvm := jvms[i]
		logInfo("%-12s %3d: %s ", displayType, jvm.javaSpecificationVersion, jvm.javaHome)
	}
}

func sortCandidates(jvms []Jvm, i int, j int) bool {
	if jvms[i].javaSpecificationVersion == jvms[j].javaSpecificationVersion {
		return jvms[i].javaHome > jvms[j].javaHome
	}
	return jvms[i].javaSpecificationVersion > jvms[j].javaSpecificationVersion
}
