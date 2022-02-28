#!/bin/bash

set -e

# Generate an OpenAPI document for all backends.
#
# Assumptions:
#
#   1. Vault has been checked out at an appropriate version and built
#   2. vault executable is in your path
#   3. Vault isn't already running

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

# Read auth backends
flag=false
re='".*"'
while read p; do
    if [[ $p == *"credentialBackends:"* ]] ; then
        flag=true
    elif [ $flag = true ] && [[ $p = *"}"* ]]  ; then
        break
    elif [ $flag = true ] && [[ $p =~ $re ]]; then
        backend=${BASH_REMATCH[0]}
        var1=$(sed -e 's/^"//' -e 's/"$//' <<<"$backend") 
        vault auth enable ${var1}
    fi
done <../../vault/helper/builtinplugins/registry.go

# Read secrets backends
flag=false
re='".*"'
while read p; do
    if [[ $p == *"logicalBackends:"* ]] ; then
        flag=true
    elif [ $flag = true ] && [[ $p = *"}"* ]]  ; then
        break
    elif [ $flag = true ] && [[ $p =~ $re ]]; then
        backend=${BASH_REMATCH[0]}
        var1=$(sed -e 's/^"//' -e 's/"$//' <<<"$backend") 
        vault secrets enable ${var1}
    fi
done <../../vault/helper/builtinplugins/registry.go


# Enable enterprise features
if [[ ! -z "$VAULT_LICENSE" ]]
then
  vault write sys/license text="$VAULT_LICENSE"
  vault secrets enable kmip
  vault secrets enable transform
fi

# Output OpenAPI, optionally formatted
if [ "$1" == "-p" ]; then
  curl -H "X-Vault-Token: root" "http://127.0.0.1:8200/v1/sys/internal/specs/openapi" | jq > openapi.json
else
  curl -H "X-Vault-Token: root" "http://127.0.0.1:8200/v1/sys/internal/specs/openapi" > openapi.json
fi

kill $VAULT_PID
sleep 1

echo "\nopenapi.json generated."