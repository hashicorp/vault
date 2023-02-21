
terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

locals {
  instances = toset([for idx in range(var.vault_instance_count) : tostring(idx)])
}

resource "enos_remote_exec" "replication_status" {
  for_each = local.instances

  inline = ["vault read -format=json sys/replication/status"]

  transport = {
    kubernetes = {
      kubeconfig_base64 = var.kubeconfig_base64
      context_name      = var.context_name
      pod               = var.vault_pods[each.key].name
      namespace         = var.vault_pods[each.key].namespace
    }
  }
}

resource "enos_local_exec" "verify_replication_status" {

  for_each = enos_remote_exec.replication_status

  environment = {
    STATUS        = each.value.stdout
    VAULT_EDITION = var.vault_edition
  }

  content = abspath("${path.module}/scripts/smoke-verify-replication.sh")
}
