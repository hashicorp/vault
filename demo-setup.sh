#!/bin/bash

# Demo setup script for IBM Cloud auth
#
# Usage:
#   ./demo-setup.sh <api_key> <account_id> [service_id] [access_group_id]
#
# Arguments:
#   api_key          Admin API key Vault uses to verify identities and access groups.
#   account_id       IBM Cloud account ID that every login must belong to.
#   service_id       (optional) Service ID that owns api_key. Only needed for
#                    `vault write -f auth/ibmcloud/config/rotate-root`.
#                    Pass "" to skip it while still providing access_group_id.
#   access_group_id  (optional) Access group the roles are bound to. If omitted,
#                    bound_access_group_ids is not set on the roles, so login only
#                    requires the identity to belong to the account.
#
# Environment:
#   IAM_ENDPOINT        IAM endpoint (default: https://iam.test.cloud.ibm.com)
#   VAULT_ADDR/VAULT_TOKEN  Must point at a running, unsealed Vault.
set -e

if [ $# -lt 2 ] || [ $# -gt 4 ]; then
    echo "Usage: $0 <api_key> <account_id> [service_id] [access_group_id]" >&2
    exit 1
fi

API_KEY="$1"
ACCOUNT_ID="$2"
SERVICE_ID="${3:-}"
ACCESS_GROUP_ID="${4:-}"
IAM_ENDPOINT="${IAM_ENDPOINT:-https://iam.test.cloud.ibm.com}"

echo "Enabling IBM Cloud auth method..."
./bin/vault auth enable ibmcloud || echo "Auth method already enabled"

echo "Configuring IBM Cloud auth backend..."
CONFIG_ARGS=(
    "iam_endpoint=${IAM_ENDPOINT}"
    "iam_client_timeout=30"
    "api_key=${API_KEY}"
    "account_id=${ACCOUNT_ID}"
)
if [ -n "${SERVICE_ID}" ]; then
    CONFIG_ARGS+=("service_id=${SERVICE_ID}")
fi
./bin/vault write auth/ibmcloud/config "${CONFIG_ARGS[@]}"

echo "Creating policies..."
./bin/vault policy write my-policy my-policy.hcl
./bin/vault policy write my-writer-policy my-writer-policy.hcl

ROLE_EXTRA_ARGS=()
if [ -n "${ACCESS_GROUP_ID}" ]; then
    ROLE_EXTRA_ARGS+=("bound_access_group_ids=${ACCESS_GROUP_ID}")
else
    echo "No access group provided: roles will only require account membership."
fi

echo "Creating role: my-service-role..."
./bin/vault write auth/ibmcloud/role/my-service-role \
    type=apikey \
    bound_account_id="${ACCOUNT_ID}" \
    token_policies=my-policy \
    token_ttl=1h \
    token_max_ttl=24h \
    ${ROLE_EXTRA_ARGS[@]+"${ROLE_EXTRA_ARGS[@]}"}

echo "Creating role: my-service-writer-role..."
./bin/vault write auth/ibmcloud/role/my-service-writer-role \
    type=apikey \
    bound_account_id="${ACCOUNT_ID}" \
    token_policies=my-writer-policy \
    token_ttl=1h \
    token_max_ttl=24h \
    ${ROLE_EXTRA_ARGS[@]+"${ROLE_EXTRA_ARGS[@]}"}

echo "Setup complete!"
echo ""
echo "Test authentication with:"
echo "  ./bin/vault login -method=ibmcloud role=my-service-writer-role type=apikey apikey=YOUR_API_KEY iam_endpoint=${IAM_ENDPOINT}"

# Made with Bob
