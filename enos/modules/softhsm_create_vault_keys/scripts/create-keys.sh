#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

set -e

fail() {
  echo "$1" 1>&2
  exit 1
}

[[ -z "$AES_LABEL" ]] && fail "AES_LABEL env variable has not been set"
[[ -z "$HMAC_LABEL" ]] && fail "HMAC_LABEL env variable has not been set"
[[ -z "$PIN" ]] && fail "PIN env variable has not been set"
[[ -z "$SO_PIN" ]] && fail "SO_PIN env variable has not been set"
[[ -z "$TOKEN_LABEL" ]] && fail "TOKEN_LABEL env variable has not been set"
[[ -z "$TOKEN_DIR" ]] && fail "TOKEN_DIR env variable has not been set"

if ! type softhsm2-util &> /dev/null; then
  fail "unable to locate softhsm2-util in PATH. Have you installed softhsm?"
fi

if ! type pkcs11-tool &> /dev/null; then
  fail "unable to locate pkcs11-tool in PATH. Have you installed opensc?"
fi

# Create an HSM slot and return the slot number in decimal value.
create_slot() {
  sudo softhsm2-util --init-token --free --so-pin="$SO_PIN" --pin="$PIN" --label="$TOKEN_LABEL" | grep -oE '[0-9]+$'
}

# Find the location of our softhsm shared object.
find_softhsm_so() {
  sudo find /usr -type f -name libsofthsm2.so -print -quit
}

# Create key a key in the slot. Args: module, key label, id number, key type
keygen() {
  sudo pkcs11-tool --keygen --usage-sign --private --sensitive --usage-wrap \
    --module "$1" \
    -p "$PIN" \
    --token-label "$TOKEN_LABEL" \
    --label "$2" \
    --id "$3" \
    --key-type "$4"
}

# Create our softhsm slot and keys
main() {
  local slot
  if ! slot=$(create_slot); then
    fail "failed to create softhsm token slot"
  fi

  local so
  if ! so=$(find_softhsm_so); then
    fail "unable to locate libsofthsm2.so shared object"
  fi

  if ! keygen "$so" "$AES_LABEL" 1 'AES:32' 1>&2; then
    fail "failed to create AES key"
  fi

  if ! keygen "$so" "$HMAC_LABEL" 2 'GENERIC:32' 1>&2; then
    fail "failed to create HMAC key"
  fi

  # Return our seal configuration attributes as JSON
  cat << EOF
{
  "lib": "${so}",
  "slot": "${slot}",
  "pin": "${PIN}",
  "key_label": "${AES_LABEL}",
  "hmac_key_label": "${HMAC_LABEL}",
  "generate_key": "false"
}
EOF
  exit 0
}

main
