# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  ldap_output = {
    ip_version        = var.ip_version
    ldap_mount        = "ldap"
    host              = var.ldap_host
    port              = var.ldap_port
    username          = "enos"
    pw                = var.ldap_password
    vault_policy_name = local.kv_output.writer_policy_name
  }
}

output "ldap" {
  value = local.ldap_output
}

# Enable LDAP secrets engine
resource "enos_remote_exec" "secrets_enable_ldap_secret" {
  environment = {
    ENGINE            = local.ldap_output.ldap_mount
    MOUNT             = local.ldap_output.ldap_mount
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../scripts/secrets-enable.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Configuring Openldap Server and Vault LDAP
resource "enos_remote_exec" "ldap_configurations" {
  depends_on = [
    enos_remote_exec.policy_write_kv_writer,
    enos_remote_exec.secrets_enable_ldap_secret
  ]

  environment = {
    MOUNT             = local.ldap_output.ldap_mount
    LDAP_SERVER       = var.ldap_host.private_ip
    LDAP_PORT         = local.ldap_output.port
    LDAP_USERNAME     = local.ldap_output.username
    LDAP_ADMIN_PW     = local.ldap_output.pw
    POLICY_NAME       = local.ldap_output.vault_policy_name
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../scripts/ldap-configs.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}
