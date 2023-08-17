# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "artifactory_username" {
  type        = string
  description = "The username to use when testing an artifact from artifactory"
  default     = null
  sensitive   = true
}

variable "artifactory_token" {
  type        = string
  description = "The token to use when authenticating to artifactory"
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

variable "aws_region" {
  description = "The AWS region where we'll create infrastructure"
  type        = string
  default     = "us-east-1"
}

variable "aws_ssh_keypair_name" {
  description = "The AWS keypair to use for SSH"
  type        = string
  default     = "enos-ci-ssh-key"
}

variable "aws_ssh_private_key_path" {
  description = "The path to the AWS keypair private key"
  type        = string
  default     = "./support/private_key.pem"
}

variable "backend_instance_type" {
  description = "The instance type to use for the Vault backend. Must be arm64/nitro compatible"
  type        = string
  default     = "t4g.small"
}

variable "backend_license_path" {
  description = "The license for the backend if applicable (Consul Enterprise)"
  type        = string
  default     = null
}

variable "backend_log_level" {
  description = "The server log level for the backend. Supported values include 'trace', 'debug', 'info', 'warn', 'error'"
  type        = string
  default     = "trace"
}

variable "project_name" {
  description = "The description of the project"
  type        = string
  default     = "vault-enos-integration"
}

variable "rhel_distro_version" {
  description = "The version of RHEL to use"
  type        = string
  default     = "9.1" // or "8.8"
}

variable "tags" {
  description = "Tags that will be applied to infrastructure resources that support tagging"
  type        = map(string)
  default     = null
}

variable "terraform_plugin_cache_dir" {
  description = "The directory to cache Terraform modules and providers"
  type        = string
  default     = null
}

variable "tfc_api_token" {
  description = "The Terraform Cloud QTI Organization API token. This is used to download the enos Terraform provider."
  type        = string
  sensitive   = true
}

variable "ubuntu_distro_version" {
  description = "The version of ubuntu to use"
  type        = string
  default     = "22.04" // or "20.04", "18.04"
}

variable "ui_test_filter" {
  type        = string
  description = "A test filter to limit the ui tests to execute. Will be appended to the ember test command as '-f=\"<filter>\"'"
  default     = null
}

variable "ui_run_tests" {
  type        = bool
  description = "Whether to run the UI tests or not. If set to false a cluster will be created but no tests will be run"
  default     = true
}

variable "vault_artifact_type" {
  description = "The type of Vault artifact to use when installing Vault from artifactory. It should be 'package' for .deb or # .rpm package and 'bundle' for .zip bundles"
  default     = "bundle"
}

variable "vault_autopilot_initial_release" {
  description = "The Vault release to deploy before upgrading with autopilot"
  default = {
    edition = "ent"
    version = "1.11.0"
  }
}

variable "vault_artifact_path" {
  description = "Path to CRT generated or local vault.zip bundle"
  type        = string
  default     = "/tmp/vault.zip"
}

variable "vault_build_date" {
  description = "The build date for Vault artifact"
  type        = string
  default     = ""
}

variable "vault_enable_file_audit_device" {
  description = "If true the file audit device will be enabled at the path /var/log/vault_audit.log"
  type        = bool
  default     = true
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
  default     = "/opt/vault/bin"
}

variable "vault_instance_count" {
  description = "How many instances to create for the Vault cluster"
  type        = number
  default     = 3
}

variable "vault_license_path" {
  description = "The path to a valid Vault enterprise edition license. This is only required for non-oss editions"
  type        = string
  default     = null
}

variable "vault_local_build_tags" {
  description = "The build tags to pass to the Go compiler for builder:local variants"
  type        = list(string)
  default     = null
}

variable "vault_log_level" {
  description = "The server log level for Vault logs. Supported values (in order of detail) are trace, debug, info, warn, and err."
  type        = string
  default     = "trace"
}

variable "vault_product_version" {
  description = "The version of Vault we are testing"
  type        = string
  default     = null
}

variable "vault_revision" {
  description = "The git sha of Vault artifact we are testing"
  type        = string
  default     = null
}

variable "vault_upgrade_initial_release" {
  description = "The Vault release to deploy before upgrading"
  default = {
    edition = "oss"
    // Vault 1.10.5 has a known issue with retry_join.
    version = "1.10.4"
  }
}
