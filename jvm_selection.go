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
	var candidates []Jvm
	var ignored []Jvm
	for _, jvm := range jvms.Jvms {
		if rules.Matches(jvm) {
			candidates = append(candidates, *jvm)
		} else {
			ignored = append(ignored, *jvm)
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
