# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

module "create_kind_cluster" {
  source = "../modules/local_kind_cluster"
}

module "create_aks_cluster" {
  source = "../modules/azure_aks_cluster"

  node_count      = 3
  kubeconfig_path = var.kubeconfig_path
}

module "load_docker_image" {
  source = "../modules/load_docker_image"
}

module "k8s_deploy_vault" {
  source = "../modules/k8s_deploy_vault"

  vault_instance_count = var.vault_instance_count
}

module "k8s_verify_build_date" {
  source = "../modules/k8s_vault_verify_build_date"

  vault_instance_count = var.vault_instance_count
}

module "k8s_verify_replication" {
  source = "../modules/k8s_vault_verify_replication"

  vault_instance_count = var.vault_instance_count
}

module "k8s_verify_ui" {
  source = "../modules/k8s_vault_verify_ui"

  vault_instance_count = var.vault_instance_count
}

module "k8s_verify_version" {
  source = "../modules/k8s_vault_verify_version"

  vault_instance_count   = var.vault_instance_count
  vault_product_version  = var.vault_product_version
  vault_product_revision = var.vault_product_revision
}

module "k8s_verify_write_data" {
  source = "../modules/k8s_vault_verify_write_data"

  vault_instance_count = var.vault_instance_count
}

module "read_license" {
  source = "../modules/read_license"
}

module "setup_azure_auth_backend" {
  source = "../modules/setup_azure_auth_backend"
}

module "k8s_azure_wif_extra_chart_values" {
  source = "../modules/k8s_azure_wif_extra_chart_values"
}

module "setup_azure_auth_method" {
  source = "../modules/vault_setup_azure_auth_method"
}

module "format_transport_block" {
  source = "../modules/format_transport_block"
}

module "deploy_vault_agent" {
  source = "../modules/k8s_deploy_vault_agent"
}
