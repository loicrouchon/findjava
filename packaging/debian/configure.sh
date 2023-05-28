#!/usr/bin/env sh
set -x
#set -e
apt update
apt install -y devscripts debhelper golang default-jdk-headless
