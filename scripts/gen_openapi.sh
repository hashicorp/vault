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

export VAULT_ADDR=http://127.0.0.1:8200

echo "Mounting all builtin plugins..."

# Enable auth plugins
vault auth enable "alicloud"
vault auth enable "approle"
vault auth enable "aws"
vault auth enable "azure"
vault auth enable "centrify"
vault auth enable "cert"
vault auth enable "cf"
vault auth enable "gcp"
vault auth enable "github"
vault auth enable "jwt"
vault auth enable "kerberos"
vault auth enable "kubernetes"
vault auth enable "ldap"
vault auth enable "oci"
vault auth enable "oidc"
vault auth enable "okta"
vault auth enable "pcf"
vault auth enable "radius"
vault auth enable "userpass"

# Enable secrets plugins
vault secrets enable "ad"
vault secrets enable "alicloud"
vault secrets enable "aws"
vault secrets enable "azure"
vault secrets enable "consul"
vault secrets enable "database"
vault secrets enable "gcp"
vault secrets enable "gcpkms"
vault secrets enable "kubernetes"
vault secrets enable "kv"
vault secrets enable "ldap"
vault secrets enable "mongodbatlas"
vault secrets enable "nomad"
vault secrets enable "openldap"
vault secrets enable "pki"
vault secrets enable "rabbitmq"
vault secrets enable "ssh"
vault secrets enable "terraform"
vault secrets enable "totp"
vault secrets enable "transit"

# Enable enterprise features
if [[ -n "${VAULT_LICENSE}" ]]; then
    vault write sys/license text="${VAULT_LICENSE}"

    vault secrets enable "keymgmt"
    vault secrets enable "kmip"
    vault secrets enable "transform"
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

kill $VAULT_PID
sleep 1

echo
echo "openapi.json generated"
