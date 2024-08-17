package selection

import (
	"findjava/internal/jvm"
	"findjava/internal/rules"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSelectJvm(t *testing.T) {
	type TestData struct {
		description string
		rules       rules.JvmSelectionRules
		jvms        jvm.JvmsInfos
		expectedJvm []*jvm.Jvm
	}
	jvm17Ubuntu := testJvm("java-17-openjdk-amd64", 17, "Ubuntu")
	jvm21Ubuntu := testJvm("java-21-openjdk-amd64", 21, "Ubuntu")
	jvm21EclipseAdoptium := testJvm("21.0.1-tem", 21, "Eclipse Adoptium")
	jvm22EclipseAdoptium := testJvm("22.0.1-tem", 22, "Eclipse Adoptium")
	jvm22GraalVMCE := testJvm("22-graalce", 22, "GraalVM Community")
	jvmInfos := jvmsInfos(
		jvm17Ubuntu,
		jvm21Ubuntu,
		jvm21EclipseAdoptium,
		jvm22EclipseAdoptium,
		jvm22GraalVMCE,
	)
	testData := []TestData{
		{
			description: "should select most recent JVM",
			rules: rules.JvmSelectionRules{
				VersionRange: &jvm.VersionRange{Min: jvm.AllVersions, Max: jvm.AllVersions},
			},
			jvms:        jvmInfos,
			expectedJvm: []*jvm.Jvm{jvm22EclipseAdoptium, jvm22GraalVMCE},
		},
		{
			description: "should filter on JVM version (21)",
			rules: rules.JvmSelectionRules{
				VersionRange: &jvm.VersionRange{Min: 8, Max: 21},
			},
			jvms:        jvmInfos,
			expectedJvm: []*jvm.Jvm{jvm21EclipseAdoptium, jvm21Ubuntu},
		},
		{
			description: "should filter on JVM version (17)",
			rules: rules.JvmSelectionRules{
				VersionRange: &jvm.VersionRange{Min: jvm.AllVersions, Max: 17},
			},
			jvms:        jvmInfos,
			expectedJvm: []*jvm.Jvm{jvm17Ubuntu},
		},
		{
			description: "should filter on vendor (Ubuntu)",
			rules: rules.JvmSelectionRules{
				VersionRange: &jvm.VersionRange{Min: jvm.AllVersions, Max: jvm.AllVersions},
				Vendors:      []string{"Ubuntu"},
			},
			jvms:        jvmInfos,
			expectedJvm: []*jvm.Jvm{jvm21Ubuntu},
		},
		{
			description: "should filter on vendor (Eclipse Adoptium)",
			rules: rules.JvmSelectionRules{
				VersionRange: &jvm.VersionRange{Min: jvm.AllVersions, Max: jvm.AllVersions},
				Vendors:      []string{"Eclipse Adoptium"},
			},
			jvms:        jvmInfos,
			expectedJvm: []*jvm.Jvm{jvm22EclipseAdoptium},
		},
		{
			description: "should filter on vendor (GraalVM Community)",
			rules: rules.JvmSelectionRules{
				VersionRange: &jvm.VersionRange{Min: jvm.AllVersions, Max: jvm.AllVersions},
				Vendors:      []string{"GraalVM Community"},
			},
			jvms:        jvmInfos,
			expectedJvm: []*jvm.Jvm{jvm22GraalVMCE},
		},
		{
			description: "should satisfy main rules & preferred rules",
			rules: rules.JvmSelectionRules{
				VersionRange: &jvm.VersionRange{Min: jvm.AllVersions, Max: jvm.AllVersions},
				PreferredRules: &rules.JvmSelectionRules{
					VersionRange: &jvm.VersionRange{Min: jvm.AllVersions, Max: 21},
					Vendors:      []string{"Ubuntu"},
				},
			},
			jvms:        jvmInfos,
			expectedJvm: []*jvm.Jvm{jvm21Ubuntu},
		},
		{
			description: "should satisfy only main rules when not possible to satisfy preferred rules",
			rules: rules.JvmSelectionRules{
				VersionRange: &jvm.VersionRange{Min: 22, Max: jvm.AllVersions},
				PreferredRules: &rules.JvmSelectionRules{
					VersionRange: &jvm.VersionRange{Min: jvm.AllVersions, Max: jvm.AllVersions},
					Vendors:      []string{"Ubuntu"},
				},
			},
			jvms:        jvmInfos,
			expectedJvm: []*jvm.Jvm{jvm22EclipseAdoptium, jvm22GraalVMCE},
		},
		{
			description: "should fail to select a JVM when none match rules",
			rules: rules.JvmSelectionRules{
				VersionRange: &jvm.VersionRange{Min: 15, Max: 16},
			},
			jvms:        jvmInfos,
			expectedJvm: []*jvm.Jvm{nil},
		},
		// TODO find a way to add rules check for programs (this requires more test data preparation)
	}
	for _, data := range testData {
		selectedJvm := Select(&data.rules, &data.jvms)
		match := false
		for _, expectedJvm := range data.expectedJvm {
			if reflect.DeepEqual(selectedJvm, expectedJvm) {
				match = true
				break
			}
		}
		if !match {
			t.Fatalf(`## %s ## Expecting: Select("%v", jvms) to be one of %v but was %v`,
				data.description, &data.rules, data.expectedJvm, selectedJvm)
		}
	}
}

func jvmsInfos(jvms ...*jvm.Jvm) jvm.JvmsInfos {
	jvmsInfos := jvm.JvmsInfos{
		Jvms: make(map[string]*jvm.Jvm),
	}
	for _, jvm := range jvms {
		javaPath := filepath.Join(jvm.JavaHome, "bin", "java")
		jvmsInfos.Jvms[javaPath] = jvm
	}
	return jvmsInfos
}

func testJvm(name string, javaSpecificationVersion uint, javaVendor string) *jvm.Jvm {
	return &jvm.Jvm{
		JavaHome:                 fmt.Sprintf("/jvms/%s", name),
		JavaSpecificationVersion: javaSpecificationVersion,
		JavaVendor:               javaVendor,
	}
}
