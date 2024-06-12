# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "cluster_id" {
  type = string
}

variable "hosts" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The hosts that will have access to the softhsm"
}

locals {
  pin             = resource.random_string.pin.result
  aes_label       = "vault_hsm_aes_${local.pin}"
  hmac_label      = "vault_hsm_hmac_${local.pin}"
  seal_attributes = jsondecode(resource.enos_remote_exec.create_keys.stdout)
  target          = tomap({ "1" = var.hosts[0] })
  token           = "${var.cluster_id}_${local.pin}"
}

resource "random_string" "pin" {
  length  = 5
  lower   = true
  upper   = false
  numeric = true
  special = false
}

module "install" {
  source = "../softhsm_install"

  hosts         = local.target
  include_tools = true # make sure opensc is also installed as we need it to create keys
}

module "initialize" {
  source     = "../softhsm_init"
  depends_on = [module.install]

  hosts = local.target
}

// Create our keys. Our stdout contains the requried the values for the pksc11 seal stanza
// as JSON. https://developer.hashicorp.com/vault/docs/configuration/seal/pkcs11#pkcs11-parameters
resource "enos_remote_exec" "create_keys" {
  depends_on = [
    module.install,
    module.initialize,
  ]

  environment = {
    AES_LABEL   = local.aes_label
    HMAC_LABEL  = local.hmac_label
    PIN         = resource.random_string.pin.result
    TOKEN_DIR   = module.initialize.token_dir
    TOKEN_LABEL = local.token
    SO_PIN      = resource.random_string.pin.result
  }

  scripts = [abspath("${path.module}/scripts/create-keys.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

// Get our softhsm token. Stdout is a base64 encoded gzipped tarball of the softhsm token dir. This
// allows us to pass around binary data inside of Terraform's type system.
resource "enos_remote_exec" "get_keys" {
  depends_on = [enos_remote_exec.create_keys]

  environment = {
    TOKEN_DIR = module.initialize.token_dir
  }

  scripts = [abspath("${path.module}/scripts/get-keys.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

output "seal_attributes" {
  description = "Seal device specific attributes. Contains all required keys for the seal stanza"
  value       = local.seal_attributes
}

output "token_base64" {
  description = "The softhsm token and keys gzipped tarball in base64"
  value       = enos_remote_exec.get_keys.stdout
}

output "token_dir" {
  description = "The softhsm directory where tokens and keys are stored"
  value       = module.initialize.token_dir
}

output "token_label" {
  description = "The HSM slot token label"
  value       = local.token
}

output "all_attributes" {
  description = "Seal device specific attributes"
  value = merge(
    local.seal_attributes,
    {
      token_base64 = enos_remote_exec.get_keys.stdout,
      token_dir    = module.initialize.token_dir
    },
  )
}
