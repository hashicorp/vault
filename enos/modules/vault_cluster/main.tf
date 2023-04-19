terraform {
  required_providers {
    # We need to specify the provider source in each module until we publish it
    # to the public registry
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
      version = ">= 0.3.2"
    }
  }
}

data "enos_environment" "localhost" {}

locals {
  bin_path        = "${var.install_dir}/vault"
  consul_bin_path = "${var.consul_install_dir}/consul"
  key_shares = {
    "awskms" = null
    "shamir" = 5
  }
  key_threshold = {
    "awskms" = null
    "shamir" = 3
  }
  // In order to get Terraform to plan we have to use collections with keys
  // that are known at plan time. In order for our module to work our var.target_hosts
  // must be a map with known keys at plan time. Here we're creating locals
  // that keep track of index values that point to our target hosts.
  followers = toset(slice(local.instances, 1, length(local.instances)))
  instances = [for idx in range(length(var.target_hosts)) : tostring(idx)]
  leader    = toset(slice(local.instances, 0, 1))
  recovery_shares = {
    "awskms" = 5
    "shamir" = null
  }
  recovery_threshold = {
    "awskms" = 3
    "shamir" = null
  }
  seal = {
    "awskms" = {
      type = "awskms"
      attributes = {
        kms_key_id = var.awskms_unseal_key_arn
      }
    }
    "shamir" = {
      type       = "shamir"
      attributes = null
    }
  }
  storage_config = [for idx, host in var.target_hosts : (var.storage_backend == "raft" ?
    merge(
      {
        node_id = "${var.storage_node_prefix}_${idx}"
      },
      var.storage_backend_addl_config
    ) :
    {
      address = "127.0.0.1:8500"
      path    = "vault"
    })
  ]
}

resource "enos_remote_exec" "install_packages" {
  for_each = {
    for idx, host in var.target_hosts : idx => var.target_hosts[idx]
    if length(var.packages) > 0
  }

  content = templatefile("${path.module}/templates/install-packages.sh", {
    packages = join(" ", var.packages)
  })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
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

resource "enos_consul_start" "consul" {
  for_each = enos_bundle_install.consul

  bin_path = local.consul_bin_path
  data_dir = var.consul_data_dir
  config = {
    data_dir         = var.consul_data_dir
    datacenter       = "dc1"
    retry_join       = ["provider=aws tag_key=Type tag_value=${var.consul_cluster_tag}"]
    server           = false
    bootstrap_expect = 0
    log_level        = "INFO"
    log_file         = var.consul_log_file
  }
  unit_name = "consul"
  username  = "consul"

  transport = {
    ssh = {
      host = var.target_hosts[each.key].public_ip
    }
  }
}

resource "enos_vault_start" "leader" {
  depends_on = [
    enos_consul_start.consul,
    enos_bundle_install.vault,
  ]
  for_each = local.leader

  bin_path    = local.bin_path
  config_dir  = var.config_dir
  environment = var.config_env_vars
  config = {
    api_addr     = "http://${var.target_hosts[each.value].private_ip}:8200"
    cluster_addr = "http://${var.target_hosts[each.value].private_ip}:8201"
    cluster_name = var.cluster_name
    listener = {
      type = "tcp"
      attributes = {
        address     = "0.0.0.0:8200"
        tls_disable = "true"
      }
    }
    storage = {
      type       = var.storage_backend
      attributes = ({ for key, value in local.storage_config[each.key] : key => value })
    }
    seal = local.seal[var.unseal_method]
    ui   = true
  }
  license        = var.license
  manage_service = var.manage_service
  username       = "vault"
  unit_name      = "vault"

  transport = {
    ssh = {
      host = var.target_hosts[each.value].public_ip
    }
  }
}

resource "enos_vault_start" "followers" {
  depends_on = [
    enos_vault_start.leader,
  ]
  for_each = local.followers

  bin_path    = local.bin_path
  config_dir  = var.config_dir
  environment = var.config_env_vars
  config = {
    api_addr     = "http://${var.target_hosts[each.value].private_ip}:8200"
    cluster_addr = "http://${var.target_hosts[each.value].private_ip}:8201"
    cluster_name = var.cluster_name
    listener = {
      type = "tcp"
      attributes = {
        address     = "0.0.0.0:8200"
        tls_disable = "true"
      }
    }
    storage = {
      type       = var.storage_backend
      attributes = { for key, value in local.storage_config[each.key] : key => value }
    }
    seal = local.seal[var.unseal_method]
    ui   = true
  }
  license        = var.license
  manage_service = var.manage_service
  username       = "vault"
  unit_name      = "vault"

  transport = {
    ssh = {
      host = var.target_hosts[each.value].public_ip
    }
  }
}

resource "enos_vault_init" "leader" {
  depends_on = [
    enos_vault_start.followers,
  ]
  for_each = toset([
    for idx, leader in local.leader : leader
    if var.initialize_cluster
  ])

  bin_path   = local.bin_path
  vault_addr = enos_vault_start.leader[0].config.api_addr

  key_shares    = local.key_shares[var.unseal_method]
  key_threshold = local.key_threshold[var.unseal_method]

  recovery_shares    = local.recovery_shares[var.unseal_method]
  recovery_threshold = local.recovery_threshold[var.unseal_method]

  transport = {
    ssh = {
      host = var.target_hosts[each.value].public_ip
    }
  }
}

resource "enos_vault_unseal" "leader" {
  depends_on = [
    enos_vault_start.followers,
    enos_vault_init.leader,
  ]
  for_each = enos_vault_init.leader // only unseal the leader if we initialized it

  bin_path    = local.bin_path
  vault_addr  = enos_vault_start.leader[each.key].config.api_addr
  seal_type   = var.unseal_method
  unseal_keys = var.unseal_method != "shamir" ? null : coalesce(var.shamir_unseal_keys, enos_vault_init.leader[0].unseal_keys_hex)

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
    if var.unseal_method == "shamir" && var.initialize_cluster
  ])

  bin_path    = local.bin_path
  vault_addr  = enos_vault_start.followers[each.key].config.api_addr
  seal_type   = var.unseal_method
  unseal_keys = var.unseal_method != "shamir" ? null : coalesce(var.shamir_unseal_keys, enos_vault_init.leader[0].unseal_keys_hex)

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
    enos_vault_start.followers,
  ]
  for_each = {
    for idx, host in var.target_hosts : idx => host
    if var.force_unseal && !var.initialize_cluster
  }

  bin_path   = local.bin_path
  vault_addr = "http://localhost:8200"
  seal_type  = var.unseal_method
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

resource "enos_remote_exec" "vault_write_license" {
  for_each = toset([
    for idx, leader in local.leader : leader
    if var.initialize_cluster
  ])

  depends_on = [
    enos_vault_unseal.leader,
    enos_vault_unseal.maybe_force_unseal,
  ]

  content = templatefile("${path.module}/templates/vault-write-license.sh", {
    bin_path   = local.bin_path,
    root_token = coalesce(var.root_token, try(enos_vault_init.leader[0].root_token, null), "none")
    license    = coalesce(var.license, "none")
  })

  transport = {
    ssh = {
      host = var.target_hosts[each.value].public_ip
    }
  }
}
