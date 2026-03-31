# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}

variable "vault_ibm_license_path" {
  type        = string
  description = "The path to the IBM PAO license file on the target machine"
}

variable "vault_ibm_license_edition" {
  type        = string
  description = "The edition corresponding to the Vault entitlement to select in the IBM PAO license"
}

resource "enos_remote_exec" "vault_update_license_ibm" {
  for_each = var.hosts

  environment = {
    VAULT_IBM_LICENSE         = file(var.vault_ibm_license_path)
    VAULT_IBM_LICENSE_EDITION = var.vault_ibm_license_edition
  }

  scripts = [abspath("${path.module}/scripts/update-license.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
