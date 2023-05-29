# How to publish RPM Spec for Fedora

1. Decide which version will be the next version
2. Update the version in [version.txt](../../version.txt)
3. Commit everything (`main` branch) and tag with `v$VERSION` (i.e. `v1.2.3` for version `1.2.3`)
4. `git push && git push --tags`
5. Build the RPM spec with `./package.sh $VERSION`
6. Once satisfied, publish the RPM spec package to Fedora COPR with `./publish.sh`
