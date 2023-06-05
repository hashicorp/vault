#!/bin/env sh

set -eux

$VAULT_BIN_PATH audit enable file file_path="$LOG_FILE_PATH"
