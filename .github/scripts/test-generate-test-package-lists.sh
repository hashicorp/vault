#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

set -e${DEBUG+x}o pipefail

#
# This script is run to make sure that every package returned by 
# go list -test ./... (when run from the repo root, api/, and sdk/ directory)
# appear in the test_packages array defined in the sibling file
# generate-test-package-lists.sh
#
# This script is executed as part of the ci workflow triggered by pull_requests
# events. In the event that the job that runs this script fails, examine the
# output of the 'test' step in that job to obtain the list of test packages that
# are missing in the test_packages array or that should be removed from it.
#

dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source generate-test-package-lists.sh

get_module_packages() {
    for d in "../.." "../../sdk" "../../api"; do
        cd "$dir/$d"
        go list -test -json ./... | jq -r '.ForTest | select(.!=null)' |
          grep -v vault/integ | grep '^github.com/hashicorp/'
    done |sort -u
}

count_test_packages() {
    count=0
    for test_package in "${test_packages[@]}" ; do
        count=$((${count}+$(wc -w <<< "${test_package}")))
    done

    echo $count
}

contains() {
    target=$1; shift
    for i; do
        if [[ "$i" == "$target" ]]; then
            return 0
        fi
    done
    return 1
}

all_packages=( $(get_module_packages) )

for package in "${all_packages[@]}" ; do
    if ! contains $package ${test_packages[@]}; then
        echo "Error: package ${package} is not present in test_packages"
        exit 1
    fi
done

test_package_count=$(count_test_packages)
if (( ${#all_packages[@]} != $test_package_count )) ; then
    echo "Error: there are currently ${#all_packages[@]} packages in the repository but $test_package_count packages in test_packages"

    unused_packages="${test_packages[@]} "
    for ap in ${all_packages[@]} ; do
        unused_packages="$(echo "$unused_packages" | sed -r "s~$ap ~ ~" )"
    done

    echo "Packages in test_packages that aren't used: ${unused_packages// /}"
    exit 1
fi
