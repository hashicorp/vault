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

defer_stop_vault() {
    echo "Stopping Vault..."
    kill $VAULT_PID
    sleep 1
}

trap defer_stop_vault INT TERM EXIT

export VAULT_ADDR=http://127.0.0.1:8200

echo "Mounting all builtin plugins..."

# Enable auth plugins
codeLinesStarted=false

while read -r line; do
    if [[ $line == *"credentialBackends:"* ]] ; then
        codeLinesStarted=true
    elif [[ $line == *"databasePlugins:"* ]] ; then
        break
    elif [ $codeLinesStarted = true ] && [[ $line == *"consts.Deprecated"* || $line == *"consts.PendingRemoval"* || $line == *"consts.Removed"* ]] ; then
        auth_plugin_previous=""
    elif [ $codeLinesStarted = true ] && [[ $line =~ ^\s*\"(.*)\"\:.*$ ]] ; then
        auth_plugin_current=${BASH_REMATCH[1]}

        if [[ -n "${auth_plugin_previous}" ]] ; then
            echo "enabling auth plugin: ${auth_plugin_previous}"
            vault auth enable "${auth_plugin_previous}"
        fi

        auth_plugin_previous="${auth_plugin_current}"
    fi
done <../../vault/helper/builtinplugins/registry.go

if [[ -n "${auth_plugin_previous}" ]] ; then
    echo "enabling auth plugin: ${auth_plugin_previous}"
    vault auth enable blah
fi

# Enable secrets plugins
codeLinesStarted=false

while read -r line; do
    if [[ $line == *"logicalBackends:"* ]] ; then
        codeLinesStarted=true
    elif [[ $line == *"addExternalPlugins("* ]] ; then
        break
    elif [ $codeLinesStarted = true ] && [[ $line == *"consts.Deprecated"* || $line == *"consts.PendingRemoval"* || $line == *"consts.Removed"* ]] ; then
        secrets_plugin_previous=""
    elif [ $codeLinesStarted = true ] && [[ $line =~ ^\s*\"(.*)\"\:.*$ ]] ; then
        secrets_plugin_current=${BASH_REMATCH[1]}

        if [[ -n "${secrets_plugin_previous}" ]] ; then
            echo "enabling secrets plugin: ${secrets_plugin_previous}"
            vault secrets enable "${secrets_plugin_previous}"
        fi

        secrets_plugin_previous="${secrets_plugin_current}"
    fi
done <../../vault/helper/builtinplugins/registry.go

if [[ -n "${secrets_plugin_previous}" ]] ; then
    echo "enabling secrets plugin: ${secrets_plugin_previous}"
    vault secrets enable "${secrets_plugin_previous}"
fi

# Enable enterprise features
entRegFile=../../vault/helper/builtinplugins/registry_util_ent.go
if [ -f $entRegFile ] && [[ -n "${VAULT_LICENSE}" ]]; then
    vault write sys/license text="${VAULT_LICENSE}"

    codeLinesStarted=false

    while read -r line; do
        if [[ $line == *"ExternalPluginsEnt:"* ]] ; then
            codeLinesStarted=true
        elif [[ $line == *"addExtPluginsEntImpl("* ]] ; then
            break
        elif [ $codeLinesStarted = true ] && [[ $line == *"consts.Deprecated"* || $line == *"consts.PendingRemoval"* || $line == *"consts.Removed"* ]] ; then
            secrets_plugin_previous=""
        elif [ $codeLinesStarted = true ] && [[ $line =~ ^\s*\"(.*)\"\:.*$ ]] ; then
            ent_plugin_current=${BASH_REMATCH[1]}

            if [[ -n "${ent_plugin_previous}" ]] ; then
                echo "enabling enterprise plugin: ${ent_plugin_previous}"
                vault secrets enable "${ent_plugin_previous}"
            fi

            ent_plugin_previous="${ent_plugin_current}"
        fi
    done <$entRegFile

    if [[ -n "${ent_plugin_previous}" ]] ; then
        echo "enabling enterprise plugin: ${ent_plugin_previous}"
        vault secrets enable "${ent_plugin_previous}"
    fi
fi

# Output OpenAPI, optionally formatted
if [ "$1" == "-p" ]; then
    curl --header 'X-Vault-Token: root' \
         --data '{"generic_mount_paths": true}' \
            'http://127.0.0.1:8200/v1/sys/internal/specs/openapi' | jq > openapi.json
else
    curl --header 'X-Vault-Token: root' \
         --data '{"generic_mount_paths": true}' \
            'http://127.0.0.1:8200/v1/sys/internal/specs/openapi' > openapi.json
fi

echo
echo "openapi.json generated"
echo
