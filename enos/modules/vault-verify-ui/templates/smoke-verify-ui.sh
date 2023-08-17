#!/usr/bin/env bash

set -e

fail() {
	echo "$1" 1>&2
	exit 1
}
if [ "$(curl -s -o /dev/null -w "%%{redirect_url}" http://localhost:8200/)" != "http://localhost:8200/ui/" ]; then
    fail "Port 8200 not redirecting to UI"
fi
if curl -s http://localhost:8200/ui/ | grep -q 'Vault UI is not available'; then
    fail "Vault UI is not available"
fi
