# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  user_login_data = jsondecode(enos_remote_exec.auth_login_testuser.stdout)
}

resource "enos_remote_exec" "auth_login_testuser" {
  environment = {
    AUTH_PATH         = var.create_state.auth.userpass.path
    PASSWORD          = var.create_state.auth.userpass.user.password
    USERNAME          = var.create_state.auth.userpass.user.name
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/auth-userpass-login.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}
