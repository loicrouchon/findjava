package main

import (
	"io/fs"
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

func findAllJavaExecutables(javaLookUpPaths []string) JavaExecutables {
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
		for _, java := range findJavaExecutables(javaLookUpPath) {
			logDebug("  - Found %v", java)
			javaPaths[java.path] = java.timestamp
		}
	}
	return JavaExecutables{javaPaths: javaPaths}
}

func findJavaExecutables(lookUpPath string) []JavaExecutable {
	if path, err := filepath.EvalSymlinks(lookUpPath); err == nil {
		if fileInfo, err := os.Stat(path); err == nil {
			fileMode := fileInfo.Mode()
			if fileMode.IsRegular() {
				return javaExecutable(path, fileInfo)
			} else if fileInfo.Mode().IsDir() {
				return javaExecutablesForEachJvmDirectory(path)
			}
		}
	}
	return []JavaExecutable{}
}

func javaExecutable(path string, fileInfo fs.FileInfo) []JavaExecutable {
	if fileInfo.Mode()&0111 != 0 {
		return []JavaExecutable{{
			path:      path,
			timestamp: fileInfo.ModTime(),
		}}
	} else {
		logDebug("  File %s is not executable", path)
		return []JavaExecutable{}
	}
}

func javaExecutablesForEachJvmDirectory(directory string) []JavaExecutable {
	dir, err := os.Open(directory)
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
		if !file.Mode().IsRegular() {
			path := filepath.Join(directory, file.Name(), "bin", "java")
			javaPaths = append(javaPaths, findJavaExecutables(path)...)
		}
	}
	return javaPaths
}
