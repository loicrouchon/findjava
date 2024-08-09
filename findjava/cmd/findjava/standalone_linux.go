//go:build standalone_linux
// +build standalone_linux

package main

func init() {
	platform.ConfigDir = "/etc/findjava/"
	platform.MetadataExtractorDir = "./metadata-extractor"
	platform.CacheDir = "~/.cache/findjava/"
}
