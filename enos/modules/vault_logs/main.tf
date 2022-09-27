terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "vault_instances" {
  type = map(object({
    private_ip = string
    public_ip  = string
  }))
  description = "The vault cluster instances that were created"
}


variable "vault_instance_count" {
  type        = number
  description = "How many vault instances are in the cluster"
}

locals {
  public_ips = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "logs" {
  for_each = local.public_ips

  content = file("${path.module}/templates/logs.sh")

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

output "logs" {
  value = {
    for idx in range(var.vault_instance_count) : idx => {
      logs = enos_remote_exec.logs[idx].stdout
      }
    }
}
