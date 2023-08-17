
terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

locals {
  instances = {
    for idx in range(var.vault_instance_count) : idx => {
      public_ip  = values(var.vault_instances)[idx].public_ip
      private_ip = values(var.vault_instances)[idx].private_ip
    }
  }
}

resource "enos_remote_exec" "smoke-verify-ui" {
  for_each = local.instances

  content = templatefile("${path.module}/templates/smoke-verify-ui.sh", {
    vault_install_dir = var.vault_install_dir,
  })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
