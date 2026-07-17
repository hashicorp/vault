# Copyright IBM Corp. 2016, 2025
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

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "vault_ibm_license_edition" {
  type        = string
  description = "The expected Vault edition in use in the IBM PAO license."
}

resource "enos_remote_exec" "vault_verify_ibm_license_update" {
  for_each = var.hosts

  environment = {
    VAULT_ADDR                = var.vault_addr
    VAULT_CLUSTER_ADDR        = "${each.value.private_ip}:${var.vault_cluster_addr_port}"
    VAULT_TOKEN               = var.vault_root_token
    VAULT_INSTALL_DIR         = var.vault_install_dir
    VAULT_IBM_LICENSE_EDITION = var.vault_ibm_license_edition
  }

  scripts = [abspath("${path.module}/scripts/license-get.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "vault_get_ibm_license_customer_id" {
  depends_on = [enos_remote_exec.vault_verify_ibm_license_update]
  for_each   = var.hosts

  environment = {
    VAULT_ADDR         = var.vault_addr
    VAULT_CLUSTER_ADDR = "${each.value.private_ip}:${var.vault_cluster_addr_port}"
    VAULT_TOKEN        = var.vault_root_token
    VAULT_INSTALL_DIR  = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/scripts/license-inspect.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

output "customer_id" {
  value = enos_remote_exec.vault_get_ibm_license_customer_id[0].stdout
}
