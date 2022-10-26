terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "undo_logs" {
  type        = string
  description = "The variable takes either 0 or 1 to indicate disabling or enabling undo logs"
}

resource "enos_local_exec" "set_undo_logs_env_var" {
  inline = ["export VAULT_REPLICATION_USE_UNDO_LOGS=${var.undo_logs}"]
}
