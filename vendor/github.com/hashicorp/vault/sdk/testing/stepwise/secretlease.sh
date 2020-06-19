#!/bin/bash

export VAULT_ADDR=http://127.0.0.1:8200
export AGENT_DIR=/tmp/agent
mkdir -p ${AGENT_DIR?}

function generate_templates() {
	echo "Test"
}

function trap_sigterm() {
    echo "Clean shutdown of Vault.."
}

trap 'trap_sigterm' SIGINT SIGTERM

vault server -dev &
sleep 5

cat <<EOF > /tmp/policy.hcl
path "secret/*"         { capabilities = ["read"] }
EOF

vault policy write db /tmp/policy.hcl

echo "Setting up approle.."
vault auth enable approle
vault write auth/approle/role/my-role \
    secret_id_ttl=24h \
    token_num_uses=10 \
    token_ttl=24h \
    token_max_ttl=24h \
    secret_id_num_uses=40000 \
    policies="db"

echo "Creating role for approle.."
ROLE_ID=$(vault read auth/approle/role/my-role/role-id -format=json | jq -r '.data.role_id')
SECRET_ID=$(vault write -f auth/approle/role/my-role/secret-id -format=json | jq -r '.data.secret_id')