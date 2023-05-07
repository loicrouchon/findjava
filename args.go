package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

const outputModeBinary = "binary"
const outputModeJavaHome = "java.home"

type Args struct {
	logLevel       string
	configKey      string
	minJavaVersion uint
	maxJavaVersion uint
	vendors        list
	programs       list
	outputMode     string
}

type list []string

func (i *list) String() string {
	return "[" + strings.Join(*i, ", ") + "]"
}

func (i *list) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func ParseArgs(commandArgs []string) (*Args, error) {
	args := Args{}
	cmd := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	output := bytes.NewBufferString("")
	cmd.SetOutput(output)
	cmd.StringVar(&args.logLevel, "log-level", "error",
		"The log level which is one of: debug, info, warn, error. Defaults to error")
	cmd.StringVar(&args.configKey, "config-key", "",
		"If specified, will look for an optional config.<KEY>.json to load before loading the default configuration")
	cmd.UintVar(&args.minJavaVersion, "min-java-version", allVersions,
		"The minimum (inclusive) Java Language Specification version the found JVMs should provide")
	cmd.UintVar(&args.maxJavaVersion, "max-java-version", allVersions,
		"The maximum (inclusive) Java Language Specification version the found JVMs should provide")
	cmd.Var(&args.vendors, "vendors",
		"The vendors to filter on. If empty, no vendor filtering will be done")
	cmd.Var(&args.programs, "programs",
		"The programs the JVM should provide in its \"${java.home}/bin\" directory. If empty, defaults to java")
	cmd.StringVar(&args.outputMode, "output-mode", outputModeBinary,
		"The output mode of jvm-finder. Possible values are \"java.home\" (the home directory of the selected JVM) "+
			"and \"binary\" (the path to the desired binary of the selected JVM). If not specified, defaults to binary")
	if err := cmd.Parse(commandArgs); err != nil {
		return nil, fmt.Errorf("%s\n%s", err, output)
	}
	if unresolvedArgs := cmd.Args(); len(unresolvedArgs) > 0 {
		cmd.Usage()
		return nil, fmt.Errorf("unresolved arguments: %v\n%s", unresolvedArgs, output)
	}
	if err := setLogLevel(args.logLevel); err != nil {
		return nil, err
	}
	if len(args.programs) == 0 {
		args.programs = append(args.programs, "java")
	}
	if err := validateOutputMode(args); err != nil {
		return nil, err
	}
	return &args, nil
}

func validateOutputMode(args Args) error {
	if args.outputMode == outputModeJavaHome {
		return nil
	} else if args.outputMode == outputModeBinary {
		if len(args.programs) > 1 {
			return fmt.Errorf("output mode \"%s\" cannot be used when multiple programs are requested. "+
				"Use \"%s\" instead", args.outputMode, outputModeJavaHome)
		}
		return nil
	} else {
		return fmt.Errorf("invalid output mode: \"%s\". Available values are: java.home, binary", args.outputMode)
	}
}
