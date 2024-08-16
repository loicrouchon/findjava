package selection

import (
	"findjava/internal/jvm"
	"findjava/internal/log"
	"findjava/internal/rules"
	"sort"
)

// Select finds the best JVM from the [jvm.JvmsInfos] according the to the [rules.JvmSelectionRules].
// The candidates filtering process is performed in two steps:
//
//   - Determinate a list of candidates according to the [rules.JvmSelectionRules] but without considering the
//     preferred rules coming from configuration.
//   - In case at least one candidate is found, the preferred rules will then be checked to narrow down the candidates
//     list. If they cannot be fulfilled, they will be ignored.
//
// Once a list of potential candidates has been established, an election process is initiated to select which one
// of these candidates shall be selected.
// The election process return the JVM implementing the highest java.specification.version. If multiple JVMs implement
// the same java.specification.version, one will be selected. This selection process is not currently specified nor
// deterministic. Future versions of might provide rules for preferred JVM selection in such cases.
func Select(rules *rules.JvmSelectionRules, jvms *jvm.JvmsInfos) *jvm.Jvm {
	if candidates := selectCandidates(rules, jvms); len(candidates) > 0 {
		candidate := candidates[0]
		logJvmList("[SELECTED]", []jvm.Jvm{candidate})
		return &candidate
	} else {
		return nil
	}
}
func selectCandidates(rules *rules.JvmSelectionRules, jvms *jvm.JvmsInfos) []jvm.Jvm {
	candidates, ignored := filter(rules, jvms)
	sort.Slice(ignored[:], func(i, j int) bool { return sortCandidates(ignored, i, j) })
	sort.Slice(candidates[:], func(i, j int) bool { return sortCandidates(candidates, i, j) })
	logJvmList("[IGNORED]", ignored)
	logJvmList("[CANDIDATE]", candidates)
	return candidates
}

func filter(rules *rules.JvmSelectionRules, jvms *jvm.JvmsInfos) ([]jvm.Jvm, []jvm.Jvm) {
	var allJvms []jvm.Jvm
	for _, jvm := range jvms.Jvms {
		allJvms = append(allJvms, *jvm)
	}
	return filterJvmList(rules, allJvms)
}

func filterJvmList(rules *rules.JvmSelectionRules, allJvms []jvm.Jvm) ([]jvm.Jvm, []jvm.Jvm) {
	var candidates []jvm.Jvm
	var ignored []jvm.Jvm
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
			for _, ig := range preferredIgnored {
				ignored = append(ignored, ig)
			}
			return preferredCandidates, ignored
		} else {
			log.Info("Unable to satisfy preferred selection rules %v, ignoring them", rules.PreferredRules)
		}
	}
	return candidates, ignored
}

func logJvmList(displayType string, jvms []jvm.Jvm) {
	for i := len(jvms) - 1; i >= 0; i = i - 1 {
		jvm := jvms[i]
		log.Info("%-12s %3d: %s ", displayType, jvm.JavaSpecificationVersion, jvm.JavaHome)
	}
}

func sortCandidates(jvms []jvm.Jvm, i int, j int) bool {
	if jvms[i].JavaSpecificationVersion == jvms[j].JavaSpecificationVersion {
		return jvms[i].JavaHome > jvms[j].JavaHome
	}
	return jvms[i].JavaSpecificationVersion > jvms[j].JavaSpecificationVersion
}
