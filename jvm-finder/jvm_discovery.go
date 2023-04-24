package main

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func findAllJavaPaths(javaLookUpPaths []string) []string {
	javaPaths := make(map[string][]string)
	for _, javaLookUpPath := range javaLookUpPaths {
		if strings.HasPrefix(javaLookUpPath, "~") {
			usr, err := user.Current()
			if err != nil {
				logError("Unable to resolve user home directory used in path %s: %s", javaLookUpPath, err)
				os.Exit(1)
			}
			javaLookUpPath = strings.Replace(javaLookUpPath, "~", usr.HomeDir, 1)
		}
		logDebug("Checking %s", javaLookUpPath)
		for _, javaPath := range findJavaPaths(javaLookUpPath) {
			logDebug("  - Found %s", javaPath)
			resolvedJavaPath, err := filepath.EvalSymlinks(javaPath)
			if err != nil {
				logError("%s cannot be resolved %s", javaPath, err)
				os.Exit(1)
			}
			if val, ok := javaPaths[resolvedJavaPath]; ok {
				javaPaths[resolvedJavaPath] = append(val, javaPath)
			} else {
				javaPaths[resolvedJavaPath] = []string{javaPath}
			}
		}
	}
	resolvedPaths := make([]string, 0, len(javaPaths))
	for path := range javaPaths {
		resolvedPaths = append(resolvedPaths, path)
	}

	return resolvedPaths
}

func findJavaPaths(javaLookUpPath string) []string {
	if fileInfo, err := os.Stat(javaLookUpPath); err == nil {
		if !fileInfo.IsDir() {
			if fileInfo.Mode()&0111 != 0 {
				return []string{javaLookUpPath}
			} else {
				logDebug("  File %s is not executable", javaLookUpPath)
			}
		} else {
			dir, err := os.Open(javaLookUpPath)
			if err != nil {
				logError("%s", err)
				os.Exit(1)
			}
			defer dir.Close()

			files, err := dir.Readdir(-1)
			if err != nil {
				logError("%s", err)
				os.Exit(1)
			}
			javaPaths := []string{}
			for _, file := range files {
				path := filepath.Join(javaLookUpPath, file.Name())
				if file.IsDir() || isSymLink(path) {
					javaPath := filepath.Join(path, "bin", "java")
					javaPaths = append(javaPaths, findJavaPaths(javaPath)...)
				}
			}
			return javaPaths
		}
	}
	return []string{}
}

func isSymLink(path string) bool {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fileInfo.Mode()&os.ModeSymlink != 0
}
