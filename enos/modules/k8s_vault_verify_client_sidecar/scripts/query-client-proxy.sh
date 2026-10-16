#!/bin/sh
# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

set -eu

# The Vault Agent proxy listener runs on localhost with use_auto_auth_token=true,
# so it injects the auto-auth token into upstream requests automatically.
# No explicit X-Vault-Token header is needed from the app container.
# KV v1 path: /v1/secret/<key>. Update to /v1/secret/data/<key> if switching to KV v2.
curl -s -f http://127.0.0.1:8200/v1/secret/client-sidecar
