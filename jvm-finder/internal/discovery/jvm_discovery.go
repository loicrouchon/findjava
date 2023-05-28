package discovery

import (
	"fmt"
	"jvm-finder/internal/log"
	"jvm-finder/internal/utils"
	"os"
	"path/filepath"
	"time"
)

type JavaExecutables struct {
	JavaPaths map[string]time.Time
}

type JavaExecutable struct {
	path      string
	timestamp time.Time
}

func (javaExecutable *JavaExecutable) String() string {
	return fmt.Sprintf(`{timestamp: %-30s, path: %s}`, javaExecutable.timestamp, javaExecutable.path)
}

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

func findJavaExecutables(lookUpPath string) ([]JavaExecutable, error) {
	if path, err := filepath.EvalSymlinks(lookUpPath); err == nil {
		if fileInfo, err := os.Stat(path); err == nil {
			fileMode := fileInfo.Mode()
			if fileMode.IsRegular() {
				return javaExecutable(path, fileInfo), nil
			} else if fileInfo.Mode().IsDir() {
				return javaExecutablesForEachJvmDirectory(path)
			} else {
				return nil, fmt.Errorf("file %s (symlinked from %s) cannot be processed :(", path, lookUpPath)
			}
		}
	}
	return []JavaExecutable{}, nil
}

func javaExecutable(path string, fileInfo os.FileInfo) []JavaExecutable {
	if fileInfo.Mode()&0111 != 0 {
		return []JavaExecutable{{
			path:      path,
			timestamp: fileInfo.ModTime(),
		}}
	} else {
		log.Debug("  File %s is not executable", path)
		return []JavaExecutable{}
	}
}

func javaExecutablesForEachJvmDirectory(directory string) ([]JavaExecutable, error) {
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
	var javaPaths []JavaExecutable
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
