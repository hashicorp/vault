#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

fail() {
  echo "$1" 1>&2
  exit 1
}

verify_radar_scan_output_file() {
  # Given a file with a radar scan output, filter out tagged false positives and verify that no
  # other secrets remain.
  if ! jq -eMcn '[inputs] | [.[] | select((.tags == null) or (.tags | contains(["ignore_rule"]) | not ))] | length == 0' < "$2"; then
    found=$(jq -eMn '[inputs] | [.[] | select((.tags == null) or (.tags | contains(["ignore_rule"]) | not ))]' < "$2")
    fail "failed to radar secrets output: vault radar detected secrets in $1!: $found"
  fi
}

set -e

[[ -z "$AUDIT_LOG_FILE_PATH" ]] && fail "AUDIT_LOG_FILE_PATH env variable has not been set"
[[ -z "$VAULT_RADAR_INSTALL_DIR" ]] && fail "VAULT_RADAR_INSTALL_DIR env variable has not been set"
# Radar implicitly requires the following for creating the index and running radar itself
[[ -z "$VAULT_RADAR_LICENSE" ]] && fail "VAULT_RADAR_LICENSE env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"
[[ -z "$VAULT_UNIT_NAME" ]] && fail "VAULT_UNIT_NAME env variable has not been set"

radar_bin_path=${VAULT_RADAR_INSTALL_DIR}/vault-radar
test -x "$radar_bin_path" || fail "failed to scan vault audit log: unable to locate radar binary at $radar_bin_path"

# Make sure our audit log file exists.
if [ ! -f "$AUDIT_LOG_FILE_PATH" ]; then
  fail "failed to scan vault audit log: no audit logifile found at $AUDIT_LOG_FILE_PATH"
fi

# Create a readable copy of the audit log.
if ! sudo cp "$AUDIT_LOG_FILE_PATH" audit.log; then
  fail "failed to scan vault audit log: could not copy audit log for scanning"
fi

if ! sudo chmod +r audit.log; then
  fail "failed to scan vault audit log: could not make audit log copy readable"
fi

# Create a radar index file of our KVv2 secret values.
if ! out=$($radar_bin_path index vault --offline --disable-ui --outfile index.jsonl 2>&1); then
  fail "failed to generate vault-radar index of vault cluster: $out"
fi

# Write our ignore rules to avoid known false positives.
mkdir -p "$HOME/.hashicorp/vault-radar"
cat >> "$HOME/.hashicorp/vault-radar/ignore.yaml" << EOF
- secret_values:
  - "hmac-sha256:*"
EOF

# Scan the audit log for known secrets via the audit log and other secrets using radars built-in
# secret types.
if ! out=$("$radar_bin_path" scan file --offline --disable-ui -p audit.log --index-file index.jsonl -f json -o audit-secrets.json 2>&1); then
  fail "failed to scan vault audit log: vault-radar scan file failed: $out"
fi

verify_radar_scan_output_file vault-audit-log audit-secrets.json

# Scan the vault journal for known secrets via the audit log and other secrets using radars built-in
# secret types.
if ! out=$(sudo journalctl --no-pager -u "$VAULT_UNIT_NAME" -a | "$radar_bin_path" scan file --offline --disable-ui --index-file index.jsonl -f json -o journal-secrets.json 2>&1); then
  fail "failed to scan vault journal: vault-radar scan file failed: $out"
fi

verify_radar_scan_output_file vault-journal journal-secrets.json
