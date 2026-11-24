# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

# ---------------------------------------------------------------------------
# kubeclient_bin_path — path to the kubeclient binary INSIDE the
# vault-client-sidecar container image. This is the only variable you change
# when the standalone client image ships.
#
# TODAY  : default "/bin/vault"
#   The vault-enterprise image is used as the client image. The vault binary
#   does not speak the kubeclient CLI flags, so the smoke module detects the
#   fallback value and skips all binary-level tests automatically. RBAC is
#   still applied so the infrastructure is ready.
#
# FUTURE : set to the path of the kubeclient binary in the standalone image,
#   e.g. "/bin/kubeclient". The smoke module will exec into the
#   vault-client-sidecar container and run the full get-pod / patch-pod /
#   invalid-call suite automatically. No other changes required.
# ---------------------------------------------------------------------------
variable "kubeclient_bin_path" {
  description = "Path to the kubeclient binary inside the vault-client-sidecar container. Defaults to /bin/vault (skips binary smoke tests). Set to the real path when the standalone client image is available."
  type        = string
  default     = "/bin/vault"
}

variable "client_container_image_archive" {
  description = "The path to the location of the client container image archive to test as a sidecar"
  type        = string
  default     = null # If none is given we'll load client image or skip if not provided
}

variable "container_image_archive" {
  description = "The path to the location of the container image archive to test"
  type        = string
  default     = null # If none is given we'll simply load a container from a repo
}

variable "log_level" {
  description = "The server log level for Vault logs. Supported values (in order of detail) are trace, debug, info, warn, and err."
  type        = string
  default     = "trace"
}

variable "instance_count" {
  description = "How many instances to create for the Vault cluster"
  type        = number
  default     = 3
}

variable "terraform_plugin_cache_dir" {
  description = "The directory to cache Terraform modules and providers"
  type        = string
  default     = null
}

variable "vault_build_date" {
  description = "The expected vault build date"
  type        = string
  default     = ""
}

variable "vault_revision" {
  type        = string
  description = "The expected vault revision"
  default     = "ce"
}

variable "vault_version" {
  description = "The expected vault version"
  type        = string
  default     = "1.18.0"
}
