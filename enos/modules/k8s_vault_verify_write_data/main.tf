# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

locals {
  instances = toset([for idx in range(var.vault_instance_count) : tostring(idx)])
}

resource "enos_remote_exec" "smoke-enable-secrets-kv" {
  environment = {
    VAULT_TOKEN = var.vault_root_token
  }

  inline = ["${var.vault_bin_path} secrets enable -path=\"secret\" kv"]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.vault_pods[0].name
      namespace         = var.vault_pods[0].namespace
    }
  }
}

# Verify that we can enable the k/v secrets engine and write data to it.
resource "enos_remote_exec" "smoke-write-test-data" {
  depends_on = [enos_remote_exec.smoke-enable-secrets-kv]
  for_each   = local.instances

  environment = {
    VAULT_TOKEN = var.vault_root_token
  }

  inline = ["${var.vault_bin_path} kv put secret/test smoke${each.key}=fire"]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.vault_pods[each.key].name
      namespace         = var.vault_pods[each.key].namespace
    }
  }
}
