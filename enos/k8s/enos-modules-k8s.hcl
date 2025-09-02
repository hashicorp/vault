# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

module "create_kind_cluster" {
  source = "../modules/local_kind_cluster"
}

module "load_docker_image" {
  source = "../modules/load_docker_image"
}

module "k8s_deploy_vault" {
  source = "../modules/k8s_deploy_vault"

  vault_instance_count = var.instance_count
}

module "k8s_verify_build_date" {
  source = "../modules/k8s_vault_verify_build_date"

  vault_instance_count = var.instance_count
}

module "k8s_verify_replication" {
  source = "../modules/k8s_vault_verify_replication"

  vault_instance_count = var.instance_count
}

module "k8s_verify_ui" {
  source = "../modules/k8s_vault_verify_ui"

  vault_instance_count = var.instance_count
}

module "k8s_verify_version" {
  source = "../modules/k8s_vault_verify_version"

  vault_instance_count   = var.instance_count
  vault_product_version  = var.vault_version
  vault_product_revision = var.vault_revision
}

module "k8s_verify_write_data" {
  source = "../modules/k8s_vault_verify_write_data"

  vault_instance_count = var.instance_count
}

module "read_license" {
  source = "../modules/read_license"
}
