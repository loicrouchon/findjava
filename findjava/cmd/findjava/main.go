package main

import (
	"findjava/internal/config"
	"findjava/internal/console"
	"findjava/internal/discovery"
	"findjava/internal/jvm"
	"findjava/internal/log"
	"findjava/internal/rules"
	"findjava/internal/selection"
	"findjava/linker"
	"fmt"
	"os"
	"path/filepath"
)

var Version = "dev"

func main() {
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		log.Die(err)
	}
	platform := config.Platform{
		ConfigDir:            linker.ConfigDir,
		CacheDir:             linker.CacheDir,
		MetadataExtractorDir: linker.MetadataExtractorDir,
	}
	if args.version {
		console.Writer.Printf("findjava %s\n", Version)
		_ = platform.Resolve() // Prints platform information at debug level
		return
	}
	if args == nil {
		os.Exit(0)
	}
	cfg, err := platform.LoadConfig(args.ConfigKey)
	if err != nil {
		log.Die(err)
	}
	javaExecutables, err := discovery.FindAllJavaExecutables(&cfg.JvmsLookupPaths)
	if err != nil {
		log.Die(err)
	}
	metaDataFetcher := &jvm.MetadataReader{Classpath: cfg.JvmsMetadataExtractorPath}
	jvmInfos, err := jvm.LoadJvmsInfos(metaDataFetcher, cfg.JvmsMetadataCachePath, &javaExecutables)
	if err != nil {
		log.Die(err)
	}
	rules := rules.SelectionRules(cfg, args.MinJavaVersion, args.MaxJavaVersion, args.Vendors, args.Programs)
	if jvms := selection.Select(rules, &jvmInfos); len(jvms) > 0 {
		jvm := jvms[0]
		selection.LogJvmList("[SELECTED]", jvms[0:1])
		if err := processOutput(args, jvm); err != nil {
			log.Die(err)
		}
	} else {
		log.Die(fmt.Errorf("unable to find a JVM matching requirements %s", rules))
	}
}

func processOutput(args *Args, jvm jvm.Jvm) error {
	if args.OutputMode == outputModeJavaHome {
		console.Writer.Printf("%s\n", jvm.JavaHome)
		return nil
	}
	if args.OutputMode == outputModeBinary {
		for _, program := range args.Programs {
			console.Writer.Printf("%s\n", filepath.Join(jvm.JavaHome, "bin", program))
		}
		return nil
	}
	return fmt.Errorf("unsupported output-mode \"%s\"", args.OutputMode)
}
