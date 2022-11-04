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

variable "secondary_leader_public_ip" {
  type        = string
  description = "Vault secondary cluster leader Public IP address"
}

variable "secondary_leader_private_ip" {
  type        = string
  description = "Vault secondary cluster leader Private IP address"
}

variable "secondary_vault_root_token" {
  type        = string
  description = "The vault root token of secondary cluster"
}

resource "enos_remote_exec" "verify_kv_on_secondary" {
  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_TOKEN       = var.secondary_vault_root_token
  }

  inline = ["${var.vault_install_dir}/vault operator raft list-peers"]

  transport = {
    ssh = {
      host = var.secondary_leader_public_ip
    }
  }
}

output "list-peers" {
  value = enos_remote_exec.verify_kv_on_secondary.stdout
}