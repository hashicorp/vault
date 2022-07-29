module "az_finder" {
  source = "./modules/az_finder"
}

module "backend_consul" {
  source = "app.terraform.io/hashicorp-qti/aws-consul/enos"

  project_name    = "qti-enos-provider"
  environment     = "ci"
  common_tags     = var.tags
  ssh_aws_keypair = "enos-ci-ssh-keypair"

  # Set this to a real license vault if using an Enterprise edition of Consul
  consul_license = "none"
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

  project_name      = "qti-enos-provider"
  environment       = "ci"
  common_tags       = var.tags
  ami_architectures = ["amd64", "arm64"]
}

module "read_license" {
  source = "./modules/read_license"
}

module "vault_cluster" {
  source = "app.terraform.io/hashicorp-qti/aws-vault/enos"

  project_name    = "vault-enos-integration"
  environment     = "ci"
  common_tags     = var.tags
  ssh_aws_keypair = "enos-ci-ssh-keypair"
}


module "vault_upgrade" {
  source = "./modules/vault_upgrade"

  vault_install_dir = var.vault_install_dir
}

module "vault_verify_version" {
  source = "./modules/vault_verify_version"

  vault_install_dir = var.vault_install_dir
}
