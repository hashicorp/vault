// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

module "autopilot_upgrade_storageconfig" {
  source = "./modules/autopilot_upgrade_storageconfig"
}

module "backend_consul" {
  source = "./modules/backend_consul"

  license   = var.backend_license_path == null ? null : file(abspath(var.backend_license_path))
  log_level = var.backend_log_level
}

module "backend_raft" {
  source = "./modules/backend_raft"
}

// Find any artifact in Artifactory. Requires the version, revision, and edition.
module "build_artifactory" {
  source = "./modules/build_artifactory_artifact"
}

// Find any released RPM or Deb in Artifactory. Requires the version, edition, distro, and distro
// version.
module "build_artifactory_package" {
  source = "./modules/build_artifactory_package"
}

// A shim "build module" suitable for use when using locally pre-built artifacts or a zip bundle
// from releases.hashicorp.com. When using a local pre-built artifact it requires the local
// artifact path. When using a release zip it does nothing as you'll need to configure the
// vault_cluster module with release info instead.
module "build_crt" {
  source = "./modules/build_crt"
}

// Build the local branch and package it into a zip artifact. Requires the goarch, goos, build tags,
// and bundle path.
module "build_local" {
  source = "./modules/build_local"
}

module "create_vpc" {
  source = "./modules/create_vpc"

  environment = "ci"
  common_tags = var.tags
}

module "ec2_info" {
  source = "./modules/ec2_info"
}

module "get_local_metadata" {
  source = "./modules/get_local_metadata"
}

module "generate_dr_operation_token" {
  source = "./modules/generate_dr_operation_token"

  vault_install_dir = var.vault_install_dir
}

module "generate_failover_secondary_token" {
  source = "./modules/generate_failover_secondary_token"

  vault_install_dir = var.vault_install_dir
}

module "generate_secondary_public_key" {
  source = "./modules/generate_secondary_public_key"

  vault_install_dir = var.vault_install_dir
}

module "generate_secondary_token" {
  source = "./modules/generate_secondary_token"

  vault_install_dir = var.vault_install_dir
}

module "install_packages" {
  source = "./modules/install_packages"
}

module "read_license" {
  source = "./modules/read_license"
}

module "replication_data" {
  source = "./modules/replication_data"
}

module "seal_awskms" {
  source = "./modules/seal_awskms"

  cluster_ssh_keypair = var.aws_ssh_keypair_name
  common_tags         = var.tags
}

module "seal_shamir" {
  source = "./modules/seal_shamir"

  cluster_ssh_keypair = var.aws_ssh_keypair_name
  common_tags         = var.tags
}

module "seal_pkcs11" {
  source = "./modules/seal_pkcs11"

  cluster_ssh_keypair = var.aws_ssh_keypair_name
  common_tags         = var.tags
}

module "shutdown_node" {
  source = "./modules/shutdown_node"
}

module "shutdown_multiple_nodes" {
  source = "./modules/shutdown_multiple_nodes"
}

module "start_vault" {
  source = "./modules/start_vault"

  install_dir = var.vault_install_dir
  log_level   = var.vault_log_level
}

module "stop_vault" {
  source = "./modules/stop_vault"
}

// create target instances using ec2:CreateFleet
module "target_ec2_fleet" {
  source = "./modules/target_ec2_fleet"

  common_tags  = var.tags
  project_name = var.project_name
  ssh_keypair  = var.aws_ssh_keypair_name
}

// create target instances using ec2:RunInstances
module "target_ec2_instances" {
  source = "./modules/target_ec2_instances"

  common_tags   = var.tags
  ports_ingress = values(global.ports)
  project_name  = var.project_name
  ssh_keypair   = var.aws_ssh_keypair_name
}

// don't create instances but satisfy the module interface
module "target_ec2_shim" {
  source = "./modules/target_ec2_shim"

  common_tags   = var.tags
  ports_ingress = values(global.ports)
  project_name  = var.project_name
  ssh_keypair   = var.aws_ssh_keypair_name
}

// create target instances using ec2:RequestSpotFleet
module "target_ec2_spot_fleet" {
  source = "./modules/target_ec2_spot_fleet"

  common_tags  = var.tags
  project_name = var.project_name
  ssh_keypair  = var.aws_ssh_keypair_name
}

module "vault_agent" {
  source = "./modules/vault_agent"

  vault_install_dir = var.vault_install_dir
  vault_agent_port  = global.ports["vault_agent"]["port"]
}

module "vault_proxy" {
  source = "./modules/vault_proxy"

  vault_install_dir = var.vault_install_dir
  vault_proxy_port  = global.ports["vault_proxy"]["port"]
}

module "vault_verify_agent_output" {
  source = "./modules/vault_verify_agent_output"
}

module "vault_cluster" {
  source = "./modules/vault_cluster"

  install_dir     = var.vault_install_dir
  consul_license  = var.backend_license_path == null ? null : file(abspath(var.backend_license_path))
  cluster_tag_key = global.vault_tag_key
  log_level       = var.vault_log_level
}

module "vault_get_cluster_ips" {
  source = "./modules/vault_get_cluster_ips"

  vault_install_dir = var.vault_install_dir
}

module "vault_failover_demote_dr_primary" {
  source = "./modules/vault_failover_demote_dr_primary"

  vault_install_dir = var.vault_install_dir
}

module "vault_failover_promote_dr_secondary" {
  source = "./modules/vault_failover_promote_dr_secondary"

  vault_install_dir = var.vault_install_dir
}

module "vault_failover_update_dr_primary" {
  source = "./modules/vault_failover_update_dr_primary"

  vault_install_dir = var.vault_install_dir
}

module "vault_raft_remove_peer" {
  source            = "./modules/vault_raft_remove_peer"
  vault_install_dir = var.vault_install_dir
}

module "vault_setup_dr_primary" {
  source = "./modules/vault_setup_dr_primary"

  vault_install_dir = var.vault_install_dir
}

module "vault_setup_perf_primary" {
  source = "./modules/vault_setup_perf_primary"

  vault_install_dir = var.vault_install_dir
}

module "vault_setup_replication_secondary" {
  source = "./modules/vault_setup_replication_secondary"

  vault_install_dir = var.vault_install_dir
}

module "vault_step_down" {
  source = "./modules/vault_step_down"

  vault_install_dir = var.vault_install_dir
}

module "vault_test_ui" {
  source = "./modules/vault_test_ui"

  ui_run_tests = var.ui_run_tests
}

module "vault_unseal_replication_followers" {
  source = "./modules/vault_unseal_replication_followers"

  vault_install_dir = var.vault_install_dir
}

module "vault_upgrade" {
  source = "./modules/vault_upgrade"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_autopilot" {
  source = "./modules/vault_verify_autopilot"

  vault_autopilot_upgrade_status = "await-server-removal"
  vault_install_dir              = var.vault_install_dir
}

module "vault_verify_dr_replication" {
  source = "./modules/vault_verify_dr_replication"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_secrets_engines_create" {
  source = "./modules/verify_secrets_engines/modules/create"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_secrets_engines_read" {
  source = "./modules/verify_secrets_engines/modules/read"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_default_lcq" {
  source = "./modules/vault_verify_default_lcq"

  vault_autopilot_default_max_leases = "300000"
}

module "vault_verify_performance_replication" {
  source = "./modules/vault_verify_performance_replication"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_raft_auto_join_voter" {
  source = "./modules/vault_verify_raft_auto_join_voter"

  vault_install_dir       = var.vault_install_dir
  vault_cluster_addr_port = global.ports["vault_cluster"]["port"]
}

module "vault_verify_replication" {
  source = "./modules/vault_verify_replication"
}

module "vault_verify_ui" {
  source = "./modules/vault_verify_ui"
}

module "vault_verify_undo_logs" {
  source = "./modules/vault_verify_undo_logs"

  vault_install_dir = var.vault_install_dir
}

module "vault_wait_for_cluster_unsealed" {
  source = "./modules/vault_wait_for_cluster_unsealed"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_version" {
  source = "./modules/vault_verify_version"

  vault_install_dir = var.vault_install_dir
}

module "vault_wait_for_leader" {
  source = "./modules/vault_wait_for_leader"

  vault_install_dir = var.vault_install_dir
}

module "vault_wait_for_seal_rewrap" {
  source = "./modules/vault_wait_for_seal_rewrap"

  vault_install_dir = var.vault_install_dir
}

module "verify_log_secrets" {
  source = "./modules/verify_log_secrets"

  radar_license_path = var.vault_radar_license_path != null ? abspath(var.vault_radar_license_path) : null
}

module "verify_seal_type" {
  source = "./modules/verify_seal_type"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_billing_start_date" {
  source = "./modules/vault_verify_billing_start_date"

  vault_install_dir       = var.vault_install_dir
  vault_instance_count    = var.vault_instance_count
  vault_cluster_addr_port = global.ports["vault_cluster"]["port"]
}
