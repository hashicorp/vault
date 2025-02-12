# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Read and Verify SSH role configuration
resource "enos_remote_exec" "ssh_verify_role" {
  for_each = var.hosts

  environment = {
    ROLE_NAME         = var.create_state.ssh.role.name
    KEY_TYPE          = var.create_state.ssh.role.key_type
    DEFAULT_USER      = var.create_state.ssh.role.default_user
    PORT              = var.create_state.ssh.role.port
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = local.user_login_data.auth.client_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/ssh-verify-role.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# Read and Verify SSH CA configuration
resource "enos_remote_exec" "ssh_verify_ca" {
  for_each = var.hosts

  environment = {
    CA_KEY_TYPE       = var.create_state.ssh.ca.key_type
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = local.user_login_data.auth.client_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/ssh-verify-ca.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

// Read and Verify Signed SSH Key
resource "enos_remote_exec" "ssh_verify_signed_key" {
  for_each = var.hosts

  environment = {
    ROLE_NAME         = var.create_state.ssh.role.name
    PUBLIC_KEY_PATH   = var.public_key_path
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = local.user_login_data.auth.client_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/ssh-verify-signed-key.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

// Read and Verify OTP Credential
resource "enos_remote_exec" "ssh_verify_otp" {
  for_each = var.hosts

  environment = {
    OTP_ROLE_NAME     = var.create_state.ssh.otp_role.name
    TARGET_IP         = each.value.public_ip
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = local.user_login_data.auth.client_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/ssh-verify-otp.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
