# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

variable "leader_host" {
  type = object({
    private_ip = string
    public_ip  = string
  })
  description = "The vault cluster host that is the leader"
}

variable "leader_public_ip" {
  type        = string
  description = "The public IP of the Vault leader"
}

variable "vault_root_token" {
  type        = string
  description = "The vault root token"
}

variable "test_names" {
  type        = list(string)
  description = "List of specific tests to run (e.g., ['TestStepdownAndLeaderElection', 'TestUnsealedStatus']). Empty list runs all tests."
  default     = []
}

variable "test_package" {
  type        = string
  description = "The Go package path for the tests (e.g., ./vault/external_tests/blackbox)"
}

variable "vault_addr" {
  type        = string
  description = "The full Vault address (for cloud environments). If provided, takes precedence over leader_public_ip."
  default     = null
}

variable "vault_namespace" {
  type        = string
  description = "The Vault namespace to operate in (for HCP environments). Optional."
  default     = null
}
