# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "kmip_listen_address" {
  type        = string
  description = "The KMIP listen address for the Vault server"
  default     = "0.0.0.0"
}

locals {
  kmip_scope_name  = "kmip_scope"
  kmip_role_name   = "kmip_role"
  kmip_cert_format = "pem"
  kmip_mount_path  = "kmip"

  // Response data - only access if Vault Enterprise (count > 0)
  server_ca   = var.vault_edition == "ce" ? "" : enos_remote_exec.kmip_configure[0].stdout
  client_cert = var.vault_edition == "ce" ? "" : enos_remote_exec.kmip_generate_certificate[0].stdout

  kmip_output = {
    server_ca      = local.server_ca
    client_cert    = local.client_cert
    test_server_ip = var.integration_host_state.kmip.host.public_ip
    port           = var.ports.kmip.port
  }
}

output "kmip" {
  value = local.kmip_output
}

resource "enos_remote_exec" "secrets_enable_kmip_secret" {
  environment = {
    ENGINE            = "kmip"
    MOUNT             = local.kmip_mount_path
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  // Only perform KMIP operations for Vault Enterprise
  // The KMIP secrets engine is not available in Vault CE
  count = var.vault_edition == "ce" ? 0 : 1

  scripts = [abspath("${path.module}/../../scripts/secrets-enable.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

resource "enos_remote_exec" "kmip_configure" {
  depends_on = [enos_remote_exec.secrets_enable_kmip_secret]
  environment = {
    MOUNT             = local.kmip_mount_path
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
    KMIP_MOUNT        = local.kmip_mount_path
    KMIP_LISTEN_ADDR  = var.kmip_listen_address
    KMIP_PORT         = var.ports.kmip.port
  }

  // Only perform KMIP operations for Vault Enterprise
  // The KMIP secrets engine is not available in Vault CE
  count = var.vault_edition == "ce" ? 0 : 1

  scripts = [abspath("${path.module}/../../scripts/kmip/kmip-configure.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Creating KMIP Scope
resource "enos_remote_exec" "kmip_create_scope" {
  depends_on = [enos_remote_exec.kmip_configure]

  environment = {
    REQPATH           = "${local.kmip_mount_path}/scope/${local.kmip_scope_name}"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  // Only perform KMIP operations for Vault Enterprise
  // The KMIP secrets engine is not available in Vault CE
  count = var.vault_edition == "ce" ? 0 : 1

  scripts = [abspath("${path.module}/../../scripts/write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Creating KMIP Role
resource "enos_remote_exec" "kmip_create_role" {
  depends_on = [enos_remote_exec.kmip_create_scope]

  environment = {
    REQPATH = "${local.kmip_mount_path}/scope/${local.kmip_scope_name}/role/${local.kmip_role_name}"
    PAYLOAD = jsonencode({
      operation_all = true,
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  // Only perform KMIP operations for Vault Enterprise
  // The KMIP secrets engine is not available in Vault CE
  count = var.vault_edition == "ce" ? 0 : 1

  scripts = [abspath("${path.module}/../../scripts/write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Generating KMIP Certificate
resource "enos_remote_exec" "kmip_generate_certificate" {
  depends_on = [enos_remote_exec.kmip_create_role]

  environment = {
    MOUNT             = local.kmip_mount_path
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
    SCOPE_NAME        = local.kmip_scope_name
    ROLE_NAME         = local.kmip_role_name
    CERT_FORMAT       = local.kmip_cert_format
  }

  // Only perform KMIP operations for Vault Enterprise
  // The KMIP secrets engine is not available in Vault CE
  count = var.vault_edition == "ce" ? 0 : 1

  scripts = [abspath("${path.module}/../../scripts/kmip/kmip-generate-cert.sh")]
  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Managing KMIP Roles
resource "enos_remote_exec" "kmip_manage_roles" {
  depends_on = [enos_remote_exec.kmip_generate_certificate]
  environment = {
    REQPATH = "${local.kmip_mount_path}/scope/${local.kmip_scope_name}/role/${local.kmip_role_name}"
    PAYLOAD = jsonencode({
      operation_activate = true,
      operation_create   = true,
      operation_get      = true
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  // Only perform KMIP operations for Vault Enterprise
  // The KMIP secrets engine is not available in Vault CE
  count = var.vault_edition == "ce" ? 0 : 1

  scripts = [abspath("${path.module}/../../scripts/write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}
