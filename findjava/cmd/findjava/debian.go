//go:build debian
// +build debian

package main

func init() {
	platform.ConfigDir = "/etc/findjava/"
	platform.MetadataExtractorDir = "/usr/share/findjava/metadata-extractor"
	platform.CacheDir = "~/.cache/findjava/"
}
