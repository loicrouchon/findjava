//go:build standalone_linux
// +build standalone_linux

package linker

func init() {
	setConfigDir("/etc/findjava/")
	setCacheDir("~/.cache/findjava/")
	setMetadataExtractorDir("./metadata-extractor")
}
