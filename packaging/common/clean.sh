#!/usr/bin/env sh
set -ex

if [ "$#" -ne 2 ]; then
  echo "Error: Invalid number of arguments. Expected 2 arguments."
  echo "Usage $0 BUILD_DIR"
  exit 1
fi

rm -rf "$1"
mkdir -p "$1/$2"
