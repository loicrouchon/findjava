package discovery

import (
	"findjava/internal/log"
	"findjava/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type JavaExecutables struct {
	// A map where keys are the absolute java executables path with all symbolic links resolved and keys are the
	// modification timestamp of the java executable.
	JavaPaths map[string]time.Time
}

type javaExecutable struct {
	path      string
	timestamp time.Time
}

func (javaExecutable *javaExecutable) String() string {
	return fmt.Sprintf(`{timestamp: %-30s, path: %s}`, javaExecutable.timestamp, javaExecutable.path)
}

// FindAllJavaExecutables locates all JVMs present in the `jvmLookupPaths`.
// Each path must either be an absolute path or a path relative to the user home (~).
//
// Paths are preprocessed to resolve environment variables (e.g. $JAVA_HOME, $JAVA_HOME/bin/java, etc.) and
// the user's home directory (e.g. ~/.sdkman/candidates/java).
//
// JVMs will be discovered for a given path in the following use cases:
//
//   - The path points to a file (after resolving symbolic links) that is executable
//     (e.g. /usr/bin/java, $JAVA_HOME/bin/java, etc.).
//   - The path points to a directory that contains (after resolving symbolic links) a bin/java executable
//     (e.g. $JAVA_HOME, $GRAALVM_HOME, etc.).
//   - If no bin/java executable is found, all direct subdirectories will be checked for <subdirectory>/bin/java
//     executables (e.g. /usr/lib/jvm, ~/.sdkman/candidates/java, /System/Volumes/Data/Library/Java/JavaVirtualMachines,
//     etc.). This will not recurse into subdirectories of subdirectories.
func FindAllJavaExecutables(javaLookUpPaths *[]string) (JavaExecutables, error) {
	javaPaths := make(map[string]time.Time)
	for _, javaLookUpPath := range *javaLookUpPaths {
		log.Debug("Checking %s", javaLookUpPath)
		javaExecutables, err := findJavaExecutables(javaLookUpPath)
		if err != nil {
			return JavaExecutables{}, err
		}
		for _, java := range javaExecutables {
			log.Debug("  - Found %v", &java)
			javaPaths[java.path] = java.timestamp
		}
	}
	return JavaExecutables{JavaPaths: javaPaths}, nil
}

func findJavaExecutables(lookUpPath string) ([]javaExecutable, error) {
	if path, err := filepath.EvalSymlinks(lookUpPath); err == nil {
		if fileInfo, err := os.Stat(path); err == nil {
			fileMode := fileInfo.Mode()
			if fileMode.IsRegular() {
				return findSingleJavaExecutable(path, fileInfo), nil
			} else if fileInfo.Mode().IsDir() {
				return findJavaExecutablesForEachJvmDirectory(path)
			} else {
				return nil, fmt.Errorf("file %s (symlinked from %s) cannot be processed :(", path, lookUpPath)
			}
		}
	}
	return []javaExecutable{}, nil
}

func findSingleJavaExecutable(path string, fileInfo os.FileInfo) []javaExecutable {
	if fileInfo.Mode()&0111 != 0 {
		return []javaExecutable{{
			path:      path,
			timestamp: fileInfo.ModTime(),
		}}
	} else {
		log.Debug("  File %s is not executable", path)
		return []javaExecutable{}
	}
}

func findJavaExecutablesForEachJvmDirectory(directory string) ([]javaExecutable, error) {
	if java, err := findJavaExecutables(filepath.Join(directory, "bin", "java")); len(java) == 1 {
		return nil, err
	}
	dir, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	defer utils.CloseFile(dir)

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	var javaPaths []javaExecutable
	for _, file := range files {
		if !file.Mode().IsRegular() {
			path := filepath.Join(directory, file.Name(), "bin", "java")
			javaExecutables, err := findJavaExecutables(path)
			if err != nil {
				return nil, err
			}
			javaPaths = append(javaPaths, javaExecutables...)
		}
	}
	return javaPaths, nil
}
