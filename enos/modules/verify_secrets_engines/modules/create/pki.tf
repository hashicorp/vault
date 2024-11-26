# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  pki_mount                  = "pki_secret"     # secret
  pki_issuer_name            = "issuer"
  pki_common_name            = "common"
  pki_default_ttl            = "72h"
  pki_test_data_path_prefix  = "smoke"
  tmp_test_results           = "tmp_test_results"

  // Response data
#   identity_group_kv_writers_data = jsondecode(enos_remote_exec.identity_group_kv_writers.stdout).data

  // Output
  pki_output = {
    mount              = local.pki_mount
    issuer_name        = local.pki_issuer_name
    common_name        = local.pki_common_name
  }
  test = {
    path_prefix  = local.pki_test_data_path_prefix
  }
}

output "pki" {
  value = local.pki_output
}

# Enable pki secrets engine
resource "enos_remote_exec" "secrets_enable_pki_secret" {
  environment = {
    ENGINE            = "pki"
    MOUNT             = local.pki_mount
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

# Issue RSA Certificate
resource "enos_remote_exec" "pki_issue_rsa_cert" {
  depends_on = [enos_remote_exec.secrets_enable_pki_secret]

  environment = {
    MOUNT             = local.pki_mount
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
    COMMON_NAME       = local.pki_common_name
    TTL               = local.pki_default_ttl
    TMP_TEST_RESULTS  = local.tmp_test_results
  }

  scripts = [abspath("${path.module}/../../scripts/kv-pki-issue-certificates.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# # Configure AIA
# resource "enos_remote_exec" "policy_write_kv_writer" {
#   depends_on = [
#     enos_remote_exec.secrets_enable_kv_secret,
#   ]
#   environment = {
#     POLICY_NAME       = local.kv_write_policy_name
#     POLICY_CONFIG     = <<-EOF
#       path "${local.kv_mount}/*" {
#         capabilities = ["create", "update", "read", "delete", "list"]
#       }
#     EOF
#     VAULT_ADDR        = var.vault_addr
#     VAULT_TOKEN       = var.vault_root_token
#     VAULT_INSTALL_DIR = var.vault_install_dir
#   }
#
#   scripts = [abspath("${path.module}/../../scripts/policy-write.sh")]
#
#   transport = {
#     ssh = {
#       host = var.leader_host.public_ip
#     }
#   }
# }
