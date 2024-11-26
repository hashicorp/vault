# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Verify PKI Certificate
resource "enos_remote_exec" "pki_issue_certificates" {
  depends_on = [enos_remote_exec.pki_issue_certificates]

  environment = {
    MOUNT             = local.pki_mount
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
    COMMON_NAME       = local.pki_common_name
    TTL               = local.pki_default_ttl
    TMP_TEST_RESULTS  = local.pki_tmp_results
  }

  scripts = [abspath("${path.module}/../../scripts/kv-pki-verify-certificates.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

