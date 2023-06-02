#!/bin/env sh

set -eux

sudo su "$SERVICE_USER" -c "VAULT_TOKEN=$VAULT_TOKEN VAULT_ADDR=$VAULT_ADDR $VAULT_BIN_PATH audit enable file file_path=$LOG_FILE_PATH"
