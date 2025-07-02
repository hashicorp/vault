# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


# List SSH roles
resource "enos_remote_exec" "ssh_list_roles" {
  for_each = var.hosts
  environment = {
    REQPATH           = "ssh/roles"
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/list.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# Read and Verify SSH CA Role configuration
resource "enos_remote_exec" "ssh_verify_ca_role" {
  for_each = var.hosts

  environment = {
    ROLE_NAME               = var.create_state.ssh.ca_role_name
    KEY_TYPE                = var.create_state.ssh.ca_role_params.key_type
    DEFAULT_USER            = var.create_state.ssh.ca_role_params.default_user
    ALLOWED_USERS           = var.create_state.ssh.ca_role_params.allowed_users
    PORT                    = tostring(var.create_state.ssh.ca_role_params.port)
    TTL                     = var.create_state.ssh.ca_role_params.ttl
    MAX_TTL                 = var.create_state.ssh.ca_role_params.max_ttl
    ALLOWED_EXTENSIONS      = var.create_state.ssh.ca_role_params.allowed_extensions
    ALLOW_USER_CERTIFICATES = tostring(var.create_state.ssh.ca_role_params.allow_user_certificates)
    ALLOW_HOST_CERTIFICATES = tostring(var.create_state.ssh.ca_role_params.allow_host_certificates)
    ALLOW_USER_KEY_IDS      = tostring(var.create_state.ssh.ca_role_params.allow_user_key_ids)
    ALLOW_EMPTY_PRINCIPALS  = tostring(var.create_state.ssh.ca_role_params.allow_empty_principals)
    ALGORITHM_SIGNER        = var.create_state.ssh.ca_role_params.algorithm_signer
    KEY_ID_FORMAT           = var.create_state.ssh.ca_role_params.key_id_format
    DEFAULT_EXTENSIONS      = jsonencode(var.create_state.ssh.ca_role_params.default_extensions)
    VAULT_ADDR              = var.vault_addr
    VAULT_TOKEN             = var.vault_root_token
    VAULT_INSTALL_DIR       = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/ssh-verify-role.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "ssh_verify_otp_role" {
  for_each = var.hosts

  environment = {
    ROLE_NAME         = var.create_state.ssh.otp_role_name
    KEY_TYPE          = var.create_state.ssh.otp_role_params.key_type
    DEFAULT_USER      = var.create_state.ssh.otp_role_params.default_user
    ALLOWED_USERS     = var.create_state.ssh.otp_role_params.allowed_users
    CIDR_LIST         = var.create_state.ssh.otp_role_params.cidr_list
    EXCLUDE_CIDR_LIST = var.create_state.ssh.otp_role_params.exclude_cidr_list
    PORT              = tostring(var.create_state.ssh.otp_role_params.port)
    TTL               = var.create_state.ssh.otp_role_params.ttl
    MAX_TTL           = var.create_state.ssh.otp_role_params.max_ttl
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
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
    CA_KEY_TYPE       = var.create_state.ssh.ca_key_type
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
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
    SIGNED_KEY        = var.create_state.ssh.data.sign_key.signed_key
    CA_KEY_TYPE       = var.create_state.ssh.ca_key_type
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
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
    OTP               = var.create_state.ssh.data.generate_otp.key
    IP                = var.create_state.ssh.test_ip
    ROLE_NAME         = var.create_state.ssh.otp_role_name
    USERNAME          = var.create_state.ssh.test_user
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/ssh-verify-otp.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_local_exec" "ssh_verify_cert" {
  environment = {
    SIGNED_KEY        = var.create_state.ssh.data.generate_cert.signed_key
    CERT_KEY_TYPE     = var.create_state.ssh.cert_key_type
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/ssh-verify-signed-key.sh")]
}
