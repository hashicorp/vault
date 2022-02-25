#!/bin/sh

set -e

# Generate an OpenAPI document for all backends.
#
# Assumptions:
#
#   1. Vault has been checked out at an appropriate version and built
#   2. vault executable is in your path
#   3. Vault isn't already running

go run genplugins.go

echo "Starting Vault..."
if pgrep -x "vault" > /dev/null
then
    echo "Vault is already running. Aborting."
    exit 1
fi

vault server -dev -dev-root-token-id=root &
sleep 2
VAULT_PID=$!

echo "Mounting all builtin backends..."

# Credential backends
for row in $(cat plugins.json | jq -r '.CredentialBackends[] | @base64'); do
    _jq() {
     echo ${row} | base64 --decode | jq -r ${1}
    }
    val=$(_jq '.name')
    vault auth enable ${val}
done

# Secrets backends 
for row in $(cat plugins.json | jq -r '.LogicalBackends[] | @base64'); do
    _jq() {
     echo ${row} | base64 --decode | jq -r ${1}
    }
    val=$(_jq '.name')
    vault secrets enable ${val}
done


# Output OpenAPI, optionally formatted
if [ "$1" == "-p" ]; then
  curl -H "X-Vault-Token: root" "http://127.0.0.1:8200/v1/sys/internal/specs/openapi" | jq > openapi.json
else
  curl -H "X-Vault-Token: root" "http://127.0.0.1:8200/v1/sys/internal/specs/openapi" > openapi.json
fi

kill $VAULT_PID
sleep 1

echo "\nopenapi.json generated."