#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -e

fail() {
  echo "$1" 1>&2
  return 1
}

actual_output=$(cat "${VAULT_AGENT_TEMPLATE_DESTINATION}")
if [[ "$actual_output" != "${VAULT_AGENT_EXPECTED_OUTPUT}" ]]; then
  fail "expected '${VAULT_AGENT_EXPECTED_OUTPUT}' to be the Agent output, but got: '$actual_output'"
fi
