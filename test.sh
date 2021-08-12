#!/bin/bash

set -euo pipefail

base="$HOME/tmp/patch-test"
vault_token=devroot
export VAULT_ADDR=http://localhost:8200

function finish {
  kill "$pid"
  rm -fr "$base"
}

trap finish EXIT

mkdir -p "$base"
touch "$base/vault.log"
touch "$base/audit.log"

vault server -dev -dev-root-token-id="$vault_token" "$@" 2> "$base/vault.log" &
pid=$!

while ! nc -w 1 localhost 8200 </dev/null; do
  echo "Waiting for Vault to start"
  sleep 1
done

vault audit enable file file_path="$base/audit.log"

cat << EOF > "$base/allow-policy.hcl"
path "secret*" {
  capabilities = [ "create", "read", "update", "delete", "list", "patch" ]
}
EOF

cat << EOF > "$base/deny-policy.hcl"
path "secret*" {
  capabilities = [ "create", "read", "update", "delete", "list" ]
}
EOF

# shellcheck disable=SC2086
VAULT_TOKEN="$vault_token" vault policy write p1 $base/allow-policy.hcl
# shellcheck disable=SC2086
VAULT_TOKEN="$vault_token" vault policy write p2 $base/deny-policy.hcl

allow_token=$(vault token create -format=json -policy="p1" | jq -r ".auth.client_token")
deny_token=$(vault token create -format=json -policy="p2" | jq -r ".auth.client_token")

VAULT_TOKEN="$vault_token" vault kv put secret/foo bar=baz quux=wibble wobble=wubble
quux=$(VAULT_TOKEN="$vault_token" vault kv get -format=json secret/foo | jq -r '.data.data.quux')
echo

echo "checking quux: it should be wibble"
if [ "$quux" != "wibble" ]; then
  echo "expected quux to be wibble but it was $quux"
  exit 1
else
  echo "test 1: passed"
fi
echo

echo "patching quux to lol using allow_token (should work)"
curl -s -X PATCH -H "X-Vault-Token: $allow_token" -d '{"data":{"quux": "lol"}}' $VAULT_ADDR/v1/secret/data/foo
quux=$(VAULT_TOKEN="$vault_token" vault kv get -format=json secret/foo | jq -r '.data.data.quux')

if [ "$quux" != "lol" ]; then
  echo "expected quux to be lol but it was $quux"
  exit 1
else
  echo "test 2: passed"
fi

echo "patching quux to lawl using deny_token (should fail)"
curl -s -X PATCH -H "X-Vault-Token: $deny_token" -d '{"data":{"quux": "lawl"}}' $VAULT_ADDR/v1/secret/data/foo
quux=$(VAULT_TOKEN="$vault_token" vault kv get -format=json secret/foo | jq -r '.data.data.quux')

if [ "$quux" == "lawl" ]; then
  echo "expected quux to be lol but it was $quux"
  exit 1
elif [ "$quux" == "lol" ]; then
  echo "test 3: passed"
else
  echo "something weird happened. quux == $quux"
  exit 2
fi
