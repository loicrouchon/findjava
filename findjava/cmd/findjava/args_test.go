package main

import (
	"findjava/test"
	"fmt"
	"testing"
)

func TestParseArgs(t *testing.T) {
	type TestData struct {
		args     []string
		expected Args
		err      error
	}
	defaults := Args{
		logLevel:   "error",
		Programs:   []string{"java"},
		OutputMode: "binary",
	}
	data := []TestData{{
		args:     []string{},
		expected: defaults,
	}, {
		args: []string{"--log-level=debug"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "debug"
		}),
	}, {
		args: []string{"--log-level=info"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "info"
		}),
	}, {
		args: []string{"--log-level=error"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
		}),
	}, {
		args: []string{"--config-key", "xyz"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.ConfigKey = "xyz"
		}),
	}, {
		args: []string{"--min-java-version", "11"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.MinJavaVersion = 11
		}),
	}, {
		args: []string{"--max-java-version", "17"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.MaxJavaVersion = 17
		}),
	}, {
		args: []string{"--vendors", "Eclipse Adoptium"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.Vendors = []string{"Eclipse Adoptium"}
		}),
	}, {
		args: []string{"--vendors", "Eclipse Adoptium", "--vendors", "GraalVM Community"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.Vendors = []string{"Eclipse Adoptium", "GraalVM Community"}
		}),
	}, {
		args: []string{"--programs", "javac"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.Programs = []string{"javac"}
		}),
	}, {
		args: []string{"--programs", "java", "--programs", "javac", "--programs", "native-image", "--output-mode", "java.home"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.Programs = []string{"java", "javac", "native-image"}
			args.OutputMode = "java.home"
		}),
	}, {
		args:     []string{"--output-mode", "binary"},
		expected: defaults,
	}, {
		args: []string{"--output-mode", "java.home"},
		expected: patch(defaults, func(args *Args) {
			args.OutputMode = "java.home"
		}),
	}}
	for _, data := range data {
		actual, err := ParseArgs(data.args)
		description := fmt.Sprintf("ParseArgs(%#v)", data.args)
		test.AssertErrorEquals(t, description, data.err, err)
		test.AssertEquals(t, description, &data.expected, actual)
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
		err:  "invalid log level: \"xoxo\". Available levels are: debug, info, warn, error",
	}, {
		args: []string{"--programs", "java", "--programs", "javac", "--programs", "native-image"},
		err: "output mode \"binary\" cannot be used when multiple programs are requested. " +
			"Use \"java.home\" instead",
	}, {
		args: []string{"--output-mode=xoxo"},
		err:  "invalid output mode: \"xoxo\". Available values are: java.home, binary",
	}}
	for _, data := range data {
		actual, err := ParseArgs(data.args)
		description := fmt.Sprintf("ParseArgs(%#v)", data.args)
		var nothing *Args
		test.AssertEquals(t, description, nothing, actual)
		test.AssertErrorContains(t, description, data.err, err)
	}
}

func patch(args Args, patchFunc func(args *Args)) Args {
	patchFunc(&args)
	return args
}
