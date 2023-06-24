# findjava

`findjava` is a tool which goal is to find the best JVM to run your java program according to your program's specific
constraints.

## Motivations

Java has a fast release cycle. Every 6 months a new version is released, with various improvements (performance,
security, features) as well as deprecations and feature removal.

Distributing a Java application which will rely on the JVM to be installed as dependency by a package manager is a
difficult exercise.
Even if the proper JVM is installed along the program, there is no guarantee it will be the one available in the
user's `$PATH`.
To complicate things even more, there is also no strong guarantee on the path where the JVM will be installed by the
package manager. There could be various reasons for this ranging from no guarantee provided by the package manager or
the JVM package might be a virtual package which will be resolved to a different JVM later on.

This creates a serious challenge for developers, packages maintainers and end-users, as they may have
multiple Java applications installed, each requiring a specific Java version.

There are three actors involved in this problem:

- The application developers: who want an easy way to ensure that their program runs with the proper JVM.
- The packages' maintainers: who do not want to provide and maintain a mechanism for selecting JVMs based on criteria
  as:
    - it is very specific to Java due to the high release frequency of the JVM
    - it is not really their responsibility to provide such a system (every package manager would need to implement it,
      probably with different solutions complicating the packaging for multiple package managers).
    - However, they want to ensure that all programs installed via their package manager will work and therefore may
      want to provide some default JVM selection metadata
- End-users, who do not want to know that the program they are running needs a JVM, and even less that the current JVM
  in the PATH is not the right one for that given program.

## Goal

The goal of findjava is to provide a solution to this problem while having a minimal impact on programs startup time.
It is not a goal of findjava to provide a way to install JVMs, as this is still the responsibility of package
managers.

## Features

- JVM discovery: Scans a list of directories, files, and environment variables to find installed JVMs
- JVM metadata extraction: Analyzes each JVM to extract its relevant metadata
- JVM filtering: minimum/maximum java specification version, vendors, programs (java, javac, native-image, ...)
- Outputs mode: path to the java.home of the selected JVM or directly to the desired binary
- Configurable at system level: JVM discovery and filtering can be configured at system level giving control to package
  managers
- Configurable at application level: Overrides are possible at application level allowing for exceptions
- Caching: JVMs metadata are automatically cached and invalidated when it detects a JVM is updated, deleted, or added.

## Installation

The goal is for findjava to be available in as many package managers for Linux, macOS and Windows as possible so that you can depend on it when packaging your application in those package managers.

For the moment, we are not there, but we're progressing.

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

> _**WORK IN PROGRESS**_: To be usable, findjava needs to be available in package managers

## Usage

To use findjava, you need to call it from your start script or command line with the following arguments:

```shell
JAVA="$(findjava --min-java-version=11)"
"$JAVA" ...
```

### Arguments

- `--min-java-version <version>`: the minimum version of the java specification required to run the application.
  If `--max-java-version` is specified, defaults to `0` meaning no minimum version filtering.If
  both `--min-java-version` and `--max-java-version` are not specified, falls back on the configuration.
- `--max-java-version <version>`: the maximum version of the java specification required to run the application.
  If `--min-java-version` is specified, defaults to `0` meaning no maximum version filtering.If
  both `--min-java-version` and `--max-java-version` are not specified, falls back on the configuration.
- `--vendors <vendor>`: (repeatable) a list of JVM vendors to choose from. If specified, findjava will only consider
  JVMs from this vendor. If not specified, not vendor filtering will happen
- `--programs <program>`: (repeatable) a list of programs the JVM must provide in their `$JAVA_HOME/bin` directory. If
  more than one program is provided, the output will automatically be in `java.home` mode. If not specified, defaults
  to `java`.
- `--output-mode <output-mode>`: the output mode of findjava. Possible values are `java.home` (the `java.home`
  directory of
  the selected JVM) and `binary` (the path to the desired binary of the selected JVM). If not specified, defaults
  to `binary`.

> Java specification versions can be specified in a simplified way as integers 1, 2, 8, 20.
> findjava will know that 1.8 and 8 are the same version.

## Configuration

> _**WORK IN PROGRESS**_

### JVM Discovery (files, directories, environment variables)

The JVM discovery is driven by the `jvmLookupPaths` configuration property.
It will scan the content each path in this property to discover JVMs.

Each of those path must be either absolute, relative to the user home (`~`).
Path processing will be performed as follows:

1. Resolve environment variables (`$JAVA_HOME`, `$JAVA_HOME/bin/java`)
2. Resolve user's home directory `~`: `~/.sdkman/candidates/java`

JVMs will be discovered for a given path in the following use cases:

* The path is pointing to a file (after symbolic links resolution) which is executable
    * Examples:
        * `/usr/bin/java`
        * `$JAVA_HOME/bin/java`
* The path is pointing to a directory which contains (after symbolic links resolution) a `bin/java` executable
    * Examples:
        * `$JAVA_HOME`
        * `$GRAALVM_HOME`
    * If no `bin/java` executable is found, all direct subdirectories will be checked for `<subdirectory>/bin/java`
      executables
        * Examples:
            * `/usr/lib/jvm`
            * `~/.sdkman/candidates/java`
            * `/System/Volumes/Data/Library/Java/JavaVirtualMachines`
        * This will not recurse to subdirectories' subdirectories

If no configuration for `jvmLookupPaths` is defined, sensible defaults depending on the operating system will be used
for the lookup. The defaults are specified below:

* > _**WORK IN PROGRESS**_

### JVM filtering

The filtering is split in two steps:

1. Applying the strong filtering constraints (i.e. specified on the command line)
2. Applying the weak filtering constraints (i.e. coming from the configuration)

The reasoning is the following: The startup script calling findjava is the most knowledgeable about the requirements
of the program it will need to run. Therefore, constraints expressed as arguments when calling `findjava` are
considered strong.

On the other side, system configuration will be considered as recommendations and `findjava` will try to fulfill those
as much as possible.

If strong constraints can be satisfied but not the recommendations from the system configuration, `findjava` will
ignore the recommendations and select a JVM only based on the strong constraints.

> **Note**: If no `--min-java-version`/`--max-java-version` is specified on the command line, `findjava` will not consider
> to have strong recommendations. In this case, if system recommendations cannot be fulfilled, it will fail.
> _This behavior might be revisited in the near future_.
>
> **Recommendation**: It is recommended to always specify the `--min-java-version`

### Multiple candidate JVMs found

In case multiple JVMs are found to be matching the filtering criteria, an election process will be started to select
which one of those shall be used.

This process will return the JVM implementing the highest `java.specification.version`.
If multiple JVMs implement the same `java.specification.version` one will be selected. This selection process is not
specified at the moment nor deterministic. Future versions of `findjava` might provide rules for preferred JVM
selection in such a case.

## Building the application

To build the application, the following dependencies are required:

* Go (>= 1.15): to build the application. The go version might be relaxed in the future
* A JDK (>= 9): to build the JVM metadata extraction
* `make`

The application can then be built with:

```shell
make
```

Once built, you can run it with

```shell
./build/go/findjava
```
