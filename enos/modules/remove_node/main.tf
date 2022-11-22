terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "node_public_ips" {
  type        = list(string)
  description = "Vault primary cluster follower Public IP addresses"
  default = [""]
}

variable "node_private_ip" {
  type        = string
  description = "Vault primary cluster leader Public IP address"
  default = ""
}

variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

locals {
  nodes = toset([for idx in range(var.vault_instance_count - 1) : tostring(idx)])
  # private_ip = var.node_private_ip != "" ? split(",", var.node_private_ip) : [""]
  node_ip = coalescelist(var.node_public_ips, split(",", var.node_private_ip))
}

resource "enos_remote_exec" "remove_node" {
  for_each = local.nodes

  inline = ["sudo halt -f -n --force --no-wall |at now + 1 minute; exit 0"]
  # inline = ["sudo shutdown -H --no-wall; exit 0"]

  transport = {
    ssh = {
      host = element(local.node_ip, each.key)
    }
  }
}
