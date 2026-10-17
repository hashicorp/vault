#!/bin/sh
# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

# Asserts that ACTUAL_CONTENT contains EXPECTED_CONTENT.
# LABEL names the check in pass/fail messages (e.g. "rendered secret" or "proxy response").

set -eu

fail() {
  echo "$1" 1>&2
  exit 1
}

if ! echo "$ACTUAL_CONTENT" | grep -q "$EXPECTED_CONTENT"; then
  fail "[$LABEL] expected '$EXPECTED_CONTENT', got: '$ACTUAL_CONTENT'"
fi

echo "[$LABEL] verification passed!"
