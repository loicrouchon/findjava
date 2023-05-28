#!/usr/bin/env sh
set -x
#set -e

cur_dir="$(dirname "$(realpath "$0")")"

echo "Configure environment: PPA URL"
PPA_URL="ppa:loicrouchon/symly"
dput "${PPA_URL}" "${cur_dir}/build/jvm-finder/"*.changes
