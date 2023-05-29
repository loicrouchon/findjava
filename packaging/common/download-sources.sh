#!/usr/bin/env sh
set -e

cur_dir="$(dirname "$(realpath "$0")")"

if [ "$#" -ne 1 ]; then
  echo "Error: Invalid number of arguments. Expected 1 argument."
  echo "Usage $0 VERSION"
  exit 1
fi

version="$1"
if [ -z $package_name ]; then
    . "${cur_dir}/../common/config.sh"
fi

upstream_tarball_name="${package_name}_${version}.tar.gz"

ref="tags/v${version}"
#branch="main" ref="heads/${branch}"
upstream_tarball_url="https://github.com/loicrouchon/jvm-finder/archive/refs/${ref}.tar.gz"

echo "Downloading upstream tarball from ${upstream_tarball_url}"
curl -sL "${upstream_tarball_url}" -o "${upstream_tarball_name}"

echo "Unpacking upstream tarball ${upstream_tarball_name}"
tar xzf "${upstream_tarball_name}"
rm -f "${upstream_tarball_name}"

