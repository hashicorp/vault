#!/usr/bin/env bash

set -e

fail() {
	echo "$1" 1>&2
	exit 1
}

if [ "${REDIRECT_URL}" != "http://localhost:8200/ui/" ]; then
    fail "Port 8200 not redirecting to UI"
fi
if [ "${UI_URL_RESULT}" != "200" ]; then
    fail "Vault UI is not available"
fi
