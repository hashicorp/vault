# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

variable "vault_root_token" {
  description = "The root token from the Docker Vault cluster"
  type        = string
  sensitive   = true
  default     = null
}

variable "vault_address" {
  description = "The address of the Vault cluster"
  type        = string
  default     = null
}

# HCP compatibility variables (accepted but ignored)
variable "cluster_id" {
  description = "Ignored - for HCP compatibility only"
  type        = string
  default     = null
}

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

# We need the container name to exec into
variable "container_name" {
  description = "The name of the Docker container running Vault (Docker only)"
  type        = string
  default     = null
}

variable "namespace_name" {
  description = "The name of the namespace to create and generate the token in"
  type        = string
  default     = "admin"
}

# Create namespace using the root token (only when all required vars are present)
resource "enos_local_exec" "docker_create_namespace" {
  count = var.vault_address != null && var.vault_root_token != null && var.container_name != null ? 1 : 0

  inline = [
    <<-EOT
      docker exec -e VAULT_ADDR=${var.vault_address} -e VAULT_TOKEN=${var.vault_root_token} \
        ${var.container_name} vault namespace create ${var.namespace_name}
    EOT
  ]
}

# Create policy in the namespace
resource "enos_local_exec" "docker_create_policy" {
  count = var.vault_address != null && var.vault_root_token != null && var.container_name != null ? 1 : 0

  inline = [
    <<-EOT
      # Write policy to a temp file in the container
      docker exec ${var.container_name} sh -c 'cat > /tmp/${var.namespace_name}-policy.hcl << EOF
path "*" {
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
EOF'

      # Apply the policy in the namespace
      docker exec -e VAULT_ADDR=${var.vault_address} -e VAULT_TOKEN=${var.vault_root_token} -e VAULT_NAMESPACE=${var.namespace_name} \
        ${var.container_name} vault policy write ${var.namespace_name}-policy /tmp/${var.namespace_name}-policy.hcl
    EOT
  ]

  depends_on = [enos_local_exec.docker_create_namespace]
}

# Create token in the namespace
resource "enos_local_exec" "docker_create_token" {
  count = var.vault_address != null && var.vault_root_token != null && var.container_name != null ? 1 : 0

  inline = [
    <<-EOT
      docker exec -e VAULT_ADDR=${var.vault_address} -e VAULT_TOKEN=${var.vault_root_token} -e VAULT_NAMESPACE=${var.namespace_name} \
        ${var.container_name} vault token create \
        -policy=${var.namespace_name}-policy \
        -ttl=24h \
        -renewable=true \
        -metadata="purpose=${var.namespace_name}-token" \
        -metadata="created_by=docker_namespace_token_module" \
        -format=json | jq -r '.auth.client_token'
    EOT
  ]

  depends_on = [enos_local_exec.docker_create_policy]
}

locals {
  # For Docker: use the created namespace token, for HCP: use root token (fallback)
  namespace_token = length(enos_local_exec.docker_create_token) > 0 ? trimspace(enos_local_exec.docker_create_token[0].stdout) : var.vault_root_token
}

output "created_at" {
  value = timestamp()
}

output "id" {
  value = "docker-${var.namespace_name}-token"
}

output "token" {
  value     = local.namespace_token
  sensitive = true
}

output "namespace" {
  value       = length(enos_local_exec.docker_create_token) > 0 ? var.namespace_name : "root"
  description = "The namespace where the token is valid"
}

output "policy" {
  value       = length(enos_local_exec.docker_create_token) > 0 ? "${var.namespace_name}-policy" : "root"
  description = "The policy assigned to the token"
}
