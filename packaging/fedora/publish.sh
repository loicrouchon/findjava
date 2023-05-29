#!/usr/bin/env sh
set -e
#set -x

cur_dir="$(dirname "$(realpath "$0")")"

cd "${cur_dir}/build"
rm -rf "fedora-copr-symly"
echo "Preparing publication (cloning fedora copr repository)"
git clone https://github.com/loicrouchon/fedora-copr-symly
cd "fedora-copr-symly"
cp ../jvm-finder/jvm-finder.spec jvm-finder.spec
git add jvm-finder.spec
git commit -m "Publish jvm-finder $(cat jvm-finder.spec | grep 'Version:' | tr -d ' ' | tr ':' ' ')"
git push
