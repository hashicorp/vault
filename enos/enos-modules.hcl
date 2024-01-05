# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

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

module "build_crt" {
  source = "./modules/build_crt"
}

module "build_local" {
  source = "./modules/build_local"
}

module "build_artifactory" {
  source = "./modules/vault_artifactory_artifact"
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

# create target instances using ec2:CreateFleet
module "target_ec2_fleet" {
  source = "./modules/target_ec2_fleet"

  common_tags  = var.tags
  project_name = var.project_name
  ssh_keypair  = var.aws_ssh_keypair_name
}

# create target instances using ec2:RunInstances
module "target_ec2_instances" {
  source = "./modules/target_ec2_instances"

  common_tags  = var.tags
  project_name = var.project_name
  ssh_keypair  = var.aws_ssh_keypair_name
}

# don't create instances but satisfy the module interface
module "target_ec2_shim" {
  source = "./modules/target_ec2_shim"

  common_tags  = var.tags
  project_name = var.project_name
  ssh_keypair  = var.aws_ssh_keypair_name
}

# create target instances using ec2:RequestSpotFleet
module "target_ec2_spot_fleet" {
  source = "./modules/target_ec2_spot_fleet"

  common_tags  = var.tags
  project_name = var.project_name
  ssh_keypair  = var.aws_ssh_keypair_name
}

module "vault_agent" {
  source = "./modules/vault_agent"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_proxy" {
  source = "./modules/vault_proxy"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_agent_output" {
  source = "./modules/vault_verify_agent_output"

  vault_instance_count = var.vault_instance_count
}

module "vault_cluster" {
  source = "./modules/vault_cluster"

  install_dir    = var.vault_install_dir
  consul_license = var.backend_license_path == null ? null : file(abspath(var.backend_license_path))
  log_level      = var.vault_log_level
}

module "vault_get_cluster_ips" {
  source = "./modules/vault_get_cluster_ips"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_raft_remove_peer" {
  source            = "./modules/vault_raft_remove_peer"
  vault_install_dir = var.vault_install_dir
}

module "vault_setup_perf_secondary" {
  source = "./modules/vault_setup_perf_secondary"

  vault_install_dir = var.vault_install_dir
}

module "vault_test_ui" {
  source = "./modules/vault_test_ui"

  ui_run_tests = var.ui_run_tests
}

module "vault_unseal_nodes" {
  source = "./modules/vault_unseal_nodes"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_upgrade" {
  source = "./modules/vault_upgrade"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}


module "vault_verify_autopilot" {
  source = "./modules/vault_verify_autopilot"

  vault_autopilot_upgrade_status = "await-server-removal"
  vault_install_dir              = var.vault_install_dir
  vault_instance_count           = var.vault_instance_count
}

module "vault_verify_raft_auto_join_voter" {
  source = "./modules/vault_verify_raft_auto_join_voter"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_undo_logs" {
  source = "./modules/vault_verify_undo_logs"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_default_lcq" {
  source = "./modules/vault_verify_default_lcq"

  vault_autopilot_default_max_leases = "300000"
  vault_instance_count               = var.vault_instance_count
}

module "vault_verify_replication" {
  source = "./modules/vault_verify_replication"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_ui" {
  source = "./modules/vault_verify_ui"

  vault_instance_count = var.vault_instance_count
}

module "vault_verify_unsealed" {
  source = "./modules/vault_verify_unsealed"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_setup_perf_primary" {
  source = "./modules/vault_setup_perf_primary"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_read_data" {
  source = "./modules/vault_verify_read_data"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_performance_replication" {
  source = "./modules/vault_verify_performance_replication"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_version" {
  source = "./modules/vault_verify_version"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_write_data" {
  source = "./modules/vault_verify_write_data"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_wait_for_leader" {
  source = "./modules/vault_wait_for_leader"

  vault_install_dir = var.vault_install_dir
}

module "vault_wait_for_seal_rewrap" {
  source = "./modules/vault_wait_for_seal_rewrap"

  vault_install_dir = var.vault_install_dir
}

module "verify_seal_type" {
  source = "./modules/verify_seal_type"

  vault_install_dir = var.vault_install_dir
}
