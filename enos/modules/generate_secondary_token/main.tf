terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.4.3"
    }
  }
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "primary_leader_public_ip" {
  type        = string
  description = "Vault primary cluster leader Public IP address"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

locals {
  token_id        = random_uuid.token_id.id
  secondary_token = enos_remote_exec.fetch_secondary_token.stdout
}
resource "random_uuid" "token_id" {}

resource "enos_remote_exec" "fetch_secondary_token" {
  depends_on = [random_uuid.token_id]
  environment = {
    VAULT_ADDR  = "http://127.0.0.1:8200"
    VAULT_TOKEN = var.vault_root_token
  }

  inline = ["${var.vault_install_dir}/vault write sys/replication/performance/primary/secondary-token id=${local.token_id} |sed -n '/^wrapping_token:/p' |awk '{print $2}'"]

  transport = {
    ssh = {
      host = var.primary_leader_public_ip
    }
  }
}

output "secondary_token" {
  value = local.secondary_token
}
