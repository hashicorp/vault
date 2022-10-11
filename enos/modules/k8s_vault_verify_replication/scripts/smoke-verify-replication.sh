#!/usr/bin/env sh

# The Vault replication smoke test, documented in
# https://docs.google.com/document/d/16sjIk3hzFDPyY5A9ncxTZV_9gnpYSF1_Vx6UA1iiwgI/edit#heading=h.kgrxf0f1et25

set -e

function fail() {
	echo "$1" 1>&2
	exit 1
}

# Replication STATUS endpoint should have data.mode disabled for OSS release
if [ "$VAULT_EDITION" == "oss" ]; then
  if [ "$(jq -r '.data.mode' <<< "$STATUS")" != "disabled" ]; then
    fail "replication data mode is not disabled for OSS release!"
  fi
else
  if [ "$(jq -r '.data.dr' <<< "$STATUS")" == "" ]; then
    fail "DR replication should be available for an ENT release!"
  fi
  if [ "$(jq -r '.data.performance' <<< "$STATUS")" == "" ]; then
    fail "Performance replication should be available for an ENT release!"
  fi
fi
