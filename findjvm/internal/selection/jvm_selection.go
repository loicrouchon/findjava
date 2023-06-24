package selection

import (
	. "findjvm/internal/jvm"
	"findjvm/internal/log"
	"findjvm/internal/rules"
	"sort"
)

func Select(rules *rules.JvmSelectionRules, jvms *JvmsInfos) []Jvm {
	candidates, ignored := filter(rules, jvms)
	sort.Slice(ignored[:], func(i, j int) bool { return sortCandidates(ignored, i, j) })
	sort.Slice(candidates[:], func(i, j int) bool { return sortCandidates(candidates, i, j) })
	LogJvmList("[IGNORED]", ignored)
	LogJvmList("[CANDIDATE]", candidates)
	return candidates
}

func filter(rules *rules.JvmSelectionRules, jvms *JvmsInfos) ([]Jvm, []Jvm) {
	var allJvms []Jvm
	for _, jvm := range jvms.Jvms {
		allJvms = append(allJvms, *jvm)
	}
	return filterJvmList(rules, allJvms)
}

func filterJvmList(rules *rules.JvmSelectionRules, allJvms []Jvm) ([]Jvm, []Jvm) {
	var candidates []Jvm
	var ignored []Jvm
	for _, jvm := range allJvms {
		if rules.Matches(&jvm) {
			candidates = append(candidates, jvm)
		} else {
			ignored = append(ignored, jvm)
		}
	}
	if len(candidates) > 0 && rules.PreferredRules != nil {
		preferredCandidates, preferredIgnored := filterJvmList(rules.PreferredRules, candidates)
		if len(preferredCandidates) > 0 {
			return preferredCandidates, preferredIgnored
		} else if !rules.VersionRange.IsBounded() {
			return nil, preferredIgnored
		} else {
			log.Info("Unable to satisfy preferred selection rules %v, ignoring them", rules.PreferredRules)
		}
	}
	return candidates, ignored
}

func LogJvmList(displayType string, jvms []Jvm) {
	for i := len(jvms) - 1; i >= 0; i = i - 1 {
		jvm := jvms[i]
		log.Info("%-12s %3d: %s ", displayType, jvm.JavaSpecificationVersion, jvm.JavaHome)
	}
}

func sortCandidates(jvms []Jvm, i int, j int) bool {
	if jvms[i].JavaSpecificationVersion == jvms[j].JavaSpecificationVersion {
		return jvms[i].JavaHome > jvms[j].JavaHome
	}
	return jvms[i].JavaSpecificationVersion > jvms[j].JavaSpecificationVersion
}
