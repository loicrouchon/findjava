package main

import (
	"os"
	"os/user"
	"regexp"
	"strings"
)

var envVarsRegexp, _ = regexp.Compile(`\$([a-zA-Z0-9_]+)`)

func resolvePaths(paths []string) []string {
	var resolvedPaths []string
	for _, path := range paths {
		resolvedPath := resolvePath(path)
		if resolvedPath != "" {
			resolvedPaths = append(resolvedPaths, resolvedPath)
		}
	}
	return resolvedPaths
}

func resolvePath(path string) string {
	validPath := true
	if strings.Contains(path, "$") {
		path = string(envVarsRegexp.ReplaceAllFunc([]byte(path),
			func(match []byte) []byte { return expandEnvVar(path, &validPath, match) }))
	}
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			logInfo("Unable to resolve user home directory -> cannot process path %s", path)
			validPath = false
		} else {
			path = strings.Replace(path, "~", usr.HomeDir, 1)
		}
	}
	if !validPath {
		return ""
	}
	return path
}

func expandEnvVar(path string, validPath *bool, envVarName []byte) []byte {
	envVar := string(envVarName)[1:]
	if value, found := os.LookupEnv(envVar); found {
		return []byte(value)
	}
	logInfo("Env var %s not found -> cannot process path %s", envVar, path)
	*validPath = false
	return nil
}
