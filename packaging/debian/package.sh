#!/usr/bin/env sh
set -ex

package_name="jvm-finder"
version="$1"
package_version_dir="${package_name}-${version}"
upstream_tarball="${package_name}_${version}.orig.tar.gz"
# TODO
#upstream_tarball_url="https://github.com/loicrouchon/jvm-finder/archive/refs/tags/v${version}.tar.gz"
upstream_tarball_url="https://github.com/loicrouchon/jvm-finder/archive/refs/heads/main.zip"
cur_dir="$(dirname "$(realpath "$0")")"

echo "Building debian source package ${package_version_dir}"

rm -rf "build"
mkdir -p "build/${package_name}"
cd "build/${package_name}" || exit 1

echo "Downloading upstream tarball from ${upstream_tarball_url}"
curl -sL "${upstream_tarball_url}" -o "${upstream_tarball}"

echo "Unpacking upstream tarball ${upstream_tarball}"
tar xzf "${upstream_tarball}"
rm -f "${upstream_tarball}"

echo "Repacking upstream tarball ${upstream_tarball} (get rid off root level directory)"
find . -maxdepth 1 -type d -name "${package_name}-*" -exec mv {} "${package_version_dir}" \;
cd "${package_version_dir}"
tar czf "../${upstream_tarball}" *

echo "Add debian dir"
cd "../${package_version_dir}" || exit 1
cp -R "${cur_dir}/debian" "debian"

if ! (cat debian/changelog | grep -E "^${package_name} " | head -n 1 | grep -q "${package_name} (${version}-"); then
  echo "Package $package_name does not have a changelog entry for version $version"
  exit 1
fi

echo "Configure environment:"
echo "Configure environment: maintainer name and email"
export DEBEMAIL="loic@loicrouchon.com"
export DEBFULLNAME="Loic Rouchon"
echo "Configure environment: GPG fingerprint"
GPG_KEY_FINGERPRINT="C3BB9448B16C971103E876BF3A091A0DF2799262"

echo "Configure environment: PPA URL"
PPA_URL="ppa:loicrouchon/jvm-finder"

echo "Build source package"
dpkg-buildpackage --sign-key="${GPG_KEY_FINGERPRINT}" --build=source

cd ..
ls -l
echo ".dsc content:"
cat ./*.dsc
echo ".buildinfo content:"
cat ./*.buildinfo
echo ".changes content:"
cat ./*.changes

dput "${PPA_URL}" ./*.changes
