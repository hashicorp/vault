#!/usr/bin/env bash
# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

fail() {
  echo "$1" 1>&2
  exit 1
}

if ! out=$(sudo systemctl stop consul 2>&1); then
  fail "failed to stop consul: $out: $(sudo systemctl status consul)"
fi

if ! out=$(sudo systemctl daemon-reload 2>&1); then
  fail "failed to daemon-reload systemd: $out" 1>&2
fi

if ! out=$(sudo systemctl start consul 2>&1); then
  fail "failed to start consul: $out: $(sudo systemctl status consul)"
fi
