# findjava

```shell
#!/bin/sh
JAVA="$(findjava --min-java-version=11)"
"$JAVA" ... # my application command line
```

findjava is a tool whose goal is to find the best Java runtime for an application, allowing it to be distributed in a
simple and reliable way.

* Reliable, as it ensures the best Java executable matching the application's constraints will be selected for
  execution. It is not the end user's job to understand the technical runtime requirements of the application they are
  trying to run. Package managers do a great job of installing the required Java runtime, but it is mostly outside their
  scope to provide tooling for locating the appropriate runtime when executing the application, especially if multiple
  versions of the runtime are available and can be installed in parallel.
* Simple, as it unifies the way to locate the Java executable in application start scripts for all distribution
  channels.

It is not a goal of findjava to provide a way to install JVMs, as this is the responsibility of package managers.

* [Motivations \& problem statement](#motivations--problem-statement)
  * [The actors](#the-actors)
  * [The distribution channels](#the-distribution-channels)
* [Features](#features)
* [Usage](#usage)
  * [Arguments](#arguments)
* [Configuration](#configuration)
  * [JVM Discovery (files, directories, environment variables)](#jvm-discovery-files-directories-environment-variables)
  * [JVM filtering](#jvm-filtering)
  * [Multiple candidate JVMs found](#multiple-candidate-jvms-found)
* [Implementation Guidelines](#implementation-guidelines)
  * [For Standalone Packages (zip, tar.gz, ...)](#for-standalone-packages-zip-targz-)
  * [For packages managed by a package manager](#for-packages-managed-by-a-package-manager)
* [Installation](#installation)
  * [Ubuntu (23.04 and above)](#ubuntu-2304-and-above)
  * [Fedora (37 and above)](#fedora-37-and-above)
  * [Homebrew (macOS/Linux)](#homebrew-macoslinux)
* [Building the application](#building-the-application)

## Motivations & problem statement

Java has a fast release cycle. Every six months, a new version is released with various improvements (performance,
security, features) as well as deprecations and feature removals.

Distributing a Java application that relies on the JVM to be installed as a dependency by a package manager is a
difficult exercise. Even if the proper JVM is installed along with the program, there is no guarantee it will be the one
available in the end user's `$PATH`. To complicate matters further, there is also no strong guarantee regarding the path
where the JVM will be installed by the package manager. There could be various reasons for this, ranging from no
guarantee provided by the package manager to the JVM package possibly being a "virtual" package, which will be resolved
to a different JVM later on.

This creates a serious challenge for developers, package maintainers, and end users, as they may have multiple
applications installed, each requiring specific Java version(s).

There are two important dimensions to this problem: the actors and the distribution channels.

### The actors

There are three different kinds of actors:

| Actor                       | Description                                                        | Needs                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| --------------------------- | ------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Java application developers | Develop the application                                            | <ul><li>Want a single way to locate the Java executable (_simplicity_)</li><li>Want the located Java executable to match the application's runtime constraints (_reliability_)</li></ul>                                                                                                                                                                                                                                                                                                                                                                    |
| Package maintainers         | Responsible for the packaging of a particular distribution channel | <ul><li>Do not want to provide and maintain a mechanism for selecting JVMs based on criteria such as:<ul><li>It is very specific to Java due to the high release frequency of the JVM.</li><li>It is not really their responsibility to provide such a system (every package manager would need to implement it, probably with different solutions complicating the packaging for multiple package managers).</li></ul></li><li>Want to align the Java executable selection rules to match the package manager's dependency rules (_reliability_)</li></ul> |
| End users                   | Use the packaged application                                       | <ul><li>Want the application to always run properly without having to worry about the technical details of the application runtime (what is in the `$PATH`, what is the `$JAVA_HOME`, etc.). Ideally, end users shouldn't have to be aware that the application requires a Java runtime, and even less how to fix issues related to it (_reliability_)</li></ul>                                                                                                                                                                                            |

### The distribution channels

There is a potentially arbitrarily high number of possible distribution channels, but we can classify them into three
different categories: standalone distribution, package manager distribution, or battery-included distribution.

| Distribution     | Description                                                                                                                                                                                                                                                                                                                                                                            |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Standalone       | Typically, it is an archive (zip, tar.gz, etc.) that the end user would manually download from a particular source, such as GitHub releases or the editor's website.<p>The downloaded archive does not contain the Java runtime, and it is the end user's responsibility to install the proper runtime through any means (e.g., direct download, SDKMAN, a package manager, etc.).</p> |
| Package manager  | A package that, in addition to containing the application itself, also includes metadata expressing a dependency on a suitable Java runtime package. It is the package manager's responsibility to install a Java runtime that matches the dependency definition.<p>If no such dependency is expressed, the distribution will be considered a standalone one.</p>                      |
| Battery-included | The battery-included distribution packages its own Java runtime. It is not a goal of findjava to address battery-included distributions, as the application already knows exactly where its Java runtime is located.                                                                                                                                                                   |

## Features

* JVM discovery: Scans a list of directories, files, and environment variables to find installed JVMs according to
  defined rules.
* JVM metadata extraction: Analyzes each JVM to extract its relevant metadata.
* JVM filtering: Filters based on minimum/maximum Java specification version, vendors, and programs (java, javac,
  native-image, etc.).
* Output mode: Provides the path desired binary of the selected JVM or the path its `java.home`.
* Configurable at the system level: JVM discovery and filtering can be configured at the system level, giving control to
  package managers.
* Configurable at the application level: Overrides are possible at the application level, allowing for exceptions.
* Caching: JVM metadata are automatically cached and invalidated when findjava detects that a JVM has been updated,
  deleted, or added.

## Usage

To use findjava, you need to call it from your start script or command line with the following arguments:

```shell
JAVA="$(findjava --min-java-version=11)"
"$JAVA" ...
```

### Arguments

* `--min-java-version <version>`: The minimum version of the Java specification required to run the application. If
  `--max-java-version` is specified, it defaults to `0`, meaning no minimum version filtering. If both
  `--min-java-version` and `--max-java-version` are not specified, it falls back on the configuration.
* `--max-java-version <version>`: The maximum version of the Java specification required to run the application. If
  `--min-java-version` is specified, it defaults to `0`, meaning no maximum version filtering. If both
  `--min-java-version` and `--max-java-version` are not specified, it falls back on the configuration.
* `--vendors <vendor>`: (repeatable) A list of JVM vendors to choose from. If specified, findjava will only consider
  JVMs from these vendors. If not specified, no vendor filtering will occur.
* `--programs <program>`: (repeatable) A list of programs that the JVM must provide in its `$JAVA_HOME/bin` directory.
  If more than one program is provided, the output will automatically switch to `java.home` mode. If not specified, it
  defaults to `java`.
* `--output-mode <output-mode>`: The output mode of findjava. Possible values are `java.home` (the `java.home` directory
  of the selected JVM) and `binary` (the path to the desired binary of the selected JVM). If not specified, it defaults
  to `binary`.

> Java specification versions can be specified in a simplified way as integers (e.g., 1, 2, 8, 20). findjava will
> recognize that versions 1.8 and 8 are equivalent.

## Configuration

> _**WORK IN PROGRESS**_

### JVM Discovery (files, directories, environment variables)

The JVM discovery is driven by the `jvmLookupPaths` configuration property. It will scan the content of each path in
this property to discover JVMs.

Each of these paths must be either absolute or relative to the user home (`~`). Path processing will be performed as
follows:

1. Resolve environment variables (e.g. `$JAVA_HOME`, `$JAVA_HOME/bin/java`, etc.).
2. Resolve the user's home directory `~` (e.g. `~/.sdkman/candidates/java`).

JVMs will be discovered for a given path in the following use cases:

* The path points to a file (after resolving symbolic links) that is executable.
  * Examples:
    * `/usr/bin/java`
    * `$JAVA_HOME/bin/java`
* The path points to a directory that contains (after resolving symbolic links) a `bin/java` executable.
  * Examples:
    * `$JAVA_HOME`
    * `$GRAALVM_HOME`
  * If no `bin/java` executable is found, all direct subdirectories will be checked for `<subdirectory>/bin/java`
      executables.
    * Examples:
      * `/usr/lib/jvm`
      * `~/.sdkman/candidates/java`
      * `/System/Volumes/Data/Library/Java/JavaVirtualMachines`
    * This will not recurse into subdirectories of subdirectories.

If no configuration for `jvmLookupPaths` is defined, sensible defaults depending on the operating system will be used
for the lookup. The defaults are specified below:

* > _**WORK IN PROGRESS**_

### JVM filtering

The filtering is split into two steps:

1. Applying the strong filtering constraints (i.e., specified on the command line).
2. Applying the weak filtering constraints (i.e., coming from the configuration).

The reasoning is as follows: The startup script calling findjava is the most knowledgeable about the requirements of the
program it needs to run. Therefore, constraints expressed as arguments when calling findjava are considered strong.

On the other hand, system configuration will be considered as recommendations, and findjava will try to fulfill those as
much as possible.

If strong constraints can be satisfied but not the recommendations from the system configuration, findjava will ignore
the recommendations and select a JVM based solely on the strong constraints.

> **Note:** If no `--min-java-version`/`--max-java-version` is specified on the command line, findjava will not consider
> having strong recommendations. In this case, if system recommendations cannot be fulfilled, findjava will fail.
> _This behavior might be revisited in the near future_.

> **Recommendation:** It is recommended to always specify the `--min-java-version` option.

### Multiple candidate JVMs found

In case multiple JVMs are found to match the filtering criteria, an election process will be initiated to select which
one of these shall be used.

This process will return the JVM implementing the highest `java.specification.version`. If multiple JVMs implement the
same `java.specification.version`, one will be selected. This selection process is not currently specified nor
deterministic. Future versions of findjava might provide rules for preferred JVM selection in such cases.

## Implementation Guidelines

### For Standalone Packages (zip, tar.gz, ...)

> **Reminder:** For standalone packages, the Java runtime must be installed by the end user beforehand.

To select the desired Java runtime, the application startup shell script needs to call findjava and specify at least the
minimum required Java version to run the application.

This can be done as shown below:

```shell
JAVA="$(findjava --min-java-version=11)"
"$JAVA" ...
```

findjava will look up Java runtimes according to its configuration located in `<FINDJAVA_CONFIG_DIR>/config.conf`, where
`<FINDJAVA_CONFIG_DIR>` is the path to the directory holding the findjava configuration (e.g., `/etc/findjava` under
Linux).

Additional options can be added to narrow down the Java runtime resolution process, such as specifying the maximum Java
version. For a complete list of options, refer to the [usage section](#usage).

```shell
JAVA="$(findjava --min-java-version=11 --max-java-version=17)"
"$JAVA" ...
```

### For packages managed by a package manager

When integrating a package using findjava to locate a Java runtime into a package manager, it is essential to align
findjava rules with the package's dependency metadata.

For example, on Debian/Ubuntu systems, if findjava specifies the rule `--min-java-version=17`, then the package should
express a dependency on `java17-runtime-headless` or any other compatible Java package.

However, this might not be enough. Even though backward compatibility is a strong aspect of the Java language and
virtual machine, features can be deprecated, disabled, and eventually removed. What if, in the future, the end user
installs an additional Java runtime that is incompatible with the application?

One solution would be for every application using findjava to also specify the maximum Java version. However, this is a
lot of work for developers and package maintainers, especially when introducing new Java runtimes into the distribution.

The recommended approach is to specify these settings in a single place for all packages installed through the package
manager and only override them at the application level as an exception. The solution is called _alternative
configurations_ which findjava supports through its `--config-key` command-line parameter.

For example, on Debian, the findjava package could define two configuration files:

* `/etc/findjava/config.conf`
* `/etc/findjava/config.dpkg.conf`

The first file, `/etc/findjava/config.conf`, would be used for all calls to findjava that do not specify the
`--config-key` command-line parameter.

The second file, `/etc/findjava/config.dpkg.conf`, is an alternative configuration which would be used when
`--config-key=dpkg` is specified.

Consider the following configuration files:

For `/etc/findjava/config.conf`:

```properties
jvm.lookup.paths=$JAVA_HOME/bin/java, /usr/bin/java, /usr/lib/jvm, ~/.sdkman/candidates/java
```

For `/etc/findjava/config.dpkg.conf`:

```properties
jvm.lookup.paths=/usr/lib/jvm
java.specification.version.min=8
java.specification.version.max=21
```

Calls to findjava will be able to find any Java runtime installed in various locations, subject only to the constraints
specified by the start script when calling findjava.

In contrast, calls to `findjava --config-key=dpkg` will only find Java runtimes installed in `/usr/lib/jvm`, where the
Java version is between 8 (inclusive) and 21 (inclusive).

This provides more control to package managers, ensuring that a package installed through the package manager will have
findjava rules in sync with the package manager's capabilities.

## Installation

The goal is for findjava to be available in as many package managers for Linux, macOS, and Windows as possible, so that
you can depend on it when packaging your application for those package managers.

Currently, we are not fully there yet, but we are making progress.

### Ubuntu (23.04 and above)

```shell
sudo apt install software-properties-common
sudo add-apt-repository ppa:loicrouchon/symly
sudo apt update
sudo apt install findjava
```

### Fedora (37 and above)

```shell
sudo dnf install 'dnf-command(copr)'
sudo dnf copr enable loicrouchon/symly
sudo dnf install findjava
```

### Homebrew (macOS/Linux)

```shell
brew install loicrouchon/symly/findjava
```

## Building the application

To build the application, the following dependencies are required:

* Go (>= 1.15): To build the application. The Go version might be relaxed in the future.
* A JDK (>= 9): To build the JVM metadata extraction.
* `make`: For build automation.

The application can then be built with:

```shell
make
```

> Be aware of the default configuration when building the application.
> By default, a development build will be created (configuration from [main.go](findjava/cmd/findjava/main.go)).
> This can be changed to one of the following tags: [darwin](findjava/linker/standalone_macos.go),
> [standalone_linux](findjava/linker/standalone_linux.go), [debian](findjava/linker/debian.go).
>
> To do so, set the `GO_TAGS` environment variable before calling `make` as follows: `GO_TAGS="-tags <TAG>"`.
>
> Example to build a binary customized for debian:
>
> ```shell
> GO_TAGS="-tags debian" make
> ```
>
> Another alternative is to override individual settings using go build -ldflags by setting the `GO_LD_FLAGS`
> environment variable before calling `make`. For example:
>
> ```shell
> GO_LD_FLAGS="-X 'findjava/linker.MetadataExtractorDir=/usr/share/findjava/metadata-extractor'" make
> ```
>
> Note that both approaches could be combined. The following example will use the `standalone_linux` configuration
> and patch the `linker.MetadataExtractorDir` variable:
>
> ```shell
> GO_TAGS="-tags standalone_linux" GO_LD_FLAGS="-X 'findjava/linker.MetadataExtractorDir=/usr/share/findjava/metadata-extractor'" make
> ```
>
> Full documentation about the overridable variables is available in the [linker package documentation](findjava/linker/doc.go).

Once built, you can run it with:

```shell
./build/dist/findjava
```
