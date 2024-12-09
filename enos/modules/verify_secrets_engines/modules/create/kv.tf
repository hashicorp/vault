# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  group_name_kv_writers     = "kv_writers" # identity/group/name/kv_writers
  kv_mount                  = "secret"     # secret
  kv_write_policy_name      = "kv_writer"  # sys/policy/kv_writer
  kv_test_data_path_prefix  = "smoke"
  kv_test_data_value_prefix = "fire"
  kv_version                = 2

  // Response data
  identity_group_kv_writers_data = jsondecode(enos_remote_exec.identity_group_kv_writers.stdout).data

  // Output
  kv_output = {
    reader_group_name  = local.group_name_kv_writers
    writer_policy_name = local.kv_write_policy_name
    mount              = local.kv_mount
    version            = local.kv_version
    test = {
      path_prefix  = local.kv_test_data_path_prefix
      value_prefix = local.kv_test_data_value_prefix
    }
    data = {
      identity_group_kv_writers = local.identity_group_kv_writers_data
    }
  }
}

output "kv" {
  value = local.kv_output
}

# Enable kv secrets engine
resource "enos_remote_exec" "secrets_enable_kv_secret" {
  environment = {
    ENGINE            = "kv"
    MOUNT             = local.kv_mount
    SECRETS_META      = "-version=${local.kv_version}"
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

# Create a group policy that allows writing to our kv store
resource "enos_remote_exec" "policy_write_kv_writer" {
  depends_on = [
    enos_remote_exec.secrets_enable_kv_secret,
  ]
  environment = {
    POLICY_NAME       = local.kv_write_policy_name
    POLICY_CONFIG     = <<-EOF
      path "${local.kv_mount}/*" {
        capabilities = ["create", "update", "read", "delete", "list"]
      }
    EOF
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/policy-write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Create kv_writers group and add our testuser to it
resource "enos_remote_exec" "identity_group_kv_writers" {
  environment = {
    REQPATH = "identity/group"
    PAYLOAD = jsonencode({
      member_entity_ids = [local.user_entity_data.id], // Created in identity.tf
      name              = local.group_name_kv_writers,
      policies          = [local.kv_write_policy_name],
    })
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

// Write test data as our user.
resource "enos_remote_exec" "kv_put_secret_test" {
  depends_on = [
    enos_remote_exec.secrets_enable_kv_secret,
  ]
  for_each = var.hosts

  environment = {
    MOUNT             = local.kv_mount
    SECRET_PATH       = "${local.kv_test_data_path_prefix}-${each.key}"
    KEY               = "${local.kv_test_data_path_prefix}-${each.key}"
    VALUE             = "${local.kv_test_data_value_prefix}-${each.key}"
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = local.user_login_data.auth.client_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/kv-put.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
