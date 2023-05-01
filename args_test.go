package main

import (
	"fmt"
	"testing"
)

func TestParseArgs(t *testing.T) {
	type TestData struct {
		args     []string
		expected Args
		err      error
	}
	data := []TestData{{
		args: []string{},
		expected: Args{
			logLevel: "error",
		},
	}, {
		args: []string{"--log-level=debug"},
		expected: Args{
			logLevel: "debug",
		},
	}, {
		args: []string{"--log-level=info"},
		expected: Args{
			logLevel: "info",
		},
	}, {
		args: []string{"--log-level=error"},
		expected: Args{
			logLevel: "error",
		},
	}, {
		args: []string{"--config-key", "xyz"},
		expected: Args{
			logLevel:  "error",
			configKey: "xyz",
		},
	}, {
		args: []string{"--min-java-version", "11"},
		expected: Args{
			logLevel:       "error",
			minJavaVersion: 11,
		},
	}, {
		args: []string{"--max-java-version", "17"},
		expected: Args{
			logLevel:       "error",
			maxJavaVersion: 17,
		},
	}, {
		args: []string{"--vendors", "Eclipse Adoptium"},
		expected: Args{
			logLevel: "error",
			vendors:  []string{"Eclipse Adoptium"},
		},
	}, {
		args: []string{"--vendors", "Eclipse Adoptium", "--vendors", "GraalVM Community"},
		expected: Args{
			logLevel: "error",
			vendors:  []string{"Eclipse Adoptium", "GraalVM Community"},
		},
	}}
	for _, data := range data {
		actual, err := ParseArgs(data.args)
		description := fmt.Sprintf("ParseArgs(%#v)", data.args)
		assertErrorEquals(t, description, data.err, err)
		assertEquals(t, description, &data.expected, actual)
	}
}

func TestParseArgsErrors(t *testing.T) {
	type TestData struct {
		args []string
		err  string
	}
	data := []TestData{{
		args: []string{"--unknown-flag"},
		err:  "flag provided but not defined: -unknown-flag",
	}, {
		args: []string{"unresolved argument"},
		err:  "unresolved arguments: [unresolved argument]",
	}, {
		args: []string{"--log-level=xoxo"},
		err:  "invalid log level: 'xoxo'. Available levels are: debug, info, error",
	}}
	for _, data := range data {
		actual, err := ParseArgs(data.args)
		description := fmt.Sprintf("ParseArgs(%#v)", data.args)
		var nothing *Args
		assertEquals(t, description, nothing, actual)
		assertErrorContains(t, description, data.err, err)
	}
}
