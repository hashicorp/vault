#!/bin/env sh
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


set -eux

$VAULT_BIN_PATH audit enable file file_path="$LOG_FILE_PATH"
