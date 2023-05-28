//go:build linux
// +build linux

package main

func init() {
	platform.ConfigDir = "/etc/jvm-finder/config.json"
	platform.MetadataExtractorDir = "./metadata-extractor"
	platform.CacheDir = "~/.cache/jvm-finder/cache.json"
}
