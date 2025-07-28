# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  ldap_login_data = jsondecode(enos_remote_exec.ldap_verify_configs.stdout)
}

# Verifying Vault LDAP Configurations
resource "enos_remote_exec" "ldap_verify_configs" {

  environment = {
    MOUNT             = "${var.create_state.ldap.ldap_mount}"
    LDAP_SERVER       = "${var.create_state.ldap.host.private_ip}"
    LDAP_PORT         = "${var.create_state.ldap.port}"
    LDAP_USERNAME     = "${var.create_state.ldap.username}"
    LDAP_ADMIN_PW     = "${var.create_state.ldap.pw}"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../scripts/ldap-verify-configs")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}
