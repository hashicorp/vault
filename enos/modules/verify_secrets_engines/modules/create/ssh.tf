# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  otp_role_name = "ssh_role_otp"
  otp_role_params = {
    key_type               = "otp"
    default_user           = local.ssh_test_user
    default_user_template  = false
    allowed_users          = local.ssh_test_user
    allowed_users_template = false
    cidr_list              = local.ssh_test_ip
    exclude_cidr_list      = "10.0.0.0/8"
    port                   = 22
    ttl                    = "30m"
    max_ttl                = "1h"
  }

  ca_role_name = "ssh_role_ca"
  ca_role_params = {
    key_type                = "ca"
    default_user            = local.ssh_test_user
    default_user_template   = false
    allow_user_certificates = true
    allow_host_certificates = true
    allowed_users           = local.ssh_test_user
    allowed_users_template  = false
    port                    = 22
    ttl                     = "1h"
    max_ttl                 = "2h"
    key_id_format           = "custom-keyid-{{token_display_name}}"
    allowed_extensions      = "*"
    default_extensions = {
      "permit-pty" = ""
    }
    allow_user_key_ids     = true
    allow_empty_principals = false
    algorithm_signer       = "default"
  }

  is_fips = can(regex("fips", lower(var.vault_edition)))

  ca_key_types = local.is_fips ? [
    "ssh-rsa", "ecdsa-sha2-nistp256", "ecdsa-sha2-nistp384", "ecdsa-sha2-nistp521"
    ] : [
    "ssh-rsa", "ecdsa-sha2-nistp256", "ecdsa-sha2-nistp384", "ecdsa-sha2-nistp521", "ssh-ed25519"
  ]

  ssh_mount      = "ssh"
  ca_key_type    = local.ca_key_types[random_integer.ca_key_type_idx.result]
  cert_key_types = ["rsa", "ed25519", "ec"]
  cert_key_type  = local.cert_key_types[random_integer.cert_key_idx.result]
  ssh_test_ip    = "192.168.1.1/32"
  ssh_test_user  = "testuser"
  ssh_public_key = tls_private_key.test_ssh_key.public_key_openssh

  # Map ca_key_type to valid Terraform tls_private_key algorithm and curve
  key_algorithm_map = {
    "ssh-rsa"             = "RSA"
    "ecdsa-sha2-nistp256" = "ECDSA"
    "ecdsa-sha2-nistp384" = "ECDSA"
    "ecdsa-sha2-nistp521" = "ECDSA"
    "ssh-ed25519"         = "ED25519"
  }

  ecdsa_curve_map = {
    "ecdsa-sha2-nistp256" = "P256"
    "ecdsa-sha2-nistp384" = "P384"
    "ecdsa-sha2-nistp521" = "P521"
  }

  # Extract the corresponding algorithm and curve
  key_algorithm = lookup(local.key_algorithm_map, local.ca_key_type, "RSA")
  ecdsa_curve   = lookup(local.ecdsa_curve_map, local.ca_key_type, null)

  // Response data
  ssh_sign_key_data      = jsondecode(enos_remote_exec.ssh_sign_key.stdout).data
  ssh_generate_otp_data  = jsondecode(enos_remote_exec.ssh_generate_otp.stdout).data
  ssh_generate_cert_data = jsondecode(enos_remote_exec.ssh_generate_cert.stdout).data

  // Output
  ssh_output = {
    ca_role_name    = local.ca_role_name
    otp_role_name   = local.otp_role_name
    mount           = local.ssh_mount
    ca_key_type     = local.ca_key_type
    cert_key_type   = local.cert_key_type
    test_ip         = local.ssh_test_ip
    test_user       = local.ssh_test_user
    otp_role_params = local.otp_role_params
    ca_role_params  = local.ca_role_params
    data = {
      sign_key      = local.ssh_sign_key_data
      generate_otp  = local.ssh_generate_otp_data
      generate_cert = local.ssh_generate_cert_data
    }
  }
}

resource "tls_private_key" "test_ssh_key" {
  algorithm = local.key_algorithm

  # Conditionally set ecdsa_curve only for ECDSA keys
  ecdsa_curve = local.key_algorithm == "ECDSA" ? local.ecdsa_curve : null

  rsa_bits = local.key_algorithm == "RSA" ? 2048 : null
}

resource "random_integer" "ca_key_type_idx" {
  min = 0
  max = length(local.ca_key_types) - 1
}

resource "random_integer" "cert_key_idx" {
  min = 0
  max = length(local.cert_key_types) - 1
}

output "ssh" {
  value = local.ssh_output
}

# Enable SSH secrets engine
resource "enos_remote_exec" "secrets_enable_ssh" {
  environment = {
    ENGINE            = local.ssh_mount
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
    PAYLOAD           = jsonencode({ key_type = local.ca_key_type })
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

resource "enos_remote_exec" "ssh_create_ca_role" {
  depends_on = [enos_remote_exec.ssh_configure_ca]
  environment = {
    REQPATH           = "ssh/roles/${local.ca_role_name}"
    PAYLOAD           = jsonencode(local.ca_role_params)
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

# Create SSH OTP role
resource "enos_remote_exec" "ssh_create_otp_role" {
  depends_on = [enos_remote_exec.secrets_enable_ssh]
  environment = {
    REQPATH           = "ssh/roles/${local.otp_role_name}"
    PAYLOAD           = jsonencode(local.otp_role_params)
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
  depends_on = [enos_remote_exec.ssh_create_ca_role]
  environment = {
    REQPATH           = "ssh/sign/ssh_role_ca"
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
  depends_on = [enos_remote_exec.ssh_create_otp_role]
  environment = {
    REQPATH           = "ssh/creds/ssh_role_otp"
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

# Generate SSH Certificate and Key
resource "enos_remote_exec" "ssh_generate_cert" {
  depends_on = [enos_remote_exec.ssh_create_ca_role]

  environment = {
    REQPATH           = "ssh/issue/ssh_role_ca"
    PAYLOAD           = jsonencode({ key_type = local.cert_key_type })
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