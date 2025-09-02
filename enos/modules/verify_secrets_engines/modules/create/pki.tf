# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  pki_mount       = "pki" # secret
  pki_issuer_name = "issuer"
  pki_common_name = "common"
  pki_default_ttl = "72h"
  pki_test_dir    = "tmp-test-results"

  // Output
  pki_output = {
    common_name = local.pki_common_name
    issuer_name = local.pki_issuer_name
    mount       = local.pki_mount
    ttl         = local.pki_default_ttl
    test_dir    = local.pki_test_dir
  }

}

output "pki" {
  value = local.pki_output
}

# Enable pki secrets engine
resource "enos_remote_exec" "secrets_enable_pki_secret" {
  environment = {
    ENGINE            = local.pki_mount
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
resource "enos_remote_exec" "pki_issue_certificates" {
  depends_on = [enos_remote_exec.secrets_enable_pki_secret]
  for_each   = var.hosts

  environment = {
    MOUNT             = local.pki_mount
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
    COMMON_NAME       = local.pki_common_name
    ISSUER_NAME       = local.pki_issuer_name
    TTL               = local.pki_default_ttl
    TEST_DIR          = local.pki_test_dir
  }

  scripts = [abspath("${path.module}/../../scripts/pki-issue-certificates.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
