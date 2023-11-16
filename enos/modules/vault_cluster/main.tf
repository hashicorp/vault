# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    # We need to specify the provider source in each module until we publish it
    # to the public registry
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
      version = ">= 0.4.0"
    }
  }
}

data "enos_environment" "localhost" {}

locals {
  audit_device_file_path = "/var/log/vault/vault_audit.log"
  bin_path               = "${var.install_dir}/vault"
  consul_bin_path        = "${var.consul_install_dir}/consul"
  enable_audit_devices   = var.enable_audit_devices && var.initialize_cluster
  // In order to get Terraform to plan we have to use collections with keys
  // that are known at plan time. In order for our module to work our var.target_hosts
  // must be a map with known keys at plan time. Here we're creating locals
  // that keep track of index values that point to our target hosts.
  followers = toset(slice(local.instances, 1, length(local.instances)))
  instances = [for idx in range(length(var.target_hosts)) : tostring(idx)]
  key_shares = {
    "awskms" = null
    "shamir" = 5
  }
  key_threshold = {
    "awskms" = null
    "shamir" = 3
  }
  leader = toset(slice(local.instances, 0, 1))
  recovery_shares = {
    "awskms" = 5
    "shamir" = null
  }
  recovery_threshold = {
    "awskms" = 3
    "shamir" = null
  }
  vault_service_user = "vault"
}

resource "enos_bundle_install" "consul" {
  for_each = {
    for idx, host in var.target_hosts : idx => var.target_hosts[idx]
    if var.storage_backend == "consul"
  }

  destination = var.consul_install_dir
  release     = merge(var.consul_release, { product = "consul" })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_bundle_install" "vault" {
  for_each = var.target_hosts

  destination = var.install_dir
  release     = var.release == null ? var.release : merge({ product = "vault" }, var.release)
  artifactory = var.artifactory_release
  path        = var.local_artifact_path

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "install_packages" {
  depends_on = [
    enos_bundle_install.vault, // Don't race for the package manager locks with vault install
  ]
  for_each = {
    for idx, host in var.target_hosts : idx => var.target_hosts[idx]
    if length(var.packages) > 0
  }

  environment = {
    PACKAGES = join(" ", var.packages)
  }

  scripts = [abspath("${path.module}/scripts/install-packages.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_consul_start" "consul" {
  for_each = enos_bundle_install.consul

  bin_path = local.consul_bin_path
  data_dir = var.consul_data_dir
  config = {
    data_dir         = var.consul_data_dir
    datacenter       = "dc1"
    retry_join       = ["provider=aws tag_key=${var.backend_cluster_tag_key} tag_value=${var.backend_cluster_name}"]
    server           = false
    bootstrap_expect = 0
    license          = var.consul_license
    log_level        = var.consul_log_level
    log_file         = var.consul_log_file
  }
  license   = var.consul_license
  unit_name = "consul"
  username  = "consul"

  transport = {
    ssh = {
      host = var.target_hosts[each.key].public_ip
    }
  }
}

module "start_vault" {
  source = "../start_vault"

  depends_on = [
    enos_consul_start.consul,
    enos_bundle_install.vault,
  ]

  cluster_name            = var.cluster_name
  config_dir              = var.config_dir
  install_dir             = var.install_dir
  license                 = var.license
  log_level               = var.log_level
  manage_service          = var.manage_service
  seal_ha_beta            = var.seal_ha_beta
  seal_key_name           = var.seal_key_name
  seal_key_name_secondary = var.seal_key_name_secondary
  seal_type               = var.seal_type
  seal_type_secondary     = var.seal_type_secondary
  service_username        = local.vault_service_user
  storage_backend         = var.storage_backend
  storage_backend_attrs   = var.storage_backend_addl_config
  storage_node_prefix     = var.storage_node_prefix
  target_hosts            = var.target_hosts
}

resource "enos_vault_init" "leader" {
  depends_on = [
    module.start_vault,
  ]
  for_each = toset([
    for idx, leader in local.leader : leader
    if var.initialize_cluster
  ])

  bin_path   = local.bin_path
  vault_addr = module.start_vault.leader[0].config.api_addr

  key_shares    = local.key_shares[var.seal_type]
  key_threshold = local.key_threshold[var.seal_type]

  recovery_shares    = local.recovery_shares[var.seal_type]
  recovery_threshold = local.recovery_threshold[var.seal_type]

  transport = {
    ssh = {
      host = var.target_hosts[each.value].public_ip
    }
  }
}

resource "enos_vault_unseal" "leader" {
  depends_on = [
    module.start_vault,
    enos_vault_init.leader,
  ]
  for_each = enos_vault_init.leader // only unseal the leader if we initialized it

  bin_path    = local.bin_path
  vault_addr  = module.start_vault.leader[each.key].config.api_addr
  seal_type   = var.seal_type
  unseal_keys = var.seal_type != "shamir" ? null : coalesce(var.shamir_unseal_keys, enos_vault_init.leader[0].unseal_keys_hex)

  transport = {
    ssh = {
      host = var.target_hosts[tolist(local.leader)[0]].public_ip
    }
  }
}

resource "enos_vault_unseal" "followers" {
  depends_on = [
    enos_vault_init.leader,
    enos_vault_unseal.leader,
  ]
  // Only unseal followers if we're not using an auto-unseal method and we've
  // initialized the cluster
  for_each = toset([
    for idx, follower in local.followers : follower
    if var.seal_type == "shamir" && var.initialize_cluster
  ])

  bin_path    = local.bin_path
  vault_addr  = module.start_vault.followers[each.key].config.api_addr
  seal_type   = var.seal_type
  unseal_keys = var.seal_type != "shamir" ? null : coalesce(var.shamir_unseal_keys, enos_vault_init.leader[0].unseal_keys_hex)

  transport = {
    ssh = {
      host = var.target_hosts[each.value].public_ip
    }
  }
}

// Force unseal the cluster. This is used if the vault-cluster module is used
// to add additional nodes to a cluster via auto-pilot, or some other means.
// When that happens we'll want to set initialize_cluster to false and
// force_unseal to true.
resource "enos_vault_unseal" "maybe_force_unseal" {
  depends_on = [
    module.start_vault.followers,
  ]
  for_each = {
    for idx, host in var.target_hosts : idx => host
    if var.force_unseal && !var.initialize_cluster
  }

  bin_path   = local.bin_path
  vault_addr = "http://localhost:8200"
  seal_type  = var.seal_type
  unseal_keys = coalesce(
    var.shamir_unseal_keys,
    try(enos_vault_init.leader[0].unseal_keys_hex, null),
  )

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# We need to ensure that the directory used for audit logs is present and accessible to the vault
# user on all nodes, since logging will only happen on the leader.
resource "enos_remote_exec" "create_audit_log_dir" {
  depends_on = [
    module.start_vault,
    enos_vault_unseal.leader,
    enos_vault_unseal.followers,
    enos_vault_unseal.maybe_force_unseal,
  ]
  for_each = toset([
    for idx, host in toset(local.instances) : idx
    if var.enable_audit_devices
  ])

  environment = {
    LOG_FILE_PATH = local.audit_device_file_path
    SERVICE_USER  = local.vault_service_user
  }

  scripts = [abspath("${path.module}/scripts/create_audit_log_dir.sh")]

  transport = {
    ssh = {
      host = var.target_hosts[each.value].public_ip
    }
  }
}

resource "enos_remote_exec" "enable_audit_devices" {
  depends_on = [
    enos_remote_exec.create_audit_log_dir,
  ]
  for_each = toset([
    for idx in local.leader : idx
    if local.enable_audit_devices
  ])

  environment = {
    VAULT_TOKEN    = enos_vault_init.leader[each.key].root_token
    VAULT_ADDR     = "http://127.0.0.1:8200"
    VAULT_BIN_PATH = local.bin_path
    LOG_FILE_PATH  = local.audit_device_file_path
    SERVICE_USER   = local.vault_service_user
  }

  scripts = [abspath("${path.module}/scripts/enable_audit_logging.sh")]

  transport = {
    ssh = {
      host = var.target_hosts[each.key].public_ip
    }
  }
}

resource "enos_local_exec" "wait_for_install_packages" {
  depends_on = [
    enos_remote_exec.install_packages,
  ]

  inline = ["true"]
}
