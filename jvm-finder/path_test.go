package main

import (
	"os"
	"os/user"
	"testing"
)

func TestResolvePath(t *testing.T) {
	type TestData struct {
		path, expectedPath string
	}
	var userHome string
	if u, err := user.Current(); err != nil {
		dierr(err)
	} else {
		userHome = u.HomeDir
	}
	envVarHome := os.Getenv("HOME")
	testData := []TestData{
		{path: "", expectedPath: ""},
		{path: "$HOME", expectedPath: envVarHome},
		{path: "$HOME/jdks/bin/java", expectedPath: envVarHome + "/jdks/bin/java"},
		{path: "~/jdks/bin/java", expectedPath: userHome + "/jdks/bin/java"},
	}
	for _, data := range testData {
		actualPath := resolvePath(data.path)
		if actualPath != data.expectedPath {
			t.Fatalf(`Expecting resolvePath("%s") == %s but was %s`,
				data.path, data.expectedPath, actualPath)
		}
	}
}
