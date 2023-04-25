package main

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type JavaExecutables struct {
	javaPaths map[string]time.Time
}

type JavaExecutable struct {
	path      string
	timestamp time.Time
}

func findAllJavaPaths(javaLookUpPaths []string) JavaExecutables {
	javaPaths := make(map[string]time.Time)
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
		for _, java := range findJavaPaths(javaLookUpPath) {
			logDebug("  - Found %s", java)
			resolvedJavaPath, err := filepath.EvalSymlinks(java.path)
			if err != nil {
				logError("%s cannot be resolved %s", java, err)
				os.Exit(1)
			}
			javaPaths[resolvedJavaPath] = java.timestamp
		}
	}
	return JavaExecutables{javaPaths: javaPaths}
}

func findJavaPaths(javaLookUpPath string) []JavaExecutable {
	if fileInfo, err := os.Stat(javaLookUpPath); err == nil {
		if !fileInfo.IsDir() {
			if fileInfo.Mode()&0111 != 0 {
				return []JavaExecutable{JavaExecutable{
					path:      javaLookUpPath,
					timestamp: fileInfo.ModTime(),
				}}
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
			javaPaths := []JavaExecutable{}
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
	return []JavaExecutable{}
}

func isSymLink(path string) bool {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fileInfo.Mode()&os.ModeSymlink != 0
}
