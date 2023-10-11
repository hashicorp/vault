#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


# The Vault replication smoke test, documented in
# https://docs.google.com/document/d/16sjIk3hzFDPyY5A9ncxTZV_9gnpYSF1_Vx6UA1iiwgI/edit#heading=h.kgrxf0f1et25

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

# Replication STATUS endpoint should have data.mode disabled for CE release
if [ "$VAULT_EDITION" == "ce" ]; then
  if [ "$(echo "${STATUS}" | jq -r '.data.mode')" != "disabled" ]; then
    fail "replication data mode is not disabled for CE release!"
  fi
else
  if [ "$(echo "${STATUS}" | jq -r '.data.dr')" == "" ]; then
    fail "DR replication should be available for an ENT release!"
  fi
  if [ "$(echo "${STATUS}" | jq -r '.data.performance')" == "" ]; then
    fail "Performance replication should be available for an ENT release!"
  fi
fi
