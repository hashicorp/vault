# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
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

variable "artifactory_host" { default = null }
variable "artifactory_repo" { default = null }
variable "artifactory_username" { default = null }
variable "artifactory_token" { default = null }
variable "arch" { default = null }
variable "artifact_path" { default = null }
variable "artifact_type" { default = null }
variable "distro" { default = null }
variable "edition" { default = null }
variable "revision" { default = null }
variable "product_version" { default = null }

module "local_metadata" {
  source = "../get_local_metadata"
}

resource "enos_local_exec" "build" {
  scripts = [abspath("${path.module}/scripts/build.sh")]

  environment = {
    BASE_VERSION       = module.local_metadata.version_base
    BIN_PATH           = "dist"
    BUNDLE_PATH        = var.bundle_path,
    GO_TAGS            = join(" ", var.build_tags)
    GOARCH             = var.goarch
    GOOS               = var.goos
    PRERELEASE_VERSION = module.local_metadata.version_pre
    VERSION_METADATA   = module.local_metadata.version_meta
  }
}
