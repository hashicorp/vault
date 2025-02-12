# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  ssh_role_name  = "ssh_role"
  ssh_mount      = "ssh"
  ssh_key_types  = ["otp", "ca"]
  ssh_key_type   = local.ssh_key_types[random_integer.ssh_key_type_idx.result]
  ssh_test_ip    = "192.168.1.1"
  ssh_test_user  = "testuser"
  ssh_public_key = "ssh-rsa AAAAB3..."

  // Output
  ssh_output = {
    role_name   = local.ssh_role_name
    mount       = local.ssh_mount
    ca_key_type = local.ssh_key_type
    test = {
      ip   = local.ssh_test_ip
      user = local.ssh_test_user
    }
  }
}

resource "random_integer" "ssh_key_type_idx" {
  min = 0
  max = length(local.ssh_key_types) - 1
}

output "ssh" {
  value = local.ssh_output
}

# Enable SSH secrets engine
resource "enos_remote_exec" "secrets_enable_ssh" {
  environment = {
    ENGINE            = "ssh"
    MOUNT             = local.ssh_mount
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/secrets-enable.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Configure SSH CA
resource "enos_remote_exec" "ssh_configure_ca" {
  depends_on = [enos_remote_exec.secrets_enable_ssh]
  environment = {
    REQPATH           = "ssh/config/ca"
    PAYLOAD           = jsonencode({ key_type = local.ssh_ca_key_type })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Create SSH role
resource "enos_remote_exec" "ssh_create_role" {
  depends_on = [enos_remote_exec.ssh_configure_ca]
  environment = {
    REQPATH           = "ssh/roles/${local.ssh_role_name}"
    PAYLOAD           = jsonencode({ key_type = "ca", default_user = local.ssh_test_user, port = 22 })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Sign SSH key
resource "enos_remote_exec" "ssh_sign_key" {
  depends_on = [enos_remote_exec.ssh_create_role]
  environment = {
    REQPATH           = "ssh/sign/${local.ssh_role_name}"
    PAYLOAD           = jsonencode({ public_key = local.ssh_public_key })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Generate SSH OTP credential
resource "enos_remote_exec" "ssh_generate_otp" {
  depends_on = [enos_remote_exec.ssh_create_role]
  environment = {
    REQPATH           = "ssh/creds/${local.ssh_role_name}"
    PAYLOAD           = jsonencode({ ip = local.ssh_test_ip, username = local.ssh_test_user })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}
