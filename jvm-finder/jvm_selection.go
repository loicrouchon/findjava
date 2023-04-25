package main

import (
	"sort"
)

func (jvmInfos *JvmInfos) Select(rules *JvmSelectionRules) (*JvmInfo, bool) {
	var matchingJvms []JvmInfo
	for _, jvmInfo := range jvmInfos.jvmInfos {
		if rules.Matches(&jvmInfo) {
			matchingJvms = append(matchingJvms, jvmInfo)
			logInfo("[CANDIDATE] %s (%d)", jvmInfo.javaHome, jvmInfo.javaSpecificationVersion)
		} else {
			logInfo("[IGNORED]   %s (%d)", jvmInfo.javaHome, jvmInfo.javaSpecificationVersion)
		}
	}
	sort.Slice(matchingJvms[:], func(i, j int) bool {
		if matchingJvms[i].javaSpecificationVersion == matchingJvms[j].javaSpecificationVersion {
			return matchingJvms[i].javaHome > matchingJvms[j].javaHome
		}
		return matchingJvms[i].javaSpecificationVersion > matchingJvms[j].javaSpecificationVersion
	})
	logDebug("%v\n", matchingJvms)
	if matchingJvms != nil && len(matchingJvms) > 0 {
		return &matchingJvms[0], true
	} else {
		return nil, false
	}
}
