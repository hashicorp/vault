# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Verify PKI Certificate
resource "enos_remote_exec" "pki_verify_certificates" {
  for_each = var.hosts

  environment = {
    MOUNT             = var.create_state.pki.mount
    AUTH_PATH         = "${var.create_state.auth.userpass.path}"
    USERNAME          = "${var.create_state.auth.userpass.user.name}"
    PASSWORD          = "${var.create_state.auth.userpass.user.password}"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
    COMMON_NAME       = var.create_state.pki.common_name
    ISSUER_NAME       = var.create_state.pki.issuer_name
    TTL               = var.create_state.pki.ttl
    TEST_DIR          = var.create_state.pki.test_dir
    VERIFY_PKI_CERTS  = var.verify_pki_certs
  }

  scripts = [abspath("${path.module}/../../scripts/pki-verify-certificates.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

