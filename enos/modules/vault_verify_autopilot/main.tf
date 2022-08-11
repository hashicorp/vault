
terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

resource "enos_remote_exec" "smoke-verify-autopilot" {
  for_each = var.vault_instances

  content = templatefile("${path.module}/templates/smoke-verify-autopilot.sh", {
    vault_install_dir               = var.vault_install_dir,
    vault_token                     = var.vault_token,
    vault_autopilot_upgrade_status  = var.vault_autopilot_upgrade_status,
    vault_autopilot_upgrade_version = var.vault_autopilot_upgrade_version,
  })

  transport = {
    ssh = {
      host = each.value.public_ip
      user = var.enos_transport_user
    }
  }
}
