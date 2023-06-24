#!/usr/bin/env sh
set -e
#set -x

cur_dir="$(dirname "$(realpath "$0")")"

cd "${cur_dir}/build"
rm -rf "fedora-copr-symly"
echo "Preparing publication (cloning fedora copr repository)"
git clone git@github.com:loicrouchon/fedora-copr-symly.git
cd "fedora-copr-symly"
cp ../findjvm/findjvm.spec findjvm.spec
git add findjvm.spec
git commit -m "Publish findjvm $(cat findjvm.spec | grep 'Version:' | tr -d ' ' | tr ':' ' ')"
git push
