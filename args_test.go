package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	type TestData struct {
		args     []string
		expected Args
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
		console := setTestConsole()
		actual := ParseArgs(data.args)
		if !reflect.DeepEqual(actual, &data.expected) {
			t.Fatalf(`Expecting ParseArgs(%#v) == %#v but was %#v`,
				strings.Join(data.args, `", "`), data.expected, actual)
		}
		console.hasMessages(t, []string{}, []string{})
	}
}
