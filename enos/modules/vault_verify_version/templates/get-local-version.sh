#!/usr/bin/env bash

set -e

fail() {
	echo "$1" 1>&2
	exit 1
}

test -x "${vault_local_binary_path}" || fail "unable to locate vault binary at ${vault_local_binary_path}"

${vault_local_binary_path} version
