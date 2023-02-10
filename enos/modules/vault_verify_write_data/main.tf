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
  description = "How many Vault instances are in the cluster"
}

variable "leader_public_ip" {
  type        = string
  description = "Vault cluster leader Public IP address"
}

variable "leader_private_ip" {
  type        = string
  description = "Vault cluster leader Private IP address"
}

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The Vault cluster instances that were created"
}

variable "vault_root_token" {
  type        = string
  description = "The Vault root token"
  default     = null
}

locals {
  instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

# We use this module to verify write data in all Enos scenarios.  Since we cannot use
# Vault token to authenticate to secondary clusters in replication scenario we add a regular user
# here to keep the authentication method and module verification consistent between all scenarios
resource "enos_remote_exec" "smoke-enable-secrets-kv" {
  # Only enable the secrets engine on the leader node
  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = ["${path.module}/scripts/smoke-enable-secrets-kv.sh"]

  transport = {
    ssh = {
      host = var.leader_public_ip
    }
  }
}

# Verify that we can enable the k/v secrets engine and write data to it.
resource "enos_remote_exec" "smoke-write-test-data" {
  depends_on = [enos_remote_exec.smoke-enable-secrets-kv]
  for_each   = local.instances

  environment = {
    VAULT_ADDR        = "http://127.0.0.1:8200"
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
    TEST_KEY          = "smoke${each.key}"
    TEST_VALUE        = "fire"
  }

  scripts = ["${path.module}/scripts/smoke-write-test-data.sh"]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
