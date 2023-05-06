# jvm-finder

`jvm-finder` is a tool which goal is to find the best JVM to run your java program according to your program's specific
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

The goal of jvm-finder is to provide a solution to this problem while having a minimal impact on programs startup time.
It is not a goal of jvm-finder to provide a way to install JVMs, as this is still the responsibility of package
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

> _**WORK IN PROGRESS**_: To be usable, jvm-finder needs to be available in package managers

## Usage

To use jvm-finder, you need to call it from your start script or command line with the following arguments:

```shell
JAVA="$(jvm-finder --min-java-version=11)"
"$JAVA" ...
```

### Arguments

- `--min-java-version <version>`: the minimum version of the java specification required to run the application.
  If `--max-java-version` is specified, defaults to `0` meaning no minimum version filtering.If
  both `--min-java-version` and `--max-java-version` are not specified, falls back on the configuration.
- `--max-java-version <version>`: the maximum version of the java specification required to run the application.
  If `--min-java-version` is specified, defaults to `0` meaning no maximum version filtering.If
  both `--min-java-version` and `--max-java-version` are not specified, falls back on the configuration.
- `--vendors <vendor>`: (repeatable) a list of JVM vendors to choose from. If specified, jvm-finder will only consider
  JVMs from this vendor. If not specified, not vendor filtering will happen
- **[TO BE IMPLEMENTED]** `--program <program>`: (repeatable) a list of programs the JVM must provide in their `$JAVA_HOME/bin` directory. If
  more than one program is provided, the output will automatically be in `java.home` mode. If not specified, defaults
  to `java`.
- **[TO BE IMPLEMENTED]** `--output <output-mode>`: the output mode of jvm-finder. Possible values are `java.home` (the `java.home` directory of
  the selected JVM) and `binary` (the path to the desired binary of the selected JVM). If not specified, defaults
  to `binary`.

> Java specification versions can be specified in a simplified way as integers 1, 2, 8, 20.
> jvm-finder will know that 1.8 and 8 are the same version.

## Configuration

> _**WORK IN PROGRESS**_
 
### JVM Discovery (files, directories, environment variables)

### JVM filtering
