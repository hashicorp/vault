# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    # We need to specify the provider source in each module until we publish it
    # to the public registry
    enos = {
      source  = "registry.terraform.io/hashicorp-forge/enos"
      version = ">= 0.4.0"
    }
  }
}

data "enos_environment" "localhost" {}

locals {
  audit_device_file_path = "/var/log/vault/vault_audit.log"
  audit_socket_port      = "9090"
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
    "pkcs11" = null
  }
  key_threshold = {
    "awskms" = null
    "shamir" = 3
    "pkcs11" = null
  }
  leader = toset(slice(local.instances, 0, 1))
  netcat_command = {
    amzn          = "nc"
    opensuse-leap = "netcat"
    rhel          = "nc"
    sles          = "nc"
    ubuntu        = "netcat"
  }
  recovery_shares = {
    "awskms" = 5
    "shamir" = null
    "pkcs11" = 5
  }
  recovery_threshold = {
    "awskms" = 3
    "shamir" = null
    "pkcs11" = 3
  }
  vault_service_user = "vault"
}

resource "enos_host_info" "hosts" {
  for_each = var.target_hosts

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

# We run install_packages before we install Vault because for some combinations of
# certain Linux distros and artifact types (e.g. SLES and RPM packages), there may
# be packages that are required to perform Vault installation (e.g. openssl).
module "install_packages" {
  source = "../install_packages"

  hosts    = var.target_hosts
  packages = var.packages
}

resource "enos_bundle_install" "vault" {
  for_each = var.target_hosts
  depends_on = [
    module.install_packages, // Don't race for the package manager locks with install_packages
  ]

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
    # GetPrivateInterfaces is a go-sockaddr template that helps Consul get the correct
    # addr in all of our default cases. This is required in the case of Amazon Linux,
    # because amzn2 has a default docker listener that will make Consul try to use the
    # incorrect addr.
    bind_addr        = "{{ GetPrivateInterfaces | include \"type\" \"IP\" | sort \"default\" |  limit 1 | attr \"address\"}}"
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
    module.install_packages,
    enos_bundle_install.vault,
  ]

  cluster_name              = var.cluster_name
  config_dir                = var.config_dir
  config_mode               = var.config_mode
  install_dir               = var.install_dir
  license                   = var.license
  log_level                 = var.log_level
  manage_service            = var.manage_service
  seal_attributes           = var.seal_attributes
  seal_attributes_secondary = var.seal_attributes_secondary
  seal_type                 = var.seal_type
  seal_type_secondary       = var.seal_type_secondary
  service_username          = local.vault_service_user
  storage_backend           = var.storage_backend
  storage_backend_attrs     = var.storage_backend_addl_config
  storage_node_prefix       = var.storage_node_prefix
  target_hosts              = var.target_hosts
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

# Add the vault install location to the PATH and set up VAULT_ADDR and VAULT_TOKEN environement
# variables in the login shell so we don't have to do it if/when we login in to a cluster node.
resource "enos_remote_exec" "configure_login_shell_profile" {
  depends_on = [
    enos_vault_init.leader,
    enos_vault_unseal.leader,
  ]
  for_each = var.target_hosts

  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_TOKEN       = var.root_token != null ? var.root_token : try(enos_vault_init.leader[0].root_token, "_")
    VAULT_INSTALL_DIR = var.install_dir
  }

  scripts = [abspath("${path.module}/scripts/set-up-login-shell-profile.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# Add a motd to assist people that might be logging in.
resource "enos_file" "motd" {
  depends_on = [
    enos_remote_exec.configure_login_shell_profile
  ]
  for_each = var.target_hosts

  destination = "/etc/motd"
  content     = <<EOF
 ▄█    █▄     ▄████████ ███    █▄   ▄█           ███
███    ███   ███    ███ ███    ███ ███       ▀█████████▄
███    ███   ███    ███ ███    ███ ███          ▀███▀▀██
███    ███   ███    ███ ███    ███ ███           ███   ▀
███    ███ ▀███████████ ███    ███ ███           ███
███    ███   ███    ███ ███    ███ ███           ███
███    ███   ███    ███ ███    ███ ███▌    ▄     ███
 ▀██████▀    ███    █▀  ████████▀  █████▄▄██    ▄████▀
                                   ▀
We've added `vault` to the PATH for you and configured
the VAULT_ADDR and VAULT_TOKEN with the root token.

Have fun!
EOF

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

  scripts = [abspath("${path.module}/scripts/create-audit-log-dir.sh")]

  transport = {
    ssh = {
      host = var.target_hosts[each.value].public_ip
    }
  }
}

# We need to ensure that the socket listener used for the audit socket device is listening on each
# node in the cluster. If we have a leader election or vault is restarted it'll fail unless the
# listener is running.
resource "enos_remote_exec" "start_audit_socket_listener" {
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
    NETCAT_COMMAND = local.netcat_command[enos_host_info.hosts[each.key].distro]
    SOCKET_PORT    = local.audit_socket_port
  }

  scripts = [abspath("${path.module}/scripts/start-audit-socket-listener.sh")]

  transport = {
    ssh = {
      host = var.target_hosts[each.value].public_ip
    }
  }
}

resource "enos_remote_exec" "enable_audit_devices" {
  depends_on = [
    enos_remote_exec.create_audit_log_dir,
    enos_remote_exec.start_audit_socket_listener,
  ]
  for_each = toset([
    for idx in local.leader : idx
    if local.enable_audit_devices
  ])

  environment = {
    LOG_FILE_PATH  = local.audit_device_file_path
    SOCKET_PORT    = local.audit_socket_port
    VAULT_ADDR     = "http://127.0.0.1:8200"
    VAULT_BIN_PATH = local.bin_path
    VAULT_TOKEN    = enos_vault_init.leader[each.key].root_token
  }

  scripts = [abspath("${path.module}/scripts/enable-audit-devices.sh")]

  transport = {
    ssh = {
      host = var.target_hosts[each.key].public_ip
    }
  }
}
