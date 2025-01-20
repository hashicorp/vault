# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "seal_type" {
  type        = string
  description = "The expected seal type"
  default     = "shamir"
}

resource "enos_remote_exec" "verify_seal_type" {
  for_each = var.vault_hosts

  scripts = [abspath("${path.module}/scripts/verify-seal-type.sh")]

  environment = {
    VAULT_ADDR         = "http://127.0.0.1:8200"
    VAULT_INSTALL_DIR  = var.vault_install_dir
    EXPECTED_SEAL_TYPE = var.seal_type
  }

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
