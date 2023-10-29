# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    # We need to specify the provider source in each module until we publish it
    # to the public registry
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
      version = ">= 0.4.7"
    }
  }
}

data "enos_environment" "localhost" {}

locals {
  bin_path = "${var.install_dir}/vault"
  environment = local.seal_secondary == null ? var.environment : merge(
    var.environment,
    { VAULT_ENABLE_SEAL_HA_BETA : tobool(var.seal_ha_beta) },
  )
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
  seals = local.seal_secondary.type == "none" ? { primary = local.seal_primary } : {
    primary   = local.seal_primary
    secondary = local.seal_secondary
  }
  seals_primary = {
    "awskms" = {
      type = "awskms"
      attributes = {
        name       = "primary"
        kms_key_id = var.seal_key_name
      }
    }
    "shamir" = {
      type       = "shamir"
      attributes = null
    }
  }
  seal_primary = local.seals_primary[var.seal_type]
  seals_secondary = {
    "awskms" = {
      type = "awskms"
      attributes = {
        name       = "secondary"
        kms_key_id = var.seal_key_name_secondary
      }
    }
    "none" = {
      type       = "none"
      attributes = null
    }
  }
  seal_secondary = local.seals_secondary[var.seal_type_secondary]
  storage_config = [for idx, host in var.target_hosts : (var.storage_backend == "raft" ?
    merge(
      {
        node_id = "${var.storage_node_prefix}_${idx}"
      },
      var.storage_backend_attrs
    ) :
    {
      address = "127.0.0.1:8500"
      path    = "vault"
    })
  ]
}

resource "enos_vault_start" "leader" {
  for_each = local.leader

  bin_path    = local.bin_path
  config_dir  = var.config_dir
  environment = local.environment
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
    log_level = var.log_level
    storage = {
      type       = var.storage_backend
      attributes = ({ for key, value in local.storage_config[each.key] : key => value })
    }
    seals = local.seals
    ui    = true
  }
  license        = var.license
  manage_service = var.manage_service
  username       = var.service_username
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
  environment = local.environment
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
    log_level = var.log_level
    storage = {
      type       = var.storage_backend
      attributes = { for key, value in local.storage_config[each.key] : key => value }
    }
    seals = local.seals
    ui    = true
  }
  license        = var.license
  manage_service = var.manage_service
  username       = var.service_username
  unit_name      = "vault"

  transport = {
    ssh = {
      host = var.target_hosts[each.value].public_ip
    }
  }
}
