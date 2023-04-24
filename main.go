package main

import (
	"fmt"
	"os"
	"os/exec"
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
		for _, javaPath := range findJavaPaths(javaLookUpPath) {
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
				fmt.Printf("File %s is not executable\n", javaLookUpPath)
			}
		} else {
			fmt.Printf("File %s is a directory\n", javaLookUpPath)

			// path := filepath.Join(javaLookUpPath, "bin", "java")
			// every single subdir

			// dir, err := os.Open(dirpath)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			// defer dir.Close()

			// // Read the directory contents
			// files, err := dir.Readdir(-1)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// 	// If the path is a directory and not a symbolic link
			// 	if file.IsDir() { // && !isSymLink(path) {
			// 		path := filepath.Join(dirpath, file.Name())
			// 		// Resolve any potential symbolic links in the path
			// 		resolvedPath, err := filepath.EvalSymlinks(path)
			// 		if err != nil {
			// 			fmt.Printf("%s cannot be resolved %s\n", path, err)
			// 			os.Exit(1)
			// 		}

			// 		// Execute "$DIR/bin/java --version" in the resolved path
			// 		javaPath := filepath.Join(jvmPath, "bin", "java")

			// 	}
			// }

			// if err != nil {
			// 	fmt.Println(err)
			// }
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
