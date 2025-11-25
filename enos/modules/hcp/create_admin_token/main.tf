# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

variable "cluster_id" {
  description = "The ID of the HCP Vault cluster."
  type        = string
}

# Docker compatibility variables (accepted but ignored)
variable "vault_root_token" {
  description = "Ignored - for Docker compatibility only"
  type        = string
  default     = null
  sensitive   = true
}

variable "vault_address" {
  description = "Ignored - for Docker compatibility only"
  type        = string
  default     = null
}

resource "hcp_vault_cluster_admin_token" "token" {
  cluster_id = var.cluster_id
}

output "created_at" {
  value = hcp_vault_cluster_admin_token.token.created_at
}

output "id" {
  value = hcp_vault_cluster_admin_token.token.id
}

output "token" {
  value = hcp_vault_cluster_admin_token.token.token
}
