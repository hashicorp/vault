# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_cluster_addr_port" {
  description = "The Raft cluster address port"
  type        = string
  default     = "8201"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

resource "enos_remote_exec" "vault_verify_billing_start_date" {
  for_each = var.hosts

  environment = {
    VAULT_ADDR              = var.vault_addr
    VAULT_CLUSTER_ADDR      = "${each.value.private_ip}:${var.vault_cluster_addr_port}"
    VAULT_INSTALL_DIR       = var.vault_install_dir
    VAULT_LOCAL_BINARY_PATH = "${var.vault_install_dir}/vault"
    VAULT_TOKEN             = var.vault_root_token
  }

  scripts = [abspath("${path.module}/scripts/verify-billing-start.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
