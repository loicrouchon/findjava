package utils

import (
	"findjvm/internal/log"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strings"
)

var envVarsRegexp, _ = regexp.Compile(`\$([a-zA-Z0-9_]+)`)

func ResolvePaths(paths []string) []string {
	var resolvedPaths []string
	for _, path := range paths {
		if resolvedPath, err := ResolvePath(path); err != nil {
			log.Info(err.Error())
		} else {
			resolvedPaths = append(resolvedPaths, resolvedPath)
		}
	}
	return resolvedPaths
}

func ResolvePath(path string) (string, error) {
	if strings.Contains(path, "$") {
		var err error
		path = string(envVarsRegexp.ReplaceAllFunc([]byte(path),
			func(match []byte) []byte { return expandEnvVar(path, &err, match) }))
		if err != nil {
			return "", err
		}
	}
	if strings.HasPrefix(path, "~") {
		if usr, err := user.Current(); err != nil {
			return "", fmt.Errorf("unable to resolve user home directory -> cannot process path %s", path)
		} else {
			path = strings.Replace(path, "~", usr.HomeDir, 1)
		}
	}
	return path, nil
}

func expandEnvVar(path string, err *error, envVarName []byte) []byte {
	envVar := string(envVarName)[1:]
	if value, found := os.LookupEnv(envVar); found {
		return []byte(value)
	}
	*err = fmt.Errorf("env var %s not found -> cannot process path %s", envVar, path)
	return nil
}
