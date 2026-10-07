# Copyright IBM Corp. 2016, 2025
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
  # The vault-client-sidecar container lives inside the app pod deployed by
  # k8s_deploy_client_sidecar. All smoke commands exec into that container
  # directly — it is already running inside the cluster, so KUBERNETES_SERVICE_HOST
  # and KUBERNETES_SERVICE_PORT are automatically set by Kubernetes, which is
  # exactly what the kubeclient binary requires.
  sidecar_container = "vault-client-sidecar"

  # Skip the binary-level smoke tests when the caller is still using the
  # vault-enterprise image as a stand-in (kubeclient_bin_path == /bin/vault).
  # The vault binary does not speak the kubeclient CLI flags, so running them
  # would produce misleading failures.
  # Flip skip_smoke to false the moment a real kubeclient path is supplied.
  skip_smoke = var.kubeclient_bin_path == "/bin/vault"
}

# ---------------------------------------------------------------------------
# Step 1 — RBAC: grant the app pod's default service account GET + PATCH on
# pods so the kubeclient binary can reach the Kubernetes API server.
# Runs unconditionally and is idempotent (kubectl apply).
# ---------------------------------------------------------------------------
resource "enos_local_exec" "setup_rbac" {
  environment = {
    KIND_CLUSTER_NAME = var.kind_cluster_name
    TARGET_NAMESPACE  = var.app_pod_namespace
  }

  scripts = [abspath("${path.module}/scripts/setup-rbac.sh")]
}

# ---------------------------------------------------------------------------
# Step 2 — Smoke: get-pod happy path
#
# Execs into the vault-client-sidecar container and calls:
#   <kubeclient_bin_path> -call=get-pod -namespace=<ns> -pod-name=<pod>
# Asserts that stdout contains JSON with the pod name.
# ---------------------------------------------------------------------------
resource "enos_remote_exec" "smoke_get_pod" {
  count      = local.skip_smoke ? 0 : 1
  depends_on = [enos_local_exec.setup_rbac]

  environment = {
    KUBECLIENT_BIN  = var.kubeclient_bin_path
    SMOKE_NAMESPACE = var.app_pod_namespace
    SMOKE_POD       = var.app_pod_name
  }

  scripts = [abspath("${path.module}/scripts/smoke-get-pod.sh")]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.app_pod_name
      namespace         = var.app_pod_namespace
      container         = local.sidecar_container
    }
  }
}

# ---------------------------------------------------------------------------
# Step 3 — Smoke: patch-pod single label (add)
#
# Adds /metadata/labels/smoke-test=passed, then re-fetches the pod and
# asserts the label is present in the returned JSON.
# ---------------------------------------------------------------------------
resource "enos_remote_exec" "smoke_patch_pod_single" {
  count      = local.skip_smoke ? 0 : 1
  depends_on = [enos_remote_exec.smoke_get_pod]

  environment = {
    KUBECLIENT_BIN  = var.kubeclient_bin_path
    SMOKE_NAMESPACE = var.app_pod_namespace
    SMOKE_POD       = var.app_pod_name
  }

  scripts = [abspath("${path.module}/scripts/smoke-patch-pod.sh")]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.app_pod_name
      namespace         = var.app_pod_namespace
      container         = local.sidecar_container
    }
  }
}

# ---------------------------------------------------------------------------
# Step 4 — Smoke: patch-pod multiple labels (comma-separated -patches flag)
#
# Sends two patches in a single call, then verifies both labels appear in
# the get-pod response. Tests the comma-parsing path in main.go.
# ---------------------------------------------------------------------------
resource "enos_remote_exec" "smoke_patch_pod_multi" {
  count      = local.skip_smoke ? 0 : 1
  depends_on = [enos_remote_exec.smoke_patch_pod_single]

  environment = {
    KUBECLIENT_BIN  = var.kubeclient_bin_path
    SMOKE_NAMESPACE = var.app_pod_namespace
    SMOKE_POD       = var.app_pod_name
  }

  scripts = [abspath("${path.module}/scripts/smoke-patch-pod-multi.sh")]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.app_pod_name
      namespace         = var.app_pod_namespace
      container         = local.sidecar_container
    }
  }
}

# ---------------------------------------------------------------------------
# Step 5 — Smoke: unsupported -call value must exit non-zero
#
# Passes -call=delete-pod and asserts the binary exits non-zero immediately
# without making any Kubernetes API call.
# ---------------------------------------------------------------------------
resource "enos_remote_exec" "smoke_invalid_call" {
  count      = local.skip_smoke ? 0 : 1
  depends_on = [enos_local_exec.setup_rbac]

  environment = {
    KUBECLIENT_BIN = var.kubeclient_bin_path
  }

  scripts = [abspath("${path.module}/scripts/smoke-invalid-call.sh")]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.app_pod_name
      namespace         = var.app_pod_namespace
      container         = local.sidecar_container
    }
  }
}
