#!/bin/bash

set -e

# Generate an OpenAPI document for all backends.
#
# Assumptions:
#
#   1. Vault has been checked out at an appropriate version and built
#   2. vault executable is in your path
#   3. Vault isn't already running
#   4. jq is installed

cd "$(dirname "${BASH_SOURCE[0]}")"

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
codeLinesStarted=false
inQuotesRegex='".*"'
while read -r line; do
    if [[ $line == *"credentialBackends:"* ]] ; then
        codeLinesStarted=true
    elif [ $codeLinesStarted = true ] && [[ $line = *"}"* ]]  ; then
        break
    elif [ $codeLinesStarted = true ] && [[ $line =~ $inQuotesRegex ]] && [[ $line != *"Deprecated"* ]] ; then
        backend=${BASH_REMATCH[0]}
        plugin=$(sed -e 's/^"//' -e 's/"$//' <<<"$backend")
        vault auth enable "${plugin}"
    fi
done <../../vault/helper/builtinplugins/registry.go

# Read secrets backends
codeLinesStarted=false
while read -r line; do
    if [[ $line == *"logicalBackends:"* ]] ; then
        codeLinesStarted=true
    elif [ $codeLinesStarted = true ] && [[ $line = *"}"* ]]  ; then
        break
    elif [ $codeLinesStarted = true ] && [[ $line =~ $inQuotesRegex ]] && [[ $line != *"Deprecated"* ]] ; then
        backend=${BASH_REMATCH[0]}
        plugin=$(sed -e 's/^"//' -e 's/"$//' <<<"$backend")
        vault secrets enable "${plugin}"
    fi
done <../../vault/helper/builtinplugins/registry.go


# Enable enterprise features
entRegFile=../../vault/helper/builtinplugins/registry_util_ent.go
if [ -f $entRegFile ] && [[ -n "$VAULT_LICENSE" ]]; then
  vault write sys/license text="$VAULT_LICENSE"

  inQuotesRegex='".*"'
  codeLinesStarted=false
  while read -r line; do
        if [[ $line == *"ExternalPluginsEnt"* ]] ; then
        codeLinesStarted=true
    elif [ $codeLinesStarted = true ] && [[ $line = *"}"* ]]  ; then
        break
    elif [ $codeLinesStarted = true ] && [[ $line =~ $inQuotesRegex ]] && [[ $line != *"Deprecated"* ]] ; then
        backend=${BASH_REMATCH[0]}
        plugin=$(sed -e 's/^"//' -e 's/"$//' <<<"$backend")
        vault secrets enable "${plugin}"
    fi
  done <$entRegFile
fi

# Output OpenAPI, optionally formatted
if [ "$1" == "-p" ]; then
  curl -H "X-Vault-Token: root" "http://127.0.0.1:8200/v1/sys/internal/specs/openapi" | jq > openapi.json
else
  curl -H "X-Vault-Token: root" "http://127.0.0.1:8200/v1/sys/internal/specs/openapi" > openapi.json
fi

kill $VAULT_PID
sleep 1

echo
echo "openapi.json generated"
