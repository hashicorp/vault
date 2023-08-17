#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


# The Vault smoke test to verify the Vault version installed

set -e

fail() {
	echo "$1" 1>&2
	exit 1
}

if [[ "${CHECK_BUILD_DATE}" == "false" ]]; then
  expected_build_date=""
else
  build_date="${BUILD_DATE}"
  if [[ "${build_date}" == "" ]]; then
    build_date=$(echo "${VAULT_STATUS}" | jq -Mr .build_date)
  fi
  expected_build_date=", built $build_date"
fi

vault_expected_version="Vault v${EXPECTED_VERSION} (${VAULT_REVISION})"

case "${VAULT_EDITION}" in
  oss) version_expected="${vault_expected_version}${expected_build_date}";;
	ent) version_expected="${vault_expected_version}${expected_build_date}";;
	ent.hsm) version_expected="${vault_expected_version}${expected_build_date} (cgo)";;
	ent.fips1402) version_expected="${vault_expected_version}${expected_build_date} (cgo)" ;;
	ent.hsm.fips1402) version_expected="${vault_expected_version}${expected_build_date} (cgo)" ;;
  *) fail "(${VAULT_EDITION}) does not match any known Vault editions"
esac

version_expected_nosha=$(echo "$version_expected" | awk '!($3="")' | sed 's/  / /' | sed -e 's/[[:space:]]*$//')

if [[ "${ACTUAL_VERSION}" == "$version_expected_nosha" ]] || [[ "${ACTUAL_VERSION}" == "$version_expected" ]]; then
	echo "Version verification succeeded!"
else
  echo "CHECK_BUILD_DATE: ${CHECK_BUILD_DATE}"
  echo "BUILD_DATE: ${BUILD_DATE}"
  echo "build_date: ${build_date}"
	fail "expected Version=$version_expected or $version_expected_nosha, got: ${ACTUAL_VERSION}"
fi
