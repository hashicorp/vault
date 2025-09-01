# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "cluster_id" {
  description = "The ID of the HCP Vault cluster."
  type        = string
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
