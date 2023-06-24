//go:build linux
// +build linux

package main

func init() {
	platform.ConfigDir = "/etc/findjava/"
	platform.MetadataExtractorDir = "./metadata-extractor"
	platform.CacheDir = "~/.cache/findjava/"
}
