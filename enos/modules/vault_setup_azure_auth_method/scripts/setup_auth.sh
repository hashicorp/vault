#!/usr/bin/env sh

set -eux

vault auth enable azure

vault write auth/azure/config \
    tenant_id="${TENANT_ID}" \
    resource=https://management.azure.com/ \
    client_id="${CLIENT_ID}" \
    client_secret="${CLIENT_SECRET}"

vault write auth/azure/role/dev-role \
    policies="root,default" \
    bound_subscription_ids="${BOUND_SUBSCRIPTION_ID}" \
    bound_resource_groups="${BOUND_RESOURCE_GROUP}"
