#!/usr/bin/env sh
set -e
#set -x

cur_dir="$(dirname "$(realpath "$0")")"

echo "Configure environment: PPA URL"
PPA_URL="ppa:loicrouchon/symly"
dput "${PPA_URL}" "${cur_dir}/build/findjvm/"*_source.changes
