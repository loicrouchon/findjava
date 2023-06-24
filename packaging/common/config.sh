#!/usr/bin/env sh
set -e

if [ -z "$version" ]; then
    echo "version is not defined"
    exit 1
fi

package_name="findjava"
package_version_dir="${package_name}-${version}"
