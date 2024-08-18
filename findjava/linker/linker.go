package linker

const defaultConfigDir = "../"
const defaultCacheDir = "../"
const defaultMetadataExtractorDir = "./metadata-extractor/"

// ConfigDir is the path to the directory holding the findjava configuration.
// Can be absolute, relative to the user's home (starts with ~),
// or relative to the findjava binary's directory.
var ConfigDir = defaultConfigDir

// CacheDir is the path to the directory in which findjava will cache the JVM metadata.
// Can be absolute, relative to the user's home (starts with ~),
// or relative to the findjava binary's directory.
var CacheDir = defaultCacheDir

// MetadataExtractorDir is the path to the directory holding the JVM metadata extractor.
// Can be absolute, relative to the user's home (starts with ~),
// or relative to the findjava binary's directory.
var MetadataExtractorDir = defaultMetadataExtractorDir

func setConfigDir(value string) {
	setIfNotOverridden(&ConfigDir, defaultConfigDir, value)
}

func setCacheDir(value string) {
	setIfNotOverridden(&CacheDir, defaultCacheDir, value)
}

func setMetadataExtractorDir(value string) {
	setIfNotOverridden(&MetadataExtractorDir, defaultMetadataExtractorDir, value)
}

func setIfNotOverridden(variable *string, defaultValue string, newValue string) {
	if *variable == defaultValue {
		*variable = newValue
	}
}
