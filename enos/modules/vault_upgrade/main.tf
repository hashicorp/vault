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

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_instance_public_ips" {
  type        = list(string)
  description = "The public IP addresses to the Vault cluster instances"
}

variable "vault_local_bundle_path" {
  type        = string
  description = "The path to the local Vault (vault.zip) bundle"
}

locals {
  public_ips = {
    for idx in range(length(var.vault_instance_public_ips)) : idx => {
      public_ip = var.vault_instance_public_ips[idx]
    }
  }
}

resource "enos_bundle_install" "upgrade_vault_binary" {
  for_each = local.public_ips

  destination = var.vault_install_dir
  path        = var.vault_local_bundle_path

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "upgrade_standby" {
  for_each   = local.public_ips
  depends_on = [enos_bundle_install.upgrade_vault_binary]

  content = templatefile("${path.module}/templates/vault-upgrade.sh", {
    vault_install_dir = var.vault_install_dir,
    upgrade_target    = "standby"
  })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "upgrade_active" {
  for_each   = local.public_ips
  depends_on = [enos_remote_exec.upgrade_standby]

  content = templatefile("${path.module}/templates/vault-upgrade.sh", {
    vault_install_dir = var.vault_install_dir,
    upgrade_target    = "active"
  })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
