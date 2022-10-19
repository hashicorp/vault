
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

resource "enos_remote_exec" "smoke-enable-secrets-kv" {
  for_each = local.instances

  content = templatefile("${path.module}/templates/smoke-enable-secrets-kv.sh", {
    instance_id       = each.key
    vault_install_dir = var.vault_install_dir,
    vault_token       = var.vault_root_token,
  })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

# Verify that we can enable the k/v secrets engine and write data to it.
resource "enos_remote_exec" "smoke-write-test-data" {
  depends_on = [enos_remote_exec.smoke-enable-secrets-kv]
  for_each   = local.instances

  content = templatefile("${path.module}/templates/smoke-write-test-data.sh", {
    test_key          = "smoke${each.key}"
    test_value        = "fire"
    vault_install_dir = var.vault_install_dir,
    vault_token       = var.vault_root_token,
  })

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
