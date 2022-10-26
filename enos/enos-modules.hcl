module "autopilot_upgrade_storageconfig" {
  source = "./modules/autopilot_upgrade_storageconfig"
}

module "az_finder" {
  source = "./modules/az_finder"
}

module "backend_consul" {
  source = "app.terraform.io/hashicorp-qti/aws-consul/enos"

  project_name    = var.project_name
  environment     = "ci"
  common_tags     = var.tags
  ssh_aws_keypair = var.aws_ssh_keypair_name

  # Set this to a real license vault if using an Enterprise edition of Consul
  consul_license = var.backend_license_path == null ? "none" : file(abspath(var.backend_license_path))
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
  source = "app.terraform.io/hashicorp-qti/aws-infra/enos"

  project_name      = var.project_name
  environment       = "ci"
  common_tags       = var.tags
  ami_architectures = ["amd64", "arm64"]
}

module "get_local_metadata" {
  source = "./modules/get_local_metadata"
}

module "read_license" {
  source = "./modules/read_license"
}

module "vault_agent" {
  source = "./modules/vault_agent"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}


module "vault_verify_agent_output" {
  source = "./modules/vault_verify_agent_output"

  vault_instance_count = var.vault_instance_count
}

module "vault_cluster" {
  source = "app.terraform.io/hashicorp-qti/aws-vault/enos"

  common_tags       = var.tags
  environment       = "ci"
  instance_count    = var.vault_instance_count
  project_name      = var.project_name
  ssh_aws_keypair   = var.aws_ssh_keypair_name
  vault_install_dir = var.vault_install_dir
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

module "vault_verify_replication" {
  source = "./modules/vault-verify-replication"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_ui" {
  source = "./modules/vault-verify-ui"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_unsealed" {
  source = "./modules/vault_verify_unsealed"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_version" {
  source = "./modules/vault_verify_version"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}

module "vault_verify_write_test_data" {
  source = "./modules/vault-verify-write-data"

  vault_install_dir    = var.vault_install_dir
  vault_instance_count = var.vault_instance_count
}
