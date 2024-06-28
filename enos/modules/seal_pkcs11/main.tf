# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

/*

A seal module that emulates using a real PKCS#11 HSM. For this we'll use softhsm2. You'll
need softhsm2 and opensc installed to get access to the userspace tools and dynamic library that
Vault Enterprise will use. Here we'll take in the vault hosts and use the one of the nodes
to generate the hsm slot and the tokens, and then we'll copy the softhsm tokens to the other nodes.

Using softhsm2 and opensc is a bit complicated but here's a cheat sheet for getting started.

$ brew install softhsm opensc
or
$ sudo apt install softhsm2 opensc

Create a softhsm slot. You can use anything you want for the pin and the supervisor pin. This will
output the slot identifier, which you'll use as the `slot` parameter in the seal config.
$ softhsm2-util --init-token --free --so-pin=1234 --pin=1234 --label="seal" | grep -oE '[0-9]+$'

You can see the slots:
$ softhsm2-util --show-slots
Or use opensc's pkcs11-tool. Make sure to use your pin for the -p flag. The module that we refer
to is the location of the shared library that we need to provide to Vault Enterprise. Depending on
your platform or installation method this could be different.
$ pkcs11-tool --module /usr/local/Cellar/softhsm/2.6.1/lib/softhsm/libsofthsm2.so -a seal -p 1234 -IL

Find yours
$ find /usr/local -type f -name libsofthsm2.so -print -quit

Your tokens will be installed in the default directories.tokendir. See man softhsm2.conf(5) for
more details. On macOS from brew this is /usr/local/var/lib/softhsm/tokens/

Vault Enterprise supports creating the HSM keys, but for softhsm2 that would require us to
initialize with one node before copying the contents. So instead we'll create an HSM key and HMAC
key that we'll copy everywhere.

$ pkcs11-tool --module /usr/local/Cellar/softhsm/2.6.1/lib/softhsm/libsofthsm2.so -a seal -p 1234 --token-label seal --keygen --usage-sign --label hsm_hmac --id 1 --key-type GENERIC:32 --private --sensitive
$ pkcs11-tool --module /usr/local/Cellar/softhsm/2.6.1/lib/softhsm/libsofthsm2.so -a seal -p 1234 --token-label seal --keygen --usage-sign --label hsm_aes --id 2 --key-type AES:32 --private --sensitive --usage-wrap

Now you should be able to configure Vault Enterprise seal stanza.
*/

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "cluster_id" {
  type        = string
  description = "The VPC ID of the cluster"
}

variable "cluster_meta" {
  type        = string
  default     = null
  description = "Any metadata that needs to be passed in. If we're creating multiple softhsm tokens this value could be a prior KEYS_BASE64"
}

variable "cluster_ssh_keypair" {
  type        = string
  description = "The ssh keypair of the vault cluster. We need this to used the inherited provider for our target"
}

variable "common_tags" {
  type    = map(string)
  default = null
}

variable "other_resources" {
  type    = list(string)
  default = []
}

resource "random_string" "id" {
  length  = 8
  numeric = false
  special = false
  upper   = false
}

module "ec2_info" {
  source = "../ec2_info"
}

locals {
  id = "${var.cluster_id}-${random_string.id.result}"
}

module "target" {
  source          = "../target_ec2_instances"
  ami_id          = module.ec2_info.ami_ids["arm64"]["ubuntu"]["22.04"]
  cluster_tag_key = local.id
  common_tags     = var.common_tags
  instance_count  = 1
  instance_types = {
    amd64 = "t3a.small"
    arm64 = "t4g.small"
  }
  // Make sure it's not too long as we use this for aws resources that size maximums that are easy
  // to hit.
  project_name = substr("vault-ci-softhsm-${local.id}", 0, 32)
  ssh_keypair  = var.cluster_ssh_keypair
  vpc_id       = var.cluster_id
}

module "create_vault_keys" {
  source = "../softhsm_create_vault_keys"

  cluster_id = var.cluster_id
  hosts      = module.target.hosts
}

// Our attributes contain all required keys for the seal stanza and our base64 encoded softhsm
// token and keys.
output "attributes" {
  description = "Seal device specific attributes"
  value       = module.create_vault_keys.all_attributes
}

// Shim for chaining seals that require IAM roles
output "resource_name" { value = null }
output "resource_names" { value = var.other_resources }
