# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

# Vault Local Build Module
#
# This module builds Vault binaries locally and produces artifacts for different use cases.
# The module enforces that at least one output artifact is created to prevent silent no-op execution.
#
# Supported workflows:
# 1. ZIP Bundle Only:     artifact_path specified, docker_bin_path null
#    - Builds binary to dist/
#    - Creates zip bundle at artifact_path
#
# 2. Docker Only:         docker_bin_path specified, artifact_path null
#    - Builds binary to dist/
#    - Copies binary to docker_bin_path for Docker image creation
#    - No zip bundle created
#
# 3. Both ZIP and Docker: Both artifact_path and docker_bin_path specified
#    - Builds binary to dist/
#    - Creates zip bundle at artifact_path
#    - Copies binary to docker_bin_path for Docker image creation
#
# The validation ensures at least one of artifact_path or docker_bin_path is specified.

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "artifact_path" {
  type        = string
  description = "Where to create the zip bundle of the Vault build. If null, no zip bundle will be created."
  default     = null
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

variable "docker_bin_path" {
  type        = string
  description = "Path to copy the built binary for Docker image creation. When specified, the binary is copied to this location for subsequent Docker image builds. If null, no Docker-specific binary copy is performed."
  default     = null

  validation {
    condition     = var.artifact_path != null || var.docker_bin_path != null
    error_message = "At least one of 'artifact_path' (for zip bundle) or 'docker_bin_path' (for Docker builds) must be specified. The module must produce at least one output artifact."
  }
}

variable "artifactory_host" { default = null }
variable "artifactory_repo" { default = null }
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
    BUNDLE_PATH        = var.artifact_path != null ? abspath(var.artifact_path) : ""
    GO_TAGS            = join(" ", var.build_tags)
    GOARCH             = var.goarch
    GOOS               = var.goos
    PRERELEASE_VERSION = module.local_metadata.version_pre
    VERSION_METADATA   = module.local_metadata.version_meta
    CUSTOM_BIN_PATH    = var.docker_bin_path != null ? abspath("${path.module}/../../../${var.docker_bin_path}") : ""
  }
}
