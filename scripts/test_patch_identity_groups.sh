#!/bin/bash

set -ex

#function finish {
#  kill "$pid"
#}
#
#trap finish EXIT

rm -rf ~/devvault
mkdir -p ~/devvault

pkill vault & sleep 1

vault server -dev -log-level=trace -dev-listen-address=0.0.0.0:8200 -dev-root-token-id=devroot 2> ~/devvault/vault.log &
pid=$!

while ! nc -w 1 localhost 8200 </dev/null; do
  echo "Waiting for Vault to start"
  sleep 1
done

export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=devroot

curl --header "X-Vault-Token: $VAULT_TOKEN" --request DELETE "$VAULT_ADDR/v1/identity/group/name/foo"

FOO_ENTITY_ID=$(curl -s --header "X-Vault-Token: $VAULT_TOKEN" --request POST --data '{"name": "foo"}' "$VAULT_ADDR/v1/identity/entity" | jq .data.id)
curl --header "X-Vault-Token: $VAULT_TOKEN" --request POST --data "{\"member_entity_ids\": [$FOO_ENTITY_ID]}" "$VAULT_ADDR/v1/identity/group/name/foo"
curl -s --header "X-Vault-Token: $VAULT_TOKEN" --request GET "$VAULT_ADDR/v1/identity/group/name/foo" | jq .data

BAR_ENTITY_ID=$(curl -s --header "X-Vault-Token: $VAULT_TOKEN" --request POST --data '{"name": "bar"}' "$VAULT_ADDR/v1/identity/entity" | jq .data.id)
curl --header "Content-Type: application/merge-patch+json" --header "X-Vault-Token: $VAULT_TOKEN" --request PATCH --data "{\"member_entity_ids\": [$BAR_ENTITY_ID]}" "$VAULT_ADDR/v1/identity/group/name/foo"
curl -s --header "X-Vault-Token: $VAULT_TOKEN" --request GET "$VAULT_ADDR/v1/identity/group/name/foo" | jq .data

curl --header "Content-Type: application/json-patch+json" --header "X-Vault-Token: $VAULT_TOKEN" --request PATCH --data "{\"patch_json\": [{\"op\": \"add\", \"path\": \"/member_entity_ids/-\", \"value\": $FOO_ENTITY_ID}]}" "$VAULT_ADDR/v1/identity/group/name/foo"
curl -s --header "X-Vault-Token: $VAULT_TOKEN" --request GET "$VAULT_ADDR/v1/identity/group/name/foo" | jq .data
