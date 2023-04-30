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

func parseVersion(version string) uint {
	switch version {
	case "1.0", "1.1":
		return 1
	case "1.2":
		return 2
	case "1.3":
		return 3
	case "1.4":
		return 4
	case "1.5":
		return 5
	case "1.6":
		return 6
	case "1.7":
		return 7
	case "1.8":
		return 8
	default:
		v, err := strconv.Atoi(version)
		if err != nil || v < 0 {
			die("JVM version %s cannot be parsed as an unsigned int", version)
		}
		return uint(v)
	}
}
