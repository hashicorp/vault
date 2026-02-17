// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

// Delete LDAP library set
// Test Case: Delete Library Set - Delete a library & all associated service accounts
resource "enos_remote_exec" "ldap_library_set_delete" {
  count = var.ldap_enabled ? 1 : 0

  environment = {
    REQPATH           = "${try(var.create_state.ldap, null) != null ? var.create_state.ldap.ldap_mount : "ldap"}/library/test-set"
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/delete.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}


