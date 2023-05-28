#!/usr/bin/env sh
set -x
#set -e

docker run -ti \
    -v "$(pwd)/../..":/workspace \
    -w /workspace/packaging/debian \
    debian:bullseye \
    bash -c """
./configure.sh
bash
"""
