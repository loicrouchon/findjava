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

build() {
    echo "Building fedora source package ${package_version_dir}"
    sed "s/\${version}/${version}/" "${cur_dir}/findjvm.spec" > "build/${package_name}/findjvm.spec"
}

"${cur_dir}/../common/clean.sh" "build" "${package_name}"
build
