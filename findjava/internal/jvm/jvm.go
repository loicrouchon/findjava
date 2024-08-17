package jvm

import (
	"fmt"
	"time"
)

// Jvm represents a JVM and its associated metadata as extracted by the [MetadataReader].
type Jvm struct {
	// The absolute path to the java executable.
	javaPath string
	// The absolute path to the `java.home` directory.
	JavaHome string
	// The `java.specification.version`
	JavaSpecificationVersion uint
	// The vendor
	JavaVendor string
	// The time at which the metadata were read
	FetchedAt time.Time
	// The system properties extracted by the [MetadataReader]
	SystemProperties map[string]string
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
