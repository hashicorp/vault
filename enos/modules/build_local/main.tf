terraform {
  required_providers {
    enos = {
      source = "hashicorp.com/qti/enos"
    }
  }
}

variable "bundle_path" {
  type    = string
  default = "/tmp/vault.zip"
}

variable "local_vault_artifact_path" {
  type    = string
  default = "/tmp/vault"
}

resource "enos_local_exec" "build" {
  content = templatefile("${path.module}/templates/build.sh", {
    bundle_path = var.bundle_path,
    vault_path  = var.local_vault_artifact_path
  })
}
