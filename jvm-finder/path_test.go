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
	os.Setenv("RESOLVE_PATH_ENV", "/resolve/path/env")
	testData := []TestData{
		{path: "", expectedPath: ""},
		{path: "$RESOLVE_PATH_NON_EXISTING_ENV", expectedPath: ""},
		{path: "$RESOLVE_PATH_ENV", expectedPath: "/resolve/path/env"},
		{path: "$RESOLVE_PATH_ENV/jdks/bin/java", expectedPath: "/resolve/path/env/jdks/bin/java"},
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
