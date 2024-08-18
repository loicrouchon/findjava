/*
Package linker defines a set of configurable variables which can be overridden at build time through go ldflags
as well as some pre-written build configurations.

The following variables can be overridden at build time thanks to go ldflags:

  - [ConfigDir]: the path to the directory holding the findjava configuration.
  - [CacheDir]: the path to the directory in which findjava will cache the JVM metadata.
  - [MetadataExtractorDir]: the path to the directory holding the JVM metadata extractor.

To override one variable, set a `-X` variable definition via go build -ldflags:

	-X 'findjava/linker.<NAME>=<VALUE>'

Where `<NAME>` is the variable name and `<VALUE>` the value to inject at build time.

The pre-written configurations are the following:

  - standalone_linux: configures the variable for a Linux environment where the metadata extractor
    should be located next to the binary.
  - standalone_macos: configures the variable for a macOS environment where the metadata extractor
    should be located next to the binary.
  - debian: configures the variable for honouring Debian packaging rules.

To activate a configuration, set a go build tag with the configuration name when building.

Note: it is is possible to activate a configuration and then override individual values via ldflags.

If none of those configuration is specified at build time, the built binary will be a development one.
See [defaultConfigDir], [defaultCacheDir] and [defaultMetadataExtractorDir] for the development values.
*/
package linker
