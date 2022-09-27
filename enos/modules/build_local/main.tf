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

variable "artifactory_host" {default = null}
variable "artifactory_repo" {default = null}
variable    "artifactory_username" {default = null}
variable    "artifactory_token" {default = null}
variable "arch" {
  default = null
}
variable "artifact_path" {
  default = null
}
variable "artifact_type" {
  default = null
}
variable "distro" {
  default = null
}
variable "edition" {
  default = null
}
variable "instance_type" {
  default = null
}
variable "revision" {
  default = null
}
variable "vault_product_version" {
  default = null
}

resource "enos_local_exec" "build" {
  content = templatefile("${path.module}/templates/build.sh", {
    bundle_path = var.bundle_path,
    build_tags  = join(" ", var.build_tags)
    goarch      = var.goarch
    goos        = var.goos
  })
}
