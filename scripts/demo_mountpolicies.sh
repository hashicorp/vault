#!/usr/bin/env bash

# Assumes VAULT_ADDR and VAULT_TOKEN are set.

set -ex

policy='
auth "approle" "my*" {
	actions = ["create-role", "update-role", "update-role-secret-id"]
    allow {
        role_name = ["test-role-1"]
    }
}
'
echo "$policy" > /tmp/policy

vault auth disable myapprole1
vault policy write my-approle /tmp/policy
vault policy read my-approle
# Should yield an empty string, since no mounts are in scope
vault policy read -compiled my-approle

vault auth enable -path=myapprole1 approle
vault policy read -compiled my-approle

vault auth enable -path=myapprole2 approle
vault policy read -compiled my-approle

