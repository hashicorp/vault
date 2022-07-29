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

variable "vault_local_binary_path" {
  type        = string
  description = "The path to the local vault binary to compare version information"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

locals {
  public_ips = {
    for idx in range(length(var.vault_instance_public_ips)) : idx => {
      public_ip = var.vault_instance_public_ips[idx]
    }
  }
}

resource "enos_local_exec" "get_expected_version" {
  content = templatefile("${path.module}/templates/get-local-version.sh", {
    vault_local_binary_path = var.vault_local_binary_path
  })
}

resource "enos_remote_exec" "verify_all_nodes_have_updated_version" {
  for_each = local.public_ips

  content = templatefile("${path.module}/templates/verify-cluster-version.sh", {
    vault_install_dir = var.vault_install_dir,
    vault_token       = var.vault_root_token,
    expected_version  = trimspace(enos_local_exec.get_expected_version.stdout)
  })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
