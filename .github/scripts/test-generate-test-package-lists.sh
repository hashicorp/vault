#!/bin/bash
set -e${DEBUG+x}o pipefail

dir=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

source generate-test-package-lists.sh

get_module_packages() {
    local package_list=($(go list -test -json ./... | jq -r '.ForTest | select(.!=null)' | grep -v vault/integ | grep '^github.com/hashicorp/'))
    
    for package in "${package_list[@]}" ; do
        # Check if the current package already exists in all_packages
        if ! grep "\b$package\b" <<< "${all_packages[@]}" &> /dev/null; then
            all_packages+=($package)
        fi
    done
}

find_packages() {
    for package in "${all_packages[@]}" ; do
        if ! grep "\b${package}\b" <<< "${test_packages[@]}" &> /dev/null ; then
            echo "Error: package ${package} is not present in test_packages"
            exit 1
        fi
    done
}

count_test_packages() {
    count=0
    for test_package in "${test_packages[@]}" ; do
        count=$((${count}+$(wc -w <<< "${test_package}")))
    done

    echo $count
}

all_packages=()

cd "$dir/../.."
get_module_packages

cd "$dir/../../sdk"
get_module_packages

cd "$dir/../../api"
get_module_packages

find_packages

test_package_count=$(count_test_packages)
if (( ${#all_packages[@]} != $test_package_count )) ; then
    echo "Error: there are currently ${#all_packages[@]} packages in the repository but $test_package_count packages in test_packages"

    unused_packages="${test_packages[@]} "
    for ap in ${all_packages[@]} ; do
        unused_packages="$(echo "$unused_packages" | sed -r "s~$ap ~ ~" )"
    done

    echo "Packages in test_packages that aren't used: ${unused_packages// /}"
fi
