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
  environment = {
    VAULT_INSTALL_DIR   = var.vault_install_dir
    VAULT_TOKEN         = var.vault_root_token
    VAULT_PROXY_PIDFILE = var.vault_proxy_pidfile
    VAULT_PROXY_ADDRESS = local.vault_proxy_address
  }

  scripts = [abspath("${path.module}/scripts/set-up-approle-and-proxy.sh")]

  transport = {
    ssh = {
      host = local.vault_instances[0].public_ip
    }
  }
}

resource "enos_remote_exec" "use_proxy" {
  environment = {
    VAULT_INSTALL_DIR   = var.vault_install_dir
    VAULT_PROXY_PIDFILE = var.vault_proxy_pidfile
    VAULT_PROXY_ADDRESS = local.vault_proxy_address
  }

  scripts = [abspath("${path.module}/scripts/use-proxy.sh")]

  transport = {
    ssh = {
      host = local.vault_instances[0].public_ip
    }
  }

  depends_on = [
    enos_remote_exec.set_up_approle_auth_and_proxy
  ]
}
