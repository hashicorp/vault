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

variable "build_tags" {
  type        = list(string)
  description = "The build tags to pass to the Go compiler"
}

variable "goarch" {
  type        = string
  description = "The Go architecture target"
  default     = "amd64"
}

variable "goos" {
  type        = string
  description = "The Go OS target"
  default     = "linux"
}

resource "enos_local_exec" "build" {
  content = templatefile("${path.module}/templates/build.sh", {
    bundle_path = var.bundle_path,
    build_tags  = join(" ", var.build_tags)
    goarch      = var.goarch
    goos        = var.goos
  })
}
