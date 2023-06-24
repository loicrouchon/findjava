//go:build linux
// +build linux

package main

func init() {
	platform.ConfigDir = "/etc/findjvm/"
	platform.MetadataExtractorDir = "./metadata-extractor"
	platform.CacheDir = "~/.cache/findjvm/"
}
