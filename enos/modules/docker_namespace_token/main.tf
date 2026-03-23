# Copyright IBM Corp. 2016, 2025
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

# Create namespace using the root token
resource "enos_local_exec" "docker_create_namespace" {
  inline = [
    <<-EOT
      docker exec -e VAULT_ADDR=${var.vault_address} -e VAULT_TOKEN=${var.vault_root_token} \
        ${var.container_name} vault namespace create ${var.namespace_name}
    EOT
  ]
}

# Create policy at root level for blackbox testing (matches HVD admin namespace permissions)
resource "enos_local_exec" "docker_create_policy" {
  inline = [
    <<-EOT
      # Write policy to a temp file in the container
      docker exec ${var.container_name} sh -c 'cat > /tmp/${var.namespace_name}-policy.hcl << EOF
# HVD admin namespace compatible policy - restricted permissions to match cloud environment
path "*" {
  capabilities = ["sudo","read","create","update","delete","list","patch","subscribe"]
  subscribe_event_types = ["*"]
}
EOF'

      # Apply the policy at root level (not in a namespace)
      docker exec -e VAULT_ADDR=${var.vault_address} -e VAULT_TOKEN=${var.vault_root_token} \
        ${var.container_name} vault policy write ${var.namespace_name}-policy /tmp/${var.namespace_name}-policy.hcl
    EOT
  ]

  depends_on = [enos_local_exec.docker_create_namespace]
}

# Create token at root level with the policy that allows namespace operations
resource "enos_local_exec" "docker_create_token" {
  inline = [
    <<-EOT
      docker exec -e VAULT_ADDR=${var.vault_address} -e VAULT_TOKEN=${var.vault_root_token} \
        ${var.container_name} vault token create \
        -policy=${var.namespace_name}-policy \
        -ttl=24h \
        -renewable=true \
        -metadata="purpose=${var.namespace_name}-hvd-compatible-token" \
        -metadata="created_by=docker_namespace_token_module" \
        -format=json | jq -r '.auth.client_token'
    EOT
  ]

  depends_on = [enos_local_exec.docker_create_policy]
}

locals {
  # Use the created namespace token
  namespace_token = trimspace(enos_local_exec.docker_create_token.stdout)
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
  value       = var.namespace_name
  description = "The namespace where the token is valid"
}

output "policy" {
  value       = "${var.namespace_name}-policy"
  description = "The HVD-compatible policy assigned to the token (matches cloud environment permissions)"
}
