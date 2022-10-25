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

variable "vault_api_addr" {
  type        = string
  description = "The API address of the Vault cluster"
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

variable "vault_local_bundle_path" {
  type        = string
  description = "The path to the local Vault (vault.zip) bundle"
}

variable "vault_local_artifact_path" {
  type        = string
  description = "The path to a locally built vault artifact to install"
  default     = null
}

variable "vault_artifactory_release" {
  type = object({
    username = string
    token    = string
    url      = string
    sha256   = string
  })
  description = "Vault release version and edition to install from artifactory.hashicorp.engineering"
  default     = null
}

variable "vault_seal_type" {
  type        = string
  description = "The Vault seal type"
}

variable "vault_unseal_keys" {
  type        = list(string)
  description = "The keys to use to unseal Vault when not using auto-unseal"
  default     = null
}

locals {
  instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
  followers      = toset([for idx in range(var.vault_instance_count - 1) : tostring(idx)])
  follower_ips   = compact(split(" ", enos_remote_exec.get_follower_public_ips.stdout))
  vault_bin_path = "${var.vault_install_dir}/vault"
}

resource "enos_bundle_install" "upgrade_vault_binary" {
  for_each = local.instances

  destination = var.vault_install_dir
  artifactory = var.vault_artifactory_release
  path        = var.vault_local_bundle_path

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

resource "enos_remote_exec" "get_leader_public_ip" {
  depends_on = [enos_bundle_install.upgrade_vault_binary]

  content = templatefile("${path.module}/templates/get-leader-public-ip.sh", {
    vault_install_dir = var.vault_install_dir,
    vault_instances   = jsonencode(local.instances)
  })

  transport = {
    ssh = {
      host = local.instances[0].public_ip
    }
  }
}

resource "enos_remote_exec" "get_follower_public_ips" {
  depends_on = [enos_bundle_install.upgrade_vault_binary]

  content = templatefile("${path.module}/templates/get-follower-public-ips.sh", {
    vault_install_dir = var.vault_install_dir,
    vault_instances   = jsonencode(local.instances)
  })

  transport = {
    ssh = {
      host = local.instances[0].public_ip
    }
  }
}

resource "enos_remote_exec" "restart_followers" {
  for_each   = local.followers
  depends_on = [enos_remote_exec.get_follower_public_ips]

  content = file("${path.module}/templates/restart-vault.sh")

  transport = {
    ssh = {
      host = trimspace(local.follower_ips[tonumber(each.key)])
    }
  }
}

resource "enos_vault_unseal" "followers" {
  depends_on = [enos_remote_exec.restart_followers]
  for_each = {
    for idx, follower in local.followers : idx => follower
    if var.vault_seal_type == "shamir"
  }
  bin_path    = local.vault_bin_path
  vault_addr  = var.vault_api_addr
  seal_type   = var.vault_seal_type
  unseal_keys = var.vault_unseal_keys

  transport = {
    ssh = {
      host = trimspace(local.follower_ips[each.key])
    }
  }
}

resource "enos_remote_exec" "restart_leader" {
  depends_on = [enos_vault_unseal.followers]

  content = file("${path.module}/templates/restart-vault.sh")

  transport = {
    ssh = {
      host = trimspace(enos_remote_exec.get_leader_public_ip.stdout)
    }
  }
}

resource "enos_vault_unseal" "leader" {
  count      = var.vault_seal_type == "shamir" ? 1 : 0
  depends_on = [enos_remote_exec.restart_leader]

  bin_path    = local.vault_bin_path
  vault_addr  = var.vault_api_addr
  seal_type   = var.vault_seal_type
  unseal_keys = var.vault_unseal_keys

  transport = {
    ssh = {
      host = trimspace(enos_remote_exec.get_leader_public_ip.stdout)
    }
  }
}
