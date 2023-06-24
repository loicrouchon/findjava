package main

import (
	"bytes"
	. "findjava/internal/jvm"
	"findjava/internal/log"
	"findjava/internal/utils"
	"flag"
	"fmt"
	"os"
)

const outputModeBinary = "binary"
const outputModeJavaHome = "java.home"

type Args struct {
	version        bool
	logLevel       string
	ConfigKey      string
	MinJavaVersion uint
	MaxJavaVersion uint
	Vendors        utils.List
	Programs       utils.List
	OutputMode     string
}

func ParseArgs(commandArgs []string) (*Args, error) {
	args := Args{}
	cmd := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	output := bytes.NewBufferString("")
	cmd.SetOutput(output)
	cmd.Usage = func() {
		output.WriteString("Usage: findjava [OPTIONS]\n\n")
		output.WriteString("Options:\n")
		cmd.PrintDefaults()
	}
	cmd.BoolVar(&args.version, "version", false, "Displays the version")
	cmd.StringVar(&args.logLevel, "log-level", "error",
		"The log level which is one of: debug, info, warn, error. Defaults to error")
	cmd.StringVar(&args.ConfigKey, "config-key", "",
		"If specified, will look for an optional config.<KEY>.json to load before loading the default configuration")
	cmd.UintVar(&args.MinJavaVersion, "min-java-version", AllVersions,
		"The minimum (inclusive) Java Language Specification version the found JVMs should provide")
	cmd.UintVar(&args.MaxJavaVersion, "max-java-version", AllVersions,
		"The maximum (inclusive) Java Language Specification version the found JVMs should provide")
	cmd.Var(&args.Vendors, "vendors",
		"The vendors to filter on. If empty, no vendor filtering will be done")
	cmd.Var(&args.Programs, "programs",
		"The programs the JVM should provide in its \"${java.home}/bin\" directory. If empty, defaults to java")
	cmd.StringVar(&args.OutputMode, "output-mode", outputModeBinary,
		"The output mode of findjava. Possible values are \"java.home\" (the home directory of the selected JVM) "+
			"and \"binary\" (the path to the desired binary of the selected JVM). If not specified, defaults to binary")
	if err := cmd.Parse(commandArgs); err != nil {
		return nil, fmt.Errorf("%s\n%s", err, output)
	}
	if unresolvedArgs := cmd.Args(); len(unresolvedArgs) > 0 {
		cmd.Usage()
		return nil, fmt.Errorf("unresolved arguments: %v\n%s", unresolvedArgs, output)
	}
	if err := log.SetLogLevel(args.logLevel); err != nil {
		return nil, err
	}
	if len(args.Programs) == 0 {
		args.Programs = append(args.Programs, "java")
	}
	if err := validateOutputMode(args); err != nil {
		return nil, err
	}
	return &args, nil
}

func validateOutputMode(args Args) error {
	if args.OutputMode == outputModeJavaHome {
		return nil
	} else if args.OutputMode == outputModeBinary {
		if len(args.Programs) > 1 {
			return fmt.Errorf("output mode \"%s\" cannot be used when multiple programs are requested. "+
				"Use \"%s\" instead", args.OutputMode, outputModeJavaHome)
		}
		return nil
	} else {
		return fmt.Errorf("invalid output mode: \"%s\". Available values are: java.home, binary", args.OutputMode)
	}
}
