#!/usr/bin/env sh
set -e
#set -x

cur_dir="$(dirname "$(realpath "$0")")"

if [ "$#" -ne 1 ]; then
  echo "Error: Invalid number of arguments. Expected 1 argument."
  echo "Usage $0 VERSION"
  exit 1
fi

version="$1"
if [ -z "$package_name" ]; then
    . "${cur_dir}/../common/config.sh"
fi

upstream_tarball_name="${package_name}_${version}.orig.tar.gz"

repack() {
    echo "Repacking upstream tarball ${upstream_tarball_name} (get rid off root level directory)"
    find . -maxdepth 1 -type d -name "${package_name}-*" -exec mv {} "${package_version_dir}" \; || true
    tar czf "${upstream_tarball_name}" "${package_version_dir}"
    echo "Add debian dir"
    cp -R "${cur_dir}/debian" "${package_version_dir}/debian"
}

check() {
    if ! (cat "${package_version_dir}/debian/changelog" | grep -E "^${package_name} " | head -n 1 | grep -q "${package_name} (${version}-"); then
      echo "Package $package_name does not have a changelog entry for version $version"
      exit 1
    fi
}

build() {
    echo "Building debian binary package ${package_version_dir} (test mode)"
    (cd "${package_version_dir}" && dpkg-buildpackage --sign-key="${GPG_KEY_FINGERPRINT}" --build=binary)
    echo "Building debian source package ${package_version_dir} (for upload)"
    (cd "${package_version_dir}" && dpkg-buildpackage --sign-key="${GPG_KEY_FINGERPRINT}" --build=source)

    ls -l
    echo ".dsc content:"
    cat ./*.dsc
    echo ".buildinfo content:"
    cat ./*.buildinfo
    echo ".changes content:"
    cat ./*.changes
}

"${cur_dir}/../common/clean.sh" "build" "${package_name}"
cd "build/${package_name}" || exit 1
"${cur_dir}/../common/download-sources.sh" "${version}"
repack
check
build
