#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -ex

make dockerbuild

docker kill vaultplg 2>/dev/null || true
tmpdir=$(mktemp -d vaultplgXXXXXX)
mkdir "$tmpdir/data"
PLUGIN_DIR="$(pwd)/pkg/linux_amd64"
docker run --rm -d -p8200:8200 --name vaultplg -v "$(pwd)/$tmpdir/data":/data -v "${PLUGIN_DIR}":/example --cap-add=IPC_LOCK -e 'VAULT_LOCAL_CONFIG=
{
  "backend": {"file": {"path": "/data"}},
  "listener": [{"tcp": {"address": "0.0.0.0:8200", "tls_disable": true}}],
  "plugin_directory": "/example",
  "log_level": "debug",
  "disable_mlock": true,
  "api_addr": "http://localhost:8200"
}
' docker.mirror.hashicorp.services/vault server
sleep 1

export VAULT_ADDR=http://localhost:8200

initoutput=$(vault operator init -key-shares=1 -key-threshold=1 -format=json)
vault operator unseal $(echo "$initoutput" | jq -r .unseal_keys_hex[0])

export VAULT_TOKEN=$(echo "$initoutput" | jq -r .root_token)

vault write sys/plugins/catalog/secret/vault-plugin-secrets-mongodbatlas \
    sha_256=$(shasum -a 256 "${PLUGIN_DIR}/vault-plugin-secrets-mongodbatlas" | cut -d' ' -f1) \
    command="vault-plugin-secrets-mongodbatlas"

vault secrets enable \
    -path="mongodbatlas" \
    -plugin-name="vault-plugin-secrets-mongodbatlas" plugin
