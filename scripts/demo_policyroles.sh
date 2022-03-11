#!/usr/bin/env bash

# Assumes VAULT_ADDR and VAULT_TOKEN are set.

set -ex

fail() {
    echo "$@_" 1>&2
    exit 1
}

mount_policy() {
    path=$1
    actions=$2
    policy=$(cat - <<EOF
    auth "approle" "myapprole*" {
        actions = [$actions]
        allow {
            role_name = ["test-role-1"]
        }
    }
EOF
)

    vault write $path policy="$policy"
}

echo '{"mount_actions": {"auth/my.*": "create-role,update-role"}}' |
  vault write sys/policies/acl/role/myrole -

vault read sys/policies/acl/role/myrole

vault auth disable myapprole1
vault auth enable -path=myapprole1 approle

vault policy write allow-role-policy-write <(cat - <<EOF
path "sys/policies/acl/role/myrole*" {
    capabilities = ["update"]
}
path "sys/policies/acl/*" {
    capabilities = ["read"]
}
EOF
)

token=$(vault token create -field=token -policy=allow-role-policy-write)

VAULT_TOKEN=$token
mount_policy sys/policies/acl/mypolicy '"create-role", "update-role"' &&
  fail "shouldn't be allowed to write using regular policy path"

mount_policy sys/policies/acl/role/myrole/mypolicy '"delete-role", "update-role"' &&
  fail "shouldn't be allowed to use delete-role action in a policy"

mount_policy sys/policies/acl/role/myrole/mypolicy '"create-role", "update-role"'

vault policy read mypolicy
vault policy read -compiled mypolicy
