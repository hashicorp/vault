
variable "artifactory_username" {
  type        = string
  description = "The username to use when connecting to artifactory"
  default     = null
}

variable "artifactory_token" {
  type        = string
  description = "The token to use when connecting to artifactory"
  default     = null
  sensitive   = true
}

variable "artifactory_host" {
  type        = string
  description = "The artifactory host to search for vault artifacts"
  default     = "https://artifactory.hashicorp.engineering/artifactory"
}

variable "artifactory_repo" {
  type        = string
  description = "The artifactory repo to search for vault artifacts"
  default     = "hashicorp-crt-stable-local*"
}
variable "arch" {}
variable "artifact_type" {}
variable "distro" {}
variable "edition" {}
variable "instance_type" {}
variable "revision" {}
variable "product_version" {}
variable "build_tags" { default = null }
variable "bundle_path" { default = null }
variable "goarch" { default = null }
variable "goos" { default = null }
