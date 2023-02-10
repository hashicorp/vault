#!/usr/bin/env bash

# The Vault replication smoke test, documented in
# https://docs.google.com/document/d/16sjIk3hzFDPyY5A9ncxTZV_9gnpYSF1_Vx6UA1iiwgI/edit#heading=h.kgrxf0f1et25

set -e

edition=${vault_edition}

function fail() {
	echo "$1" 1>&2
	exit 1
}

# Replication status endpoint should have data.mode disabled for OSS release
status=$(curl -s http://localhost:8200/v1/sys/replication/status)
if [ "$edition" == "oss" ]; then
  if [ "$(jq -r '.data.mode' <<< "$status")" != "disabled" ]; then
    fail "replication data mode is not disabled for OSS release!"
  fi
else
  if [ "$(jq -r '.data.dr' <<< "$status")" == "" ]; then
    fail "DR replication should be available for an ENT release!"
  fi
  if [ "$(jq -r '.data.performance' <<< "$status")" == "" ]; then
    fail "Performance replication should be available for an ENT release!"
  fi
fi
