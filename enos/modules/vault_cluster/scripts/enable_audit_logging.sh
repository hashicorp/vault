#!/bin/env sh
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


set -eux

$VAULT_BIN_PATH audit enable file file_path="$LOG_FILE_PATH"
