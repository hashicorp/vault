#!/usr/bin/env bash
set -euo pipefail

# require jq
command -v jq >/dev/null 2>&1 || { echo "jq is required (brew install jq)"; exit 1; }

# locate vault binary relative to this script
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
if [[ -x "${SCRIPT_DIR}/bin/vault" ]]; then
  VAULT_BIN="${SCRIPT_DIR}/bin/vault"
elif [[ -x "${SCRIPT_DIR}/../bin/vault" ]]; then
  VAULT_BIN="${SCRIPT_DIR}/../bin/vault"
else
  echo "vault binary not found near script. Build it first with: make dev" >&2
  exit 1
fi

# defaults (can be overridden via env)
export VAULT_ADDR="${VAULT_ADDR:-http://127.0.0.1:8200}"
export VAULT_TOKEN="${VAULT_TOKEN:-root}"

# portable base64 decode (macOS uses -D, GNU uses -d)
b64dec() {
  if echo "SGVsbG8=" | base64 -D >/dev/null 2>&1; then base64 -D; else base64 -d; fi
}

"$VAULT_BIN" status

"$VAULT_BIN" kv put -format=json secret/hello foo=bar >/dev/null
"$VAULT_BIN" kv get secret/hello

# enable transit if not already enabled
if ! "$VAULT_BIN" secrets list -format=json | jq -e 'has("transit/")' >/dev/null; then
  "$VAULT_BIN" secrets enable transit
fi

"$VAULT_BIN" write -f -format=json transit/keys/test >/dev/null

PLAINTEXT_B64="$(printf 'hello' | base64 | tr -d '\n')"
CT="$("$VAULT_BIN" write -format=json transit/encrypt/test plaintext="$PLAINTEXT_B64" | jq -r '.data.ciphertext')"

PT_B64="$("$VAULT_BIN" write -format=json transit/decrypt/test ciphertext="$CT" | jq -r '.data.plaintext')"
DECRYPTED="$(printf '%s' "$PT_B64" | b64dec)"

echo "Decrypted: $DECRYPTED"
[[ "$DECRYPTED" == "hello" ]] && echo "Round-trip ✔" || { echo "Round-trip ❌"; exit 1; }
