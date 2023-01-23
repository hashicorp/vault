# Shim module since CRT provided things will use the crt_bundle_path variable
variable "bundle_path" {
  default = "/tmp/vault.zip"
}

variable "build_tags" {
  default = ["ui"]
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

variable "artifactory_host" { default = null }
variable "artifactory_repo" { default = null }
variable "artifactory_username" { default = null }
variable "artifactory_token" { default = null }
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
variable "product_version" {
  default = null
}
