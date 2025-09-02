# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The Vault cluster instances that were created"
}

variable "ip_version" {
  type        = number
  description = "The IP version to use for the Vault TCP listeners"

  validation {
    condition     = contains([4, 6], var.ip_version)
    error_message = "The ip_version must be either 4 or 6"
  }
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
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

variable "vault_proxy_port" {
  type        = number
  description = "The Vault Proxy listener port"
}

variable "vault_root_token" {
  type        = string
  description = "The Vault root token"
}

locals {
  vault_proxy_address = "${var.ip_version == 4 ? "127.0.0.1" : "[::1]"}:${var.vault_proxy_port}"
}

resource "enos_remote_exec" "set_up_approle_auth_and_proxy" {
  environment = {
    VAULT_ADDR          = var.vault_addr
    VAULT_INSTALL_DIR   = var.vault_install_dir
    VAULT_PROXY_ADDRESS = local.vault_proxy_address
    VAULT_PROXY_PIDFILE = var.vault_proxy_pidfile
    VAULT_TOKEN         = var.vault_root_token
  }

  scripts = [abspath("${path.module}/scripts/set-up-approle-and-proxy.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
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
      host = var.hosts[0].public_ip
    }
  }

  depends_on = [
    enos_remote_exec.set_up_approle_auth_and_proxy
  ]
}
