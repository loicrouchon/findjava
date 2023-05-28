//go:build linux
// +build linux

package main

func init() {
	platform.ConfigDir = "/etc/jvm-finder/config.json"
	platform.MetadataExtractorDir = "/usr/share/jvm-finder/metadata-extractor/JvmMetadataExtractor.class"
	platform.CacheDir = "~/.cache/jvm-finder/cache.json"
}
