#!/usr/bin/env sh
set -x
#set -e

cur_dir="$(dirname "$(realpath "$0")")"

if [ "$#" -ne 1 ]; then
  echo "Error: Invalid number of arguments. Expected 1 argument."
  echo "Usage $0 VERSION"
  exit 1
fi

version="$1"
cd "${cur_dir}/build/jvm-finder/jvm-finder-${version}" || exit 1

dpkg-buildpackage -b -us -uc
