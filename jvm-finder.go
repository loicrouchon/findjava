package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type JvmSelectionRules struct {
	minJvmVersion int
	maxJvmVersion int
}

func (rules JvmSelectionRules) Matches(jvmInfo JvmInfo) bool {
	if rules.minJvmVersion > 0 && rules.minJvmVersion > jvmInfo.javaSpecificationVersion {
		return false
	}
	if rules.maxJvmVersion > 0 && rules.maxJvmVersion < jvmInfo.javaSpecificationVersion {
		return false
	}
	return true
}

func (rules JvmSelectionRules) String() string {
	return fmt.Sprintf("[%d..%d]}", rules.minJvmVersion, rules.maxJvmVersion)
}

type JvmInfo struct {
	javaPath                 string
	javaHome                 string
	javaSpecificationVersion int
}

func (jvmInfo JvmInfo) String() string {
	return fmt.Sprintf(
		`[%s]
java.home: %s
java.specification.version: %d
`,
		jvmInfo.javaPath,
		jvmInfo.javaHome,
		jvmInfo.javaSpecificationVersion)
}

func main() {
	args := parseArgs()
	if len(args) > 1 {
		logError("Usage jvm-finder [VERSION]")
		logError("  VERSION: A JVM version range:")
		logError("      - 17        exact version)")
		logError("      - 17..      17 or above)")
		logError("      - ..17      up to 17")
		logError("      - 11..17    From 11 to 17")
		os.Exit(1)
	}
	rules := jvmSelectionRules(args)

	var javaLookUpPaths = []string{
		"/bin/java",
		"/usr/bin/java",
		"/usr/local/bin/java",
		"/usr/lib/jvm",
		"~/.sdkman/candidates/java",
	}
	jvmInfos := loadJvmInfos("./build/jvm-finder.properties")
	javaPaths := findAllJavaPaths(javaLookUpPaths)
	for javaPath, _ := range javaPaths {
		jvmInfos.Fetch(javaPath)
	}
	jvmInfos.Save()
	var matchingJvms []JvmInfo
	for _, jvmInfo := range jvmInfos.jvmInfos {
		if rules.Matches(jvmInfo) {
			matchingJvms = append(matchingJvms, jvmInfo)
			logInfo("[CANDIDATE] %s (%d)", jvmInfo.javaHome, jvmInfo.javaSpecificationVersion)
		} else {
			logInfo("[IGNORED]   %s (%d)", jvmInfo.javaHome, jvmInfo.javaSpecificationVersion)
		}
	}
	sort.Slice(matchingJvms[:], func(i, j int) bool {
		if matchingJvms[i].javaSpecificationVersion == matchingJvms[j].javaSpecificationVersion {
			return matchingJvms[i].javaHome > matchingJvms[j].javaHome
		}
		return matchingJvms[i].javaSpecificationVersion > matchingJvms[j].javaSpecificationVersion
	})
	logDebug("%v\n", matchingJvms)
	if matchingJvms != nil && len(matchingJvms) > 0 {
		selectedJvm := matchingJvms[0]
		logInfo("[SELECTED]  %s (%d)", selectedJvm.javaHome, selectedJvm.javaSpecificationVersion)
		fmt.Printf("%s\n", filepath.Join(selectedJvm.javaHome, "bin", "java"))
	} else {
		logError("Unable to find a JVM matching requirements %s", rules)
		os.Exit(1)
	}
}

func jvmSelectionRules(args []string) JvmSelectionRules {
	var rules JvmSelectionRules
	if len(args) == 1 {
		jvmVersionRange := args[0]
		logDebug("%s", jvmVersionRange)
		jvmVersionRegex := `[\d]+(?:\.[\d]+)*`
		r := regexp.MustCompile(fmt.Sprintf(`^(?:`+
			`(?P<exact>%[1]s)`+
			`|(?:(?P<min>%[1]s)\.\.)`+
			`|(?:\.\.(?P<max>%[1]s))`+
			`|(?:(?P<min>%[1]s)\.\.(?P<max>%[1]s))`+
			`)$`,
			jvmVersionRegex))
		groupNames := r.SubexpNames()
		match := r.FindStringSubmatch(jvmVersionRange)
		var minJvmVersion int
		var maxJvmVersion int
		for i, m := range match {
			if len(m) > 0 {
				switch groupNames[i] {
				case "exact":
					minJvmVersion = parseVersion(m)
					maxJvmVersion = parseVersion(m)
				case "min":
					minJvmVersion = parseVersion(m)
				case "max":
					maxJvmVersion = parseVersion(m)
				}
			}
		}
		rules = JvmSelectionRules{
			minJvmVersion: minJvmVersion,
			maxJvmVersion: maxJvmVersion,
		}
	} else {
		rules = JvmSelectionRules{}
	}
	logDebug("%s", rules)
	return rules
}

func parseVersion(version string) int {
	switch version {
	case "1.0", "1.1":
		return 1
	case "1.2":
		return 2
	case "1.3":
		return 3
	case "1.4":
		return 4
	case "1.5":
		return 5
	case "1.6":
		return 6
	case "1.7":
		return 7
	case "1.8":
		return 8
	default:
		v, err := strconv.Atoi(version)
		if err != nil {
			logError("JVM version %s cannot be parsed as an int")
			panic(version)
		}
		return v
	}
}

func findAllJavaPaths(javaLookUpPaths []string) map[string][]string {
	javaPaths := make(map[string][]string)
	for _, javaLookUpPath := range javaLookUpPaths {
		if strings.HasPrefix(javaLookUpPath, "~") {
			usr, err := user.Current()
			if err != nil {
				log.Fatal(err)
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
	return javaPaths
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

func jvmInfo(javaPath string) JvmInfo {
	cmd := exec.Command(javaPath, "-cp", "build/classes", "JvmInfo")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logError("Fail to call %s %s", javaPath, err)
		os.Exit(1)
	}
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
		javaPath:                 javaPath,
		javaHome:                 javaHome,
		javaSpecificationVersion: parseVersion(javaSpecificationVersion),
	}
}

var logLevel string

func parseArgs() []string {
	flag.StringVar(&logLevel, "loglevel", "error", "Log level: debug, info, error")
	flag.Parse()
	if logLevel != "debug" && logLevel != "info" && logLevel != "error" {
		logError("Invalid log level: '%s'. Available levels are: debug, info, error", logLevel)
		os.Exit(1)
	}
	return flag.Args()
}

func logDebug(message string, v ...any) {
	if logLevel == "debug" {
		fmt.Fprintf(os.Stdout, "[DEBUG] %s\n", fmt.Sprintf(message, v...))
	}
}
func logInfo(message string, v ...any) {
	if logLevel == "debug" || logLevel == "info" {
		fmt.Fprintf(os.Stdout, "[INFO] %s\n", fmt.Sprintf(message, v...))
	}
}

func logError(message string, v ...any) {
	fmt.Fprintf(os.Stderr, "[ERROR] %s\n", fmt.Sprintf(message, v...))
}

type JvmInfos struct {
	path       string
	jvmInfos   map[string]JvmInfo
	dirtyCache bool
}

func loadJvmInfos(path string) JvmInfos {
	infos := make(map[string]JvmInfo)
	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		if err != nil {
			logError("Unable to read file %s: %s", path, err)
		}
		defer file.Close()

		var jvmInfo JvmInfo
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if value, found := strings.CutPrefix(line, "["); found {
				if len(jvmInfo.javaHome) > 0 && jvmInfo.javaSpecificationVersion > 0 {
					infos[jvmInfo.javaPath] = jvmInfo
				}
				jvmInfo = JvmInfo{javaPath: strings.Trim(value, "[]")}
			} else if value, ok := strings.CutPrefix(line, "java.home="); ok {
				jvmInfo.javaHome = value
			} else if value, ok := strings.CutPrefix(line, "java.specification.version="); ok {
				if version, err := strconv.Atoi(value); err == nil {
					jvmInfo.javaSpecificationVersion = version
				}
			}
		}
		if len(jvmInfo.javaHome) > 0 && jvmInfo.javaSpecificationVersion > 0 {
			infos[jvmInfo.javaPath] = jvmInfo
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return JvmInfos{
		path:       path,
		jvmInfos:   infos,
		dirtyCache: false,
	}
}

func (cache *JvmInfos) Fetch(javaPath string) {
	if _, found := cache.jvmInfos[javaPath]; !found {
		logInfo("[CACHE MISS] %s", javaPath)
		jvmInfo := jvmInfo(javaPath)
		cache.jvmInfos[javaPath] = jvmInfo
		cache.dirtyCache = true
		logDebug("%s: %s", javaPath, jvmInfo)
	}
}

func (cache *JvmInfos) Save() {
	if cache.dirtyCache {
		fmt.Printf("Saving %s\n", cache.path)
		output := ""
		for _, jvmInfo := range cache.jvmInfos {
			output = fmt.Sprintf(`%s
[%s]
java.home=%s
java.specification.version=%d
`, output, jvmInfo.javaPath, jvmInfo.javaHome, jvmInfo.javaSpecificationVersion)
		}
		logInfo(output)
		if err := os.WriteFile(cache.path, []byte(output), 0666); err != nil {
			logError("Unable to write to file %s, %s", cache.path, err)
			os.Exit(1)
		}
	}
}
