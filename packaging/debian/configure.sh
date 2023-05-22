#!/usr/bin/env sh
set -ex
# docker run -ti -v (pwd):/workspace -w /workspace debian:bullseye
apt update
apt install -y devscripts debhelper golang default-jdk-headless
