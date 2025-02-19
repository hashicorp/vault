# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

resource "enos_remote_exec" "kv_get_verify_test_data" {
  for_each = var.hosts

  environment = {
    MOUNT             = var.create_state.kv.mount
    SECRET_PATH       = "${var.create_state.kv.test.path_prefix}-${each.key}"
    KEY               = "${var.create_state.kv.test.path_prefix}-${each.key}"
    KV_VERSION        = var.create_state.kv.version
    VALUE             = "${var.create_state.kv.test.value_prefix}-${each.key}"
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = local.user_login_data.auth.client_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/kv-verify-value.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
