# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0


terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

locals {
  instances        = toset([for idx in range(var.vault_instance_count) : tostring(idx)])
  expected_version = var.vault_edition == "oss" ? var.vault_product_version : "${var.vault_product_version}-ent"
}

resource "enos_remote_exec" "release_info" {
  for_each = local.instances

  environment = {
    VAULT_BIN_PATH = var.vault_bin_path
  }

  scripts = [abspath("${path.module}/scripts/get-status.sh")]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.vault_pods[each.key].name
      namespace         = var.vault_pods[each.key].namespace
    }
  }
}

resource "enos_local_exec" "smoke-verify-version" {
  for_each = enos_remote_exec.release_info

  environment = {
    VAULT_STATUS     = jsonencode(jsondecode(each.value.stdout).status)
    ACTUAL_VERSION   = jsondecode(each.value.stdout).version
    EXPECTED_VERSION = var.vault_product_version,
    VAULT_EDITION    = var.vault_edition,
    VAULT_REVISION   = var.vault_product_revision,
    CHECK_BUILD_DATE = var.check_build_date
    BUILD_DATE       = var.vault_build_date
  }

  scripts = [abspath("${path.module}/scripts/smoke-verify-version.sh")]
}
