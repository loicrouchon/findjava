package utils

import (
	"findjava/internal/log"
	"findjava/test"
	"fmt"
	"os"
	"os/user"
	"testing"
)

func TestResolvePath(t *testing.T) {
	type TestData struct {
		path, expectedPath, err string
	}
	var userHome string
	if u, err := user.Current(); err != nil {
		log.Die(err)
	} else {
		userHome = u.HomeDir
	}
	_ = os.Setenv("RESOLVE_PATH_ENV", "/resolve/path/env")
	testData := []TestData{
		{path: "", expectedPath: ""},
		{path: "$RESOLVE_PATH_NON_EXISTING_ENV", expectedPath: "",
			err: "env var RESOLVE_PATH_NON_EXISTING_ENV not found -> cannot process path $RESOLVE_PATH_NON_EXISTING_ENV"},
		{path: "$RESOLVE_PATH_ENV", expectedPath: "/resolve/path/env"},
		{path: "$RESOLVE_PATH_ENV/jdks/bin/java", expectedPath: "/resolve/path/env/jdks/bin/java"},
		{path: "~/jdks/bin/java", expectedPath: userHome + "/jdks/bin/java"},
	}
	for _, data := range testData {
		actualPath, err := ResolvePath(data.path)
		description := fmt.Sprintf("ResolvePath(\"%s\")", data.path)
		test.AssertEquals(t, description, data.expectedPath, actualPath)
		test.AssertErrorContains(t, description, data.err, err)
	}
}
