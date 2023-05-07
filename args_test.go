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
	defaults := Args{
		logLevel:   "error",
		programs:   []string{"java"},
		outputMode: "binary",
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
			args.configKey = "xyz"
		}),
	}, {
		args: []string{"--min-java-version", "11"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.minJavaVersion = 11
		}),
	}, {
		args: []string{"--max-java-version", "17"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.maxJavaVersion = 17
		}),
	}, {
		args: []string{"--vendors", "Eclipse Adoptium"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.vendors = []string{"Eclipse Adoptium"}
		}),
	}, {
		args: []string{"--vendors", "Eclipse Adoptium", "--vendors", "GraalVM Community"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.vendors = []string{"Eclipse Adoptium", "GraalVM Community"}
		}),
	}, {
		args: []string{"--programs", "javac"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.programs = []string{"javac"}
		}),
	}, {
		args: []string{"--programs", "java", "--programs", "javac", "--programs", "native-image", "--output-mode", "java.home"},
		expected: patch(defaults, func(args *Args) {
			args.logLevel = "error"
			args.programs = []string{"java", "javac", "native-image"}
			args.outputMode = "java.home"
		}),
	}, {
		args:     []string{"--output-mode", "binary"},
		expected: defaults,
	}, {
		args: []string{"--output-mode", "java.home"},
		expected: patch(defaults, func(args *Args) {
			args.outputMode = "java.home"
		}),
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
		assertEquals(t, description, nothing, actual)
		assertErrorContains(t, description, data.err, err)
	}
}

func patch(args Args, patchFunc func(args *Args)) Args {
	patchFunc(&args)
	return args
}
