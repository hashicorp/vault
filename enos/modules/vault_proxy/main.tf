# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_root_token" {
  type        = string
  description = "The Vault root token"
}

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The Vault cluster instances that were created"
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_proxy_pidfile" {
  type        = string
  description = "The filepath where the Vault Proxy pid file is kept"
  default     = "/tmp/pidfile"
}

locals {
  vault_instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
  vault_proxy_address = "127.0.0.1:8100"
}

resource "enos_remote_exec" "set_up_approle_auth_and_proxy" {
  content = templatefile("${path.module}/templates/set-up-approle-and-proxy.sh", {
    vault_install_dir   = var.vault_install_dir
    vault_token         = var.vault_root_token
    vault_proxy_pidfile = var.vault_proxy_pidfile
    vault_proxy_address = local.vault_proxy_address
  })

  transport = {
    ssh = {
      host = local.vault_instances[0].public_ip
    }
  }
}

resource "enos_remote_exec" "use_proxy" {
  content = templatefile("${path.module}/templates/use-proxy.sh", {
    vault_install_dir   = var.vault_install_dir
    vault_proxy_pidfile = var.vault_proxy_pidfile
    vault_proxy_address = local.vault_proxy_address
  })

  transport = {
    ssh = {
      host = local.vault_instances[0].public_ip
    }
  }

  depends_on = [
    enos_remote_exec.set_up_approle_auth_and_proxy
  ]
}
