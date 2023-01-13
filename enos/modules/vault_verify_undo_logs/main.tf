terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

variable "vault_instances" {
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

variable "vault_undo_logs_status" {
  type        = string
  description = "An integer either 0 or 1 which indicates whether undo_logs are disabled or enabled"
  default     = null
}

locals {
  public_ips = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "smoke-verify-undo-logs" {
  for_each = local.public_ips

  environment = {
    VAULT_TOKEN            = var.vault_root_token
    VAULT_ADDR             = "http://localhost:8200"
    VAULT_UNDO_LOGS_STATUS = var.vault_undo_logs_status
  }

  scripts = [abspath("${path.module}/scripts/smoke-verify-undo-logs.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
