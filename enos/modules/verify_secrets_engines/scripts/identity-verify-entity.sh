#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$ENTITY_ALIAS_ID" ]] && fail "ENTITY_ALIAS_ID env variable has not been set"
[[ -z "$ENTITY_GROUP_IDS" ]] && fail "ENTITY_GROUP_IDS env variable has not been set"
[[ -z "$ENTITY_METADATA" ]] && fail "ENTITY_METADATA env variable has not been set"
[[ -z "$ENTITY_NAME" ]] && fail "ENTITY_NAME env variable has not been set"
[[ -z "$ENTITY_POLICIES" ]] && fail "ENTITY_POLICIES env variable has not been set"
[[ -z "$VAULT_ADDR" ]] && fail "VAULT_ADDR env variable has not been set"
[[ -z "$VAULT_INSTALL_DIR" ]] && fail "VAULT_INSTALL_DIR env variable has not been set"
[[ -z "$VAULT_TOKEN" ]] && fail "VAULT_TOKEN env variable has not been set"


binpath=${VAULT_INSTALL_DIR}/vault
test -x "$binpath" || fail "unable to locate vault binary at $binpath"

export VAULT_FORMAT=json
if ! output=$("$binpath" read "identity/entity/name/$ENTITY_NAME" 2>&1); then
  fail "failed to read identity/entity/name/$ENTITY_NAME: $output"
fi

if ! jq -Mec --arg ALIAS "$ENTITY_ALIAS_ID" '.data.aliases[0].id == $ALIAS' <<< "$output"; then
  fail "entity alias ID does not match, expected: $ENTITY_ALIAS_ID, got: $(jq -Mrc '.data.aliases' <<< "$output")"
fi

if ! jq -Mec --argjson GROUPS "$ENTITY_GROUP_IDS" '.data.group_ids | sort as $have | $GROUPS | sort as $want | $have == $want' <<< "$output"; then
  fail "entity group ID's do not match, expected: $ENTITY_GROUP_IDS, got: $(jq -Mrc '.data.group_ids' <<< "$output")"
fi

if ! jq -Mec --argjson METADATA "$ENTITY_METADATA" '.data.metadata == $METADATA' <<< "$output"; then
  fail "entity metadata does not match, expected: $ENTITY_METADATA, got: $(jq -Mrc '.data.metadata' <<< "$output")"
fi

if ! jq -Mec --argjson POLICIES "$ENTITY_POLICIES" '.data.policies == $POLICIES' <<< "$output"; then
  fail "entity policies do not match, expected: $ENTITY_POLICIES, got: $(jq -Mrc '.data.policies' <<< "$output")"
fi
