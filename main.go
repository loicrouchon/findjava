package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type JvmInfo struct {
	javaPaths                []string
	javaHome                 string
	javaSpecificationVersion string
}

func (jvmInfo JvmInfo) String() string {
	return fmt.Sprintf(
		`{
    java: %q
    java.home: %s
    java.specification.version: %s
}`,
		jvmInfo.javaPaths,
		jvmInfo.javaHome,
		jvmInfo.javaSpecificationVersion)
}

func main() {
	var javaLookUpPaths = []string{
		"/bin/java",
		"/usr/bin/java",
		"/usr/local/bin/java",
		"/usr/lib/jvm",
		"~/.sdkman/candidates/java",
	}
	javaPaths := findAllJavaPaths(javaLookUpPaths)
	jvmInfos := make(map[string]JvmInfo)
	for javaPath, javaSymLinks := range javaPaths {
		jvmInfo := jvmInfo(javaPath, javaSymLinks)
		jvmInfos[javaPath] = jvmInfo
		fmt.Printf("%s: %s\n", javaPath, jvmInfo)
	}
}

func findAllJavaPaths(javaLookUpPaths []string) map[string][]string {
	javaPaths := make(map[string][]string)
	for _, javaLookUpPath := range javaLookUpPaths {
		fmt.Printf("Checking %s\n", javaLookUpPath)
		if strings.HasPrefix(javaLookUpPath, "~") {
			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			javaLookUpPath = strings.Replace(javaLookUpPath, "~", usr.HomeDir, 1)
            fmt.Printf("Updated lookup path %s\n", javaLookUpPath)
		}
		for _, javaPath := range findJavaPaths(javaLookUpPath) {
			fmt.Printf("  - Found %s\n", javaPath)
			resolvedJavaPath, err := filepath.EvalSymlinks(javaPath)
			if err != nil {
				fmt.Printf("%s cannot be resolved %s\n", javaPath, err)
				os.Exit(1)
			}
			if val, ok := javaPaths[resolvedJavaPath]; ok {
				javaPaths[resolvedJavaPath] = append(val, javaPath)
			} else {
				javaPaths[resolvedJavaPath] = []string{javaPath}
			}
		}
	}
	return javaPaths
}

func findJavaPaths(javaLookUpPath string) []string {
	if fileInfo, err := os.Stat(javaLookUpPath); err == nil {
		if !fileInfo.IsDir() {
			if fileInfo.Mode()&0111 != 0 {
				return []string{javaLookUpPath}
			} else {
				fmt.Printf("  File %s is not executable\n", javaLookUpPath)
			}
		} else {
			fmt.Printf("  File %s is a directory\n", javaLookUpPath)
			dir, err := os.Open(javaLookUpPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer dir.Close()

			// Read the directory contents
			files, err := dir.Readdir(-1)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			javaPaths := []string{}
			for _, file := range files {
				path := filepath.Join(javaLookUpPath, file.Name())
				fmt.Printf("  Looking into %s\n", path)
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

// isSymLink checks if the given path is a symbolic link
func isSymLink(path string) bool {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return fileInfo.Mode()&os.ModeSymlink != 0
}

// parseJavaVersion parses the version and the JDK vendor from the output of "java --version"
func jvmInfo(javaPath string, javaSymLinks []string) JvmInfo {
	// find . -mindepth 1 -maxdepth 1 -type d -print -exec {}/bin/java -XshowSettings:properties --version \;
	cmd := exec.Command(javaPath, "-cp", "build/classes", "JvmInfo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Fail to call %s %s\n", javaPath, err)
		os.Exit(1)
	}
	// fmt.Printf("%s\n", output)
	lines := strings.Split(string(output), "\n")
	var javaSpecificationVersion string
	var javaHome string
	for _, line := range lines {
		if value, found := strings.CutPrefix(line, "java.home="); found {
			javaHome = strings.TrimSpace(value)
		}
		if value, found := strings.CutPrefix(line, "java.specification.version="); found {
			javaSpecificationVersion = strings.TrimSpace(value)
		}
	}
	return JvmInfo{
		javaPaths:                javaSymLinks,
		javaHome:                 javaHome,
		javaSpecificationVersion: javaSpecificationVersion,
	}
}
