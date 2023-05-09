#!/bin/env sh

set -eux

LOG_DIR=$(dirname "$LOG_FILE_PATH")

sudo mkdir -p "$LOG_DIR"
sudo chown vault:vault "$LOG_DIR"

"$VAULT_BIN_PATH" audit enable file file_path="$LOG_FILE_PATH"
