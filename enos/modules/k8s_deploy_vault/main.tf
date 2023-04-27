# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_version = ">= 1.0"

  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "2.6.0"
    }
  }
}

locals {
  helm_chart_settings = {
    "server.ha.enabled"             = "true"
    "server.ha.replicas"            = var.vault_instance_count
    "server.ha.raft.enabled"        = "true"
    "server.affinity"               = ""
    "server.image.repository"       = var.image_repository
    "server.image.tag"              = var.image_tag
    "server.image.pullPolicy"       = "Never" # Forces local image use
    "server.resources.requests.cpu" = "50m"
    "server.limits.memory"          = "200m"
    "server.limits.cpu"             = "200m"
    "server.ha.raft.config"         = file("${abspath(path.module)}/raft-config.hcl")
    "server.dataStorage.size"       = "100m"
    "server.logLevel"               = var.vault_log_level
  }
  all_helm_chart_settings = var.ent_license == null ? local.helm_chart_settings : merge(local.helm_chart_settings, {
    "server.extraEnvironmentVars.VAULT_LICENSE" = var.ent_license
  })

  vault_address = "http://127.0.0.1:8200"

  instance_indexes = [for idx in range(var.vault_instance_count) : tostring(idx)]

  leader_idx    = local.instance_indexes[0]
  followers_idx = toset(slice(local.instance_indexes, 1, var.vault_instance_count))
}

resource "helm_release" "vault" {
  name = "vault"

  repository = "https://helm.releases.hashicorp.com"
  chart      = "vault"

  dynamic "set" {
    for_each = local.all_helm_chart_settings

    content {
      name  = set.key
      value = set.value
    }
  }
}

data "enos_kubernetes_pods" "vault_pods" {
  kubeconfig_base64 = var.kubeconfig_base64
  context_name      = var.context_name
  namespace         = helm_release.vault.namespace
  label_selectors = [
    "app.kubernetes.io/name=vault",
    "component=server"
  ]

  depends_on = [helm_release.vault]
}

resource "enos_vault_init" "leader" {
  bin_path   = "/bin/vault"
  vault_addr = local.vault_address

  key_shares    = 5
  key_threshold = 3

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = data.enos_kubernetes_pods.vault_pods.pods[local.leader_idx].name
      namespace         = data.enos_kubernetes_pods.vault_pods.pods[local.leader_idx].namespace
    }
  }
}

resource "enos_vault_unseal" "leader" {
  bin_path    = "/bin/vault"
  vault_addr  = local.vault_address
  seal_type   = "shamir"
  unseal_keys = enos_vault_init.leader.unseal_keys_b64

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = data.enos_kubernetes_pods.vault_pods.pods[local.leader_idx].name
      namespace         = data.enos_kubernetes_pods.vault_pods.pods[local.leader_idx].namespace
    }
  }

  depends_on = [enos_vault_init.leader]
}

// We need to manually join the followers since the join request must only happen after the leader
// has been initialized. We could use retry join, but in that case we'd need to restart the follower
// pods once the leader is setup. The default helm deployment configuration for an HA cluster as
// documented here: https://learn.hashicorp.com/tutorials/vault/kubernetes-raft-deployment-guide#configure-vault-helm-chart
// uses a liveness probe that automatically restarts nodes that are not healthy. This works well for
// clusters that are configured with auto-unseal as eventually the nodes would join and unseal.
resource "enos_remote_exec" "raft_join" {
  for_each = local.followers_idx

  inline = [
    // asserts that vault is ready
    "for i in 1 2 3 4 5; do vault status > /dev/null 2>&1 && break || sleep 5; done",
    // joins the follower to the leader
    "vault operator raft join http://vault-0.vault-internal:8200"
  ]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = data.enos_kubernetes_pods.vault_pods.pods[each.key].name
      namespace         = data.enos_kubernetes_pods.vault_pods.pods[each.key].namespace
    }
  }

  depends_on = [enos_vault_unseal.leader]
}


resource "enos_vault_unseal" "followers" {
  for_each = local.followers_idx

  bin_path    = "/bin/vault"
  vault_addr  = local.vault_address
  seal_type   = "shamir"
  unseal_keys = enos_vault_init.leader.unseal_keys_b64

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = data.enos_kubernetes_pods.vault_pods.pods[each.key].name
      namespace         = data.enos_kubernetes_pods.vault_pods.pods[each.key].namespace
    }
  }

  depends_on = [enos_remote_exec.raft_join]
}

output "vault_root_token" {
  value = enos_vault_init.leader.root_token
}

output "vault_pods" {
  value = data.enos_kubernetes_pods.vault_pods.pods
}
