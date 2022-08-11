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

module "create_vpc" {
  source = "app.terraform.io/hashicorp-qti/aws-infra/enos"

  project_name      = var.project_name
  environment       = "ci"
  common_tags       = var.tags
  ami_architectures = ["amd64", "arm64"]
}

module "read_license" {
  source = "./modules/read_license"
}

module "vault_cluster" {
  source  = "app.terraform.io/hashicorp-qti/aws-vault/enos"
  version = ">= 0.7.0"

  project_name    = var.project_name
  environment     = "ci"
  common_tags     = var.tags
  ssh_aws_keypair = var.aws_ssh_keypair_name
}


module "vault_upgrade" {
  source = "./modules/vault_upgrade"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_version" {
  source = "./modules/vault_verify_version"

  vault_install_dir = var.vault_install_dir
}

module "get_vault_version" {
  source = "./modules/get_vault_version"
}

module "vault_autopilot_upgrade_storageconfig" {
  source = "./modules/vault_autopilot_upgrade_storageconfig"
}

module "verify_autopilot" {
  source = "./modules/vault_verify_autopilot"
}
