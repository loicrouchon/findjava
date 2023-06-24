package jvm

import (
	"fmt"
	"time"
)

type Jvm struct {
	javaPath                 string
	JavaHome                 string
	JavaSpecificationVersion uint
	JavaVendor               string
	FetchedAt                time.Time
	SystemProperties         map[string]string
}

func (jvm *Jvm) rebuild() error {
	jvm.JavaHome = jvm.SystemProperties["java.home"]
	jvm.JavaVendor = jvm.SystemProperties["java.vendor"]
	if specVersion, err := ParseJavaSpecificationVersion(jvm.SystemProperties["java.specification.version"]); err != nil {
		return err
	} else {
		jvm.JavaSpecificationVersion = specVersion
	}
	return nil
}

func (jvm *Jvm) String() string {
	return fmt.Sprintf(
		`[%v]
timestamp: %s
java.home: %s
java.specification.version: %d
`,
		jvm.javaPath,
		jvm.FetchedAt,
		jvm.JavaHome,
		jvm.JavaSpecificationVersion)
}
