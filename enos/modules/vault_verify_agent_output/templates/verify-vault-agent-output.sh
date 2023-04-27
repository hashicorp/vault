#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -e

fail() {
  echo "$1" 1>&2
  return 1
}

actual_output=$(cat ${vault_agent_template_destination})
if [[ "$actual_output" != "${vault_agent_expected_output}" ]]; then
  fail "expected '${vault_agent_expected_output}' to be the Agent output, but got: '$actual_output'"
fi
