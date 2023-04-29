package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type JavaExecutables struct {
	javaPaths map[string]time.Time
}

type JavaExecutable struct {
	path      string
	timestamp time.Time
}

func (javaExecutable *JavaExecutable) String() string {
	return fmt.Sprintf(`{timestamp: %s, path: %s}`, javaExecutable.timestamp, javaExecutable.path)
}

func findAllJavaExecutables(javaLookUpPaths *[]string) JavaExecutables {
	javaPaths := make(map[string]time.Time)
	for _, javaLookUpPath := range *javaLookUpPaths {
		logDebug("Checking %s", javaLookUpPath)
		for _, java := range findJavaExecutables(javaLookUpPath) {
			logDebug("  - Found %v", &java)
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
			} else {
				die("File %s (symlinked from %s) cannot be processed :(", path, lookUpPath)
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
	if java := findJavaExecutables(filepath.Join(directory, "bin", "java")); len(java) == 1 {
		return java
	}
	dir, err := os.Open(directory)
	if err != nil {
		dierr(err)
	}
	defer closeFile(dir)

	files, err := dir.Readdir(-1)
	if err != nil {
		dierr(err)
	}
	var javaPaths []JavaExecutable
	for _, file := range files {
		if !file.Mode().IsRegular() {
			path := filepath.Join(directory, file.Name(), "bin", "java")
			javaPaths = append(javaPaths, findJavaExecutables(path)...)
		}
	}
	return javaPaths
}
