# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  pki_mount                 = "pki_secret" # secret
  pki_issuer_name           = "issuer"
  pki_common_name           = "common"
  pki_default_ttl           = "72h"
  pki_test_data_path_prefix = "smoke"
  tmp_test_results          = "tmp_test_results"

  // Output
  pki_output = {
    mount        = local.pki_mount
    common_name  = local.pki_common_name
    test_results = local.tmp_test_results
  }

  test = {
    path_prefix = local.pki_test_data_path_prefix
  }
}

output "pki" {
  value = local.pki_output
}

# Verify PKI Certificate
resource "enos_remote_exec" "pki_verify_certificates" {
  for_each = var.hosts

  environment = {
    MOUNT             = "pki_secret"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
    COMMON_NAME       = "common"
    TTL               = "72h"
    TMP_TEST_RESULTS  = "tmp_test_results"
  }

  scripts = [abspath("${path.module}/../../scripts/kv-pki-verify-certificates.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

