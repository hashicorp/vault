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
  vault_server_pod   = var.vault_pods[0].name
  vault_namespace    = var.vault_pods[0].namespace
  test_app_namespace = "default"
  app_pod_name       = "vault-client-app"
}

# 1. Enable AppRole auth, write policy, role, and KV secret inside the Vault
# server pod. The script prints "role_id|secret_id" to stdout so the caller
# can split on "|" without extra round-trips.
resource "enos_remote_exec" "setup_vault_approle" {
  environment = {
    VAULT_TOKEN    = var.vault_root_token
    VAULT_BIN_PATH = var.vault_bin_path
  }

  scripts = [abspath("${path.module}/scripts/setup-vault-approle.sh")]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = local.vault_server_pod
      namespace         = local.vault_namespace
    }
  }
}

locals {
  role_id   = split("|", enos_remote_exec.setup_vault_approle.stdout)[0]
  secret_id = split("|", enos_remote_exec.setup_vault_approle.stdout)[1]
}

# 2. Deploy the Kubernetes Secret, ConfigMap, and Pod with vault-client sidecar.
# kubectl is not installed on the runner; kind control-plane containers always
# have it, so all kubectl calls are proxied through docker exec.
resource "enos_local_exec" "deploy_sidecar_pod" {
  depends_on = [enos_remote_exec.setup_vault_approle]

  environment = {
    ROLE_ID           = local.role_id
    SECRET_ID         = local.secret_id
    CLIENT_IMAGE      = "${var.client_image_repository}:${var.client_image_tag}"
    APP_POD_NAME      = local.app_pod_name
    KIND_CLUSTER_NAME = var.kind_cluster_name
  }

  scripts = [abspath("${path.module}/scripts/deploy-sidecar-pod.sh")]
}

output "app_pod_name" {
  value = local.app_pod_name
}

output "app_pod_namespace" {
  value = local.test_app_namespace
}
