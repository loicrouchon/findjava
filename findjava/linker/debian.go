//go:build debian
// +build debian

package linker

func init() {
	setConfigDir("/etc/findjava/")
	setCacheDir("~/.cache/findjava/")
	setMetadataExtractorDir("/usr/share/findjava/metadata-extractor")
}
