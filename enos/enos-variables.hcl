variable "artifact_path" {
  type        = string
  description = "The local path for dev artifact to test"
  default     = null
}

variable "artifactory_username" {
  type        = string
  description = "The username to use when connecting to artifactory"
  default     = null
  sensitive   = true
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

variable "aws_region" {
  description = "The AWS region where we'll create infrastructure"
  type        = string
  default     = "us-west-1"
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

variable "backend_edition" {
  description = "The backend release edition if applicable"
  type        = string
  default     = "oss"
}

variable "backend_instance_type" {
  description = "The instance type to use for the Vault backend"
  type        = string
  default     = "t3.small"
}

variable "backend_license_path" {
  description = "The license for the backend if applicable (Consul Enterprise)"
  type        = string
  default     = null
}

variable "project_name" {
  description = "The description of the project"
  type        = string
  default     = "vault-enos-integration"
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
  description = "The Terraform Cloud QTI Organization API token."
  type        = string
  sensitive   = true
}

variable "vault_artifact_type" {
  description = "The Vault artifact type"
  default     = "bundle"
  validation {
    condition = var.vault_artifact_type == "package" || var.vault_artifact_type == "bundle"
  }
}

variable "vault_autopilot_initial_release" {
  description = "The Vault release to deploy before upgrading with autopilot"
  default = {
    edition = "ent"
    version = "1.11.0"
  }
}

variable "vault_bundle_path" {
  description = "Path to CRT generated or local vault.zip bundle"
  type        = string
  default     = "/tmp/vault.zip"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
  default     = "/opt/vault/bin"
}

variable "vault_instance_type" {
  description = "The instance type to use for the Vault backend"
  type        = string
  default     = null
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

variable "vault_build_date" {
  description = "The build date for Vault artifact"
  type        = string
  default     = ""
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
