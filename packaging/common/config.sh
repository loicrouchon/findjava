#!/usr/bin/env sh

if [ -z "$version" ]; then
    echo "version is not defined"
    exit 1
fi

package_name="jvm-finder"
package_version_dir="${package_name}-${version}"
