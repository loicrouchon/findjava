//go:build darwin
// +build darwin

package main

func init() {
	platform.ConfigDir = "/etc/findjava/"
	platform.MetadataExtractorDir = "./metadata-extractor"
	platform.CacheDir = "~/.cache/findjava/"
}
