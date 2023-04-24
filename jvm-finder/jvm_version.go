package main

import (
	"strconv"
)

func parseVersion(version string) int {
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
		if err != nil {
			logError("JVM version %s cannot be parsed as an int")
			panic(version)
		}
		return v
	}
}
