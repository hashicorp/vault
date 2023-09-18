# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


terraform {
  required_providers {
    enos = {
      version = ">= 0.1.17"
      source  = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

locals {
  instances = toset([for idx in range(var.vault_instance_count) : tostring(idx)])
}

resource "enos_remote_exec" "curl_ui" {
  for_each = local.instances

  inline = [
    "curl -s -o /dev/null -w '%%{redirect_url}' http://localhost:8200/",
    "curl -s -o /dev/null -Iw '%%{http_code}\n' http://localhost:8200/ui/"
  ]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.vault_pods[each.key].name
      namespace         = var.vault_pods[each.key].namespace
    }
  }
}

resource "enos_local_exec" "verify_ui" {
  for_each = enos_remote_exec.curl_ui

  environment = {
    REDIRECT_URL  = split("\n", each.value.stdout)[0]
    UI_URL_RESULT = split("\n", each.value.stdout)[1]
  }

  scripts = [abspath("${path.module}/scripts/smoke-verify-ui.sh")]
}
