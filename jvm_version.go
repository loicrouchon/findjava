package main

import (
	"fmt"
	"strconv"
)

const allVersions = 0

type VersionRange struct {
	Min uint
	Max uint
}

func (versionRange *VersionRange) Matches(version uint) bool {
	if versionRange.Min != allVersions && versionRange.Min > version {
		return false
	}
	if versionRange.Max != allVersions && versionRange.Max < version {
		return false
	}
	return true
}

func (versionRange *VersionRange) String() string {
	return fmt.Sprintf("[%s..%s]}", str(versionRange.Min), str(versionRange.Max))
}

func str(version uint) string {
	if version == allVersions {
		return ""
	} else {
		return strconv.Itoa(int(version))
	}
}

func parseJavaSpecificationVersion(version string) (uint, error) {
	var javaSpecificationVersion uint
	switch version {
	case "1.0", "1.1":
		javaSpecificationVersion = 1
	case "1.2":
		javaSpecificationVersion = 2
	case "1.3":
		javaSpecificationVersion = 3
	case "1.4":
		javaSpecificationVersion = 4
	case "1.5":
		javaSpecificationVersion = 5
	case "1.6":
		javaSpecificationVersion = 6
	case "1.7":
		javaSpecificationVersion = 7
	case "1.8":
		javaSpecificationVersion = 8
	default:
		v, err := strconv.Atoi(version)
		if err != nil || v < 0 {
			return 0, fmt.Errorf("JVM version '%s' cannot be parsed as an unsigned int", version)
		}
		javaSpecificationVersion = uint(v)
	}
	return javaSpecificationVersion, nil
}
