terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "node_public_ip" {
  type        = string
  description = "Node Public IP address"
}

resource "enos_remote_exec" "shutdown_node" {
  inline = ["sudo shutdown -H --no-wall; exit 0"]

  transport = {
    ssh = {
      host = var.node_public_ip
    }
  }
}
