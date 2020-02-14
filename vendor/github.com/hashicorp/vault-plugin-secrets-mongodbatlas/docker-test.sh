#!/usr/bin/env bash

set -ex

make dockerbuild

docker kill vaultplg 2>/dev/null || true
tmpdir=$(mktemp -d vaultplgXXXXXX)
mkdir "$tmpdir/data"
docker run --rm -d -p8200:8200 --name vaultplg -v "$(pwd)/$tmpdir/data":/data -v $(pwd)/bin:/example --cap-add=IPC_LOCK -e 'VAULT_LOCAL_CONFIG=
{
  "backend": {"file": {"path": "/data"}},
  "listener": [{"tcp": {"address": "0.0.0.0:8200", "tls_disable": true}}],
  "plugin_directory": "/example",
  "log_level": "debug",
  "disable_mlock": true,
  "api_addr": "http://localhost:8200"
}
' vault server
sleep 1

export VAULT_ADDR=http://localhost:8200

initoutput=$(vault operator init -key-shares=1 -key-threshold=1 -format=json)
vault operator unseal $(echo "$initoutput" | jq -r .unseal_keys_hex[0])

export VAULT_TOKEN=$(echo "$initoutput" | jq -r .root_token)

vault write sys/plugins/catalog/secret/vault-plugin-secrets-mongodbatlas \
    sha_256=$(shasum -a 256 bin/vault-plugin-secrets-mongodbatlas | cut -d' ' -f1) \
    command="vault-plugin-secrets-mongodbatlas"

vault secrets enable \
    -path="mongodbatlas" \
    -plugin-name="vault-plugin-secrets-mongodbatlas" plugin

vault write sys/plugins/catalog/database/mongodbatlas-database-plugin \
    sha256=$(shasum -a 256 bin/mongodbatlas-database-plugin | cut -d' ' -f1) \
    command="mongodbatlas-database-plugin"

vault secrets enable database
