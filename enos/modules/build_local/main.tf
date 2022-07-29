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

resource "enos_local_exec" "build" {
  content = templatefile("${path.module}/templates/build.sh", {
    bundle_path = var.bundle_path,
  })
}
