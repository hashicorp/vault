# Copyright IBM Corp. 2016, 2026
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_version = ">= 1.0"

  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

locals {
  expected_secret = "hello-from-client-sidecar"
}

# 1. Fetch the rendered secret file written by the Vault Client sidecar to the
# shared volume and read it from the app container.
resource "enos_remote_exec" "get_rendered_secret" {
  scripts = [abspath("${path.module}/scripts/get-rendered-secret.sh")]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.app_pod_name
      namespace         = var.app_pod_namespace
      container         = "app"
    }
  }
}

# 2. Assert the rendered content matches what was seeded into Vault.
resource "enos_local_exec" "assert_rendered_secret" {
  depends_on = [enos_remote_exec.get_rendered_secret]

  environment = {
    ACTUAL_CONTENT   = enos_remote_exec.get_rendered_secret.stdout
    EXPECTED_CONTENT = local.expected_secret
    LABEL            = "rendered secret"
  }

  scripts = [abspath("${path.module}/scripts/verify-content.sh")]
}

# 3. Query the Vault Client proxy listener over localhost from the app container.
resource "enos_remote_exec" "query_client_proxy" {
  depends_on = [enos_remote_exec.get_rendered_secret]

  scripts = [abspath("${path.module}/scripts/query-client-proxy.sh")]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.app_pod_name
      namespace         = var.app_pod_namespace
      container         = "app"
    }
  }
}

# 4. Assert the proxy returned the expected secret data.
resource "enos_local_exec" "assert_proxy_response" {
  depends_on = [enos_remote_exec.query_client_proxy]

  environment = {
    ACTUAL_CONTENT   = enos_remote_exec.query_client_proxy.stdout
    EXPECTED_CONTENT = local.expected_secret
    LABEL            = "proxy response"
  }

  scripts = [abspath("${path.module}/scripts/verify-content.sh")]
}
