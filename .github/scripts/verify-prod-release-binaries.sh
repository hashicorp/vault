#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

set -e

binaries="$1"

if [ -z "$binaries" ]; then
    echo "Error: JSON input is required."
    exit 1
fi

##########################################
# Initialize result holders
##########################################

version_missing=""
variant_missing=""
os_missing=""

valid_versions_output=""
valid_variants_output=""
valid_os_output=""

##########################################
# 1. Verify Versions
##########################################

invalid_versions=$(echo "$binaries" | jq -r '.invalid_versions[]?')

valid_versions_output+="Versions:\n  "

versions=$(echo "$binaries" | jq -r '.valid_versions | keys[]')

for v in $versions; do
    valid_versions_output+="✓ $v, "
done

valid_versions_output+="\n"

if [ -n "$invalid_versions" ]; then
    for v in $invalid_versions; do
        version_missing+="Missing Versions: $v\n"
    done
fi

##########################################
# 2. Verify Variants
##########################################

for v in $versions; do
    base_version="${v%-lts}"

    valid_variants_output+="Version: $v\n  "

    req=(
        "$base_version"
        "${base_version}+ent"
        "${base_version}+ent.fips"
        "${base_version}+ent.hsm"
        "${base_version}+ent.hsm.fips"
    )

    actual_variants=$(echo "$binaries" | jq -r ".valid_versions[\"$v\"].variants[].variant")

    for r in "${req[@]}"; do
        if echo "$actual_variants" | grep -q "^$r"; then
            valid_variants_output+="✓ $r, "
        else
            variant_missing+="Version $v is missing variant: $r\n"
        fi
    done

    valid_variants_output+="\n"
done

##########################################
# 3. Verify OS
##########################################

for v in $versions; do
    base_version="${v%-lts}"
    variants=$(echo "$binaries" | jq -r ".valid_versions[\"$v\"].variants[].variant")

    valid_os_output+="Version: $v"

    for variant in $variants; do
        valid_os_output+="\n  Variant: $variant\n  "

        os_list=$(echo "$binaries" | jq -r \
            ".valid_versions[\"$v\"].variants[] | select(.variant==\"$variant\") | .os[]")

        if [[ "$variant" == "$base_version" || "$variant" == "${base_version}+ent" ]]; then
            for os in darwin freebsd linux netbsd openbsd solaris windows; do
                if grep -qx "$os" <<< "$os_list"; then
                    valid_os_output+="  ✓ $os, "
                else
                    os_missing+="Variant $variant of version $v is missing OS: $os\n"
                fi
            done
        else
            for os in $os_list; do
                if [[ "$os" == "linux" ]]; then
                    valid_os_output+="  ✓ linux"
                else
                    os_missing+="Variant $variant of version $v has invalid OS: $os\n"
                fi
            done
        fi
    done

    valid_os_output+="\n"
done

##########################################
# 4. OUTPUT LOGIC
##########################################

if [ -z "$version_missing" ] && [ -z "$variant_missing" ] && [ -z "$os_missing" ]; then
    echo "*Available Versions:*"
    printf "%b" "$valid_versions_output"
    exit 0
fi

##########################################
# If ANYTHING is missing → print full output
##########################################

echo "*Available Versions:*"
printf "%b" "$valid_versions_output"

echo -e "\n*Available Variants:*"
printf "%b" "$valid_variants_output"

echo -e "\n*Available OS:*"
printf "%b" "$valid_os_output"

if [ -n "$version_missing" ]; then
    echo -e "\n*Missing Release Binary Versions:*"
    printf "%b" "$version_missing"
fi

if [ -n "$variant_missing" ]; then
    echo -e "\n*Missing Release Binary Variants:*"
    printf "%b" "$variant_missing"
fi

if [ -n "$os_missing" ]; then
    echo -e "\n*Missing Release Binary OS:*"
    printf "%b" "$os_missing"
fi
