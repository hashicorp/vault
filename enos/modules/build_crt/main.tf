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
