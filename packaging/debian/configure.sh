#!/usr/bin/env sh
set -e
#set -x
apt update
apt install -y devscripts debhelper golang default-jdk-headless
