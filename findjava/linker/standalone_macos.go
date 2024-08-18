//go:build standalone_macos
// +build standalone_macos

package linker

func init() {
	setConfigDir("/etc/findjava/")
	setCacheDir("~/.cache/findjava/")
	setMetadataExtractorDir("./metadata-extractor")
}
