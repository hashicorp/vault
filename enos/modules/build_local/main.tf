# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "artifact_path" {
  description = "Where to create the zip bundle of the Vault build"
}

variable "build_tags" {
  type        = list(string)
  description = "The build tags to pass to the Go compiler"
}

variable "build_ui" {
  type        = bool
  description = "Whether or not we should build the UI when creating the local build"
  default     = true
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
variable "artifact_type" { default = null }
variable "distro" { default = null }
variable "distro_version" { default = null }
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
    BIN_PATH           = abspath("${path.module}/../../../dist")
    BUILD_UI           = tostring(var.build_ui)
    BUNDLE_PATH        = abspath(var.artifact_path)
    GO_TAGS            = join(" ", var.build_tags)
    GOARCH             = var.goarch
    GOOS               = var.goos
    PRERELEASE_VERSION = module.local_metadata.version_pre
    VERSION_METADATA   = module.local_metadata.version_meta
  }
}
