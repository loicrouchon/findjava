# How to publish source packages for Ubuntu

1. Decide which version will be the next version
2. Update the version in [version.txt](../../version.txt)
3. Update the version in [debian/changelog](debian/changelog).
   Version should be suffixed by `-$BUILD_NUMBER` and an optional `~$FLAVOR`.
   Also ensure the distribution version (jammy, lunar, ...) is set according to the publication target
4. Commit everything (`main` branch) and tag with `v$VERSION` (i.e. `v1.2.3` for version `1.2.3`)
5. `git push && git push --tags`
6. Initiate the source package build environment with `./docker-build-env.sh`
   This requires that the public and private gpg for signing the packages to be available in the [gpg](gpg) directory
7. Build the source package with `./package.sh $VERSION`
   This will be a binary package as well as a source package
8. Once satisfied, publish the source package to launchpad with `./publish.sh`
