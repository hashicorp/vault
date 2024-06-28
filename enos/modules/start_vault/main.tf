# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    # We need to specify the provider source in each module until we publish it
    # to the public registry
    enos = {
      source  = "registry.terraform.io/hashicorp-forge/enos"
      version = ">= 0.4.10"
    }
  }
}

locals {
  bin_path = "${var.install_dir}/vault"
  // In order to get Terraform to plan we have to use collections with keys that are known at plan
  // time. Here we're creating locals that keep track of index values that point to our target hosts.
  followers = toset(slice(local.instances, 1, length(local.instances)))
  instances = [for idx in range(length(var.target_hosts)) : tostring(idx)]
  leader    = toset(slice(local.instances, 0, 1))
  // Handle cases where we might have to distribute HSM tokens for the pkcs11 seal before starting
  // vault.
  token_base64           = try(lookup(var.seal_attributes, "token_base64", ""), "")
  token_base64_secondary = try(lookup(var.seal_attributes_secondary, "token_base64", ""), "")
  // This module currently supports up to two defined seals. Most of our locals logic here is for
  // creating the correct seal configuration.
  seals = {
    primary   = local.seal_primary
    secondary = local.seal_secondary
  }
  seals_primary = {
    awskms = {
      type = "awskms"
      attributes = merge(
        {
          name     = var.seal_alias
          priority = var.seal_priority
        }, var.seal_attributes
      )
    }
    pkcs11 = {
      type = "pkcs11"
      attributes = merge(
        {
          name     = var.seal_alias
          priority = var.seal_priority
        },
        // Strip out attributes that aren't supposed to be in seal stanza like our base64 encoded
        // softhsm blob and the token directory. We'll also inject the shared object library
        // location that we detect on the target machines. This allows use to create the token and
        // keys on a machines that have different shared object locations.
        merge(
          try({ for key, val in var.seal_attributes : key => val if key != "token_base64" && key != "token_dir" }, {}),
          # Note: the below reference has to point to a specific instance of the maybe_configure_hsm
          # module (in this case [0]) due to the maybe_configure_hsm module call using `count` to control whether it runs or not.
          try({ lib = module.maybe_configure_hsm[0].lib }, {})
        ),
      )
    }
    shamir = {
      type       = "shamir"
      attributes = null
    }
  }
  seal_primary = local.seals_primary[var.seal_type]
  seals_secondary = {
    awskms = {
      type = "awskms"
      attributes = merge(
        {
          name     = var.seal_alias_secondary
          priority = var.seal_priority_secondary
        }, var.seal_attributes_secondary
      )
    }
    pkcs11 = {
      type = "pkcs11"
      attributes = merge(
        {
          name     = var.seal_alias_secondary
          priority = var.seal_priority_secondary
        },
        merge(
          try({ for key, val in var.seal_attributes_secondary : key => val if key != "token_base64" && key != "token_dir" }, {}),
          # Note: the below reference has to point to a specific instance of the maybe_configure_hsm_secondary
          # module (in this case [0]) due to the maybe_configure_hsm_secondary module call using `count` to control whether it runs or not.
          try({ lib = module.maybe_configure_hsm_secondary[0].lib }, {})
        ),
      )
    }
    none = {
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

# You might be wondering why our start_vault module, which supports shamir, awskms, and pkcs11 seal
# types, contains sub-modules that are only used for HSM. Well, each of those seal devices has
# different requirements and as such we have some seal specific requirements before starting Vault.
#
# A Shamir seal key cannot exist until Vault has already started, so this modules responsibility for
# shamir seals is ensuring that the seal type is passed to the enos_vault_start resource. That's it.
#
# Auto-unseal with a KMS requires that we configure the enos_vault_start resource with the correct
# seal type and the attributes necessary to know which KMS key to use. Vault should automatically
# unseal if we've given it the correct configuration. As long as Vault is able to access the key
# in the KMS it should be able to start. That's normally done via roles associated to the target
# machines, which is outside the scope of this module.
#
# Auto-unseal with an HSM and PKCS#11 is more complicated because a shared object library, which is
# how we interface with the HSM, must be present on each node in order to start Vault. In the real
# world this means an actual HSM in the same rack or data center as every node in the Vault cluster,
# but in our case we're creating ephemeral infrastructure for these test scenarios and don't have a
# real HSM available. We could use CloudHSM or the like, but at the time of writing CloudHSM
# provisioning takes anywhere from 30 to 60 minutes and costs upwards of $2 dollars an hour. That's
# far too long and expensive for scenarios we'll run fairly frequently. Instead, we test using a
# software HSM. Using a software HSM solves the cost and speed problems but creates new set of
# problems. We need to ensure every node in the cluster has access to the same "HSM" and with
# softhsm that means the same software, configuration, tokens and keys. Our `seal_pkcs11` module
# takes care of creating the token and keys, but that's the end of the road for that module. It's
# our job to ensure that when we're starting Vault with a software HSM that we'll ensure the correct
# software, configuration and data are available on the nodes. That's where the following two
# modules come in. They handle installing the required software, configuring it, and distributing
# the key data that was passed in via seal attributes.
module "maybe_configure_hsm" {
  source = "../softhsm_distribute_vault_keys"
  count  = (var.seal_type == "pkcs11" || var.seal_type_secondary == "pkcs11") ? 1 : 0

  hosts        = var.target_hosts
  token_base64 = local.token_base64
}

module "maybe_configure_hsm_secondary" {
  source     = "../softhsm_distribute_vault_keys"
  depends_on = [module.maybe_configure_hsm]
  count      = (var.seal_type == "pkcs11" || var.seal_type_secondary == "pkcs11") ? 1 : 0

  hosts        = var.target_hosts
  token_base64 = local.token_base64_secondary
}

resource "enos_vault_start" "leader" {
  for_each = local.leader
  depends_on = [
    module.maybe_configure_hsm_secondary,
  ]

  bin_path    = local.bin_path
  config_dir  = var.config_dir
  config_mode = var.config_mode
  environment = var.environment
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
  config_mode = var.config_mode
  environment = var.environment
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

output "token_base64" {
  value = local.token_base64
}

output "token_base64_secondary" {
  value = local.token_base64_secondary
}
