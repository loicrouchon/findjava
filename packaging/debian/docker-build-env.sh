#!/usr/bin/env sh
set -x
#set -e

#image="debian:bullseye"
image="ubuntu:latest"
docker run -ti \
    -v "$(pwd)/../..":/workspace \
    -w /workspace/packaging/debian \
    "$image" \
    bash -c """
./configure.sh

echo 'Configure environment: maintainer name and email'
export DEBEMAIL='loic@loicrouchon.com'
export DEBFULLNAME='Loic Rouchon'
echo 'Configure environment: GPG fingerprint'
export GPG_KEY_FINGERPRINT='C3BB9448B16C971103E876BF3A091A0DF2799262'

cat gpg/*.key  | gpg --import

bash
"""


