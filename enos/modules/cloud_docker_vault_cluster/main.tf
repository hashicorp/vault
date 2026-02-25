# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

variable "min_vault_version" {
  type        = string
  description = "The minimum Vault version to deploy (e.g., 1.15.0 or v1.15.0+ent)"
}

variable "vault_edition" {
  type        = string
  description = "The edition of Vault to deploy (ent, ce, ent.fips1403)"
  default     = "ent"

  validation {
    condition     = contains(["ent", "ce", "ent.fips1403"], var.vault_edition)
    error_message = "vault_edition must be one of: ent, ce, ent.fips1403"
  }
}

variable "vault_license" {
  type        = string
  description = "The Vault Enterprise license"
  default     = null
  sensitive   = true
}

variable "cluster_name" {
  type        = string
  description = "The name of the Vault cluster"
  default     = "vault"
}

variable "container_count" {
  type        = number
  description = "Number of Vault containers to create"
  default     = 3
}

variable "vault_port" {
  type        = number
  description = "The port Vault listens on"
  default     = 8200
}

variable "use_local_build" {
  type        = bool
  description = "If true, build a local Docker image from the current branch instead of pulling from Docker Hub"
  default     = false
}

# HCP-specific variables (ignored but accepted for compatibility)
variable "network_name" {
  type        = string
  description = "Ignored - for HCP compatibility only"
  default     = ""
}

variable "tier" {
  type        = string
  description = "Ignored - for HCP compatibility only"
  default     = ""
}

# Generate a random suffix for the network name to avoid conflicts
resource "random_string" "network_suffix" {
  length  = 8
  lower   = true
  upper   = false
  numeric = true
  special = false
}

# Create Docker network
resource "docker_network" "cluster" {
  name = "${var.cluster_name}-network-${random_string.network_suffix.result}"
}

locals {
  # Parse min_vault_version to extract the version number
  # e.g., "v1.15.0+ent" -> "1.15.0" or "v1.15.0+ent-2cf0b2f" -> "1.15.0"
  vault_version = trimprefix(split("+", var.min_vault_version)[0], "v")

  image_map = {
    "ent"          = "hashicorp/vault-enterprise"
    "ce"           = "hashicorp/vault"
    "ent.fips1403" = "hashicorp/vault-enterprise-fips"
  }
  target_map = {
    "ent"          = "ubi"
    "ce"           = "ubi"
    "ent.fips1403" = "ubi-fips"
  }
  image      = local.image_map[var.vault_edition]
  tag_suffix = var.vault_edition == "ce" ? "" : "-ent"
  image_tag  = "${local.vault_version}${local.tag_suffix}"
  local_tag  = "vault-local-${var.vault_edition}:${local.vault_version}"
  dockerfile = "Dockerfile"
  target     = local.target_map[var.vault_edition]
}

# Pull image from Docker Hub (when not using local build)
resource "docker_image" "vault_remote" {
  count = var.use_local_build ? 0 : 1
  name  = "${local.image}:${local.image_tag}"
}

# Build image from local Dockerfile (when using local build)
resource "docker_image" "vault_local" {
  count        = var.use_local_build ? 1 : 0
  name         = local.local_tag
  keep_locally = true

  build {
    context     = "${path.module}/../../.."
    dockerfile  = local.dockerfile
    target      = local.target
    tag         = [local.local_tag]
    pull_parent = true
    build_args = {
      BIN_NAME         = "vault"
      TARGETOS         = "linux"
      TARGETARCH       = "amd64"
      NAME             = "vault"
      PRODUCT_VERSION  = local.vault_version
      PRODUCT_REVISION = "local"
      LICENSE_SOURCE   = "LICENSE"
      LICENSE_DEST     = "/usr/share/doc/vault"
    }
  }

}

locals {
  # Generate Vault configuration for each node
  vault_config_template = <<-EOF
    ui = true
    listener "tcp" {
      address = "0.0.0.0:${var.vault_port}"
      cluster_address = "0.0.0.0:8201"
      tls_disable = true
    }

    storage "raft" {
      path = "/vault/data"
      node_id = "node%s"
    }

    disable_mlock = true
  EOF
}

# Using tmpfs for Raft data (in-memory, no persistence needed for testing)

resource "docker_container" "vault" {
  count = var.container_count
  name  = "${var.cluster_name}-${count.index}"
  image = var.use_local_build ? docker_image.vault_local[0].name : docker_image.vault_remote[0].image_id

  networks_advanced {
    name = docker_network.cluster.name
  }

  ports {
    internal = var.vault_port
    external = var.vault_port + count.index
  }

  tmpfs = {
    "/vault/data" = "rw,noexec,nosuid,size=100m"
  }

  upload {
    content = format(local.vault_config_template, count.index)
    file    = "/vault/config/vault.hcl"
  }


  user = "root"

  env = concat(
    [
      "VAULT_API_ADDR=http://${var.cluster_name}-${count.index}:${var.vault_port}",
      "VAULT_CLUSTER_ADDR=http://${var.cluster_name}-${count.index}:8201",
      "SKIP_SETCAP=true",
      "SKIP_CHOWN=true",
    ],
    var.vault_license != null ? ["VAULT_LICENSE=${var.vault_license}"] : []
  )

  capabilities {
    add = ["IPC_LOCK"]
  }

  command = ["vault", "server", "-config=/vault/config/vault.hcl"]

  restart = "no"
}

locals {
  instance_indexes = [for idx in range(var.container_count) : tostring(idx)]
  leader_idx       = 0
  followers_idx    = range(1, var.container_count)

  vault_address   = "http://127.0.0.1:${var.vault_port}"
  leader_api_addr = "http://${var.cluster_name}-${local.leader_idx}:${var.vault_port}"
}

# Initialize Vault on the leader
resource "enos_local_exec" "init_leader" {
  inline = [
    <<-EOT
      # Wait for Vault to be ready (output to stderr to keep stdout clean)
      for i in 1 2 3 4 5 6 7 8 9 10; do
        if docker exec -e VAULT_ADDR=http://127.0.0.1:${var.vault_port} ${docker_container.vault[local.leader_idx].name} vault status 2>&1 | grep -q "Initialized.*false"; then
          break
        fi
        echo "Waiting for Vault to start (attempt $i/10)..." >&2
        sleep 2
      done

      # Initialize Vault and output JSON to stdout
      docker exec -e VAULT_ADDR=http://127.0.0.1:${var.vault_port} ${docker_container.vault[local.leader_idx].name} vault operator init \
        -key-shares=1 \
        -key-threshold=1 \
        -format=json
    EOT
  ]

  depends_on = [docker_container.vault]
}

locals {
  init_data  = jsondecode(enos_local_exec.init_leader.stdout)
  unseal_key = local.init_data.unseal_keys_b64[0]
  root_token = local.init_data.root_token
}

# Unseal the leader
resource "enos_local_exec" "unseal_leader" {
  inline = [
    "docker exec -e VAULT_ADDR=http://127.0.0.1:${var.vault_port} ${docker_container.vault[local.leader_idx].name} vault operator unseal ${local.unseal_key}"
  ]

  depends_on = [enos_local_exec.init_leader]
}

# Join followers to Raft cluster and unseal them
resource "enos_local_exec" "join_followers" {
  count = length(local.followers_idx)

  inline = [
    <<-EOT
      # Wait for Vault to be ready
      for i in 1 2 3 4 5; do
        docker exec -e VAULT_ADDR=http://127.0.0.1:${var.vault_port} ${docker_container.vault[local.followers_idx[count.index]].name} vault status > /dev/null 2>&1 && break || sleep 5
      done

      # Join the Raft cluster
      docker exec -e VAULT_ADDR=http://127.0.0.1:${var.vault_port} ${docker_container.vault[local.followers_idx[count.index]].name} \
        vault operator raft join ${local.leader_api_addr}

      # Unseal the follower
      docker exec -e VAULT_ADDR=http://127.0.0.1:${var.vault_port} ${docker_container.vault[local.followers_idx[count.index]].name} \
        vault operator unseal ${local.unseal_key}
    EOT
  ]

  depends_on = [enos_local_exec.unseal_leader]
}

# Outputs that match HCP module interface
output "cloud_provider" {
  value       = "docker"
  description = "The cloud provider (docker for local)"
}

output "cluster_id" {
  value       = var.cluster_name
  description = "The cluster identifier"
}

output "created_at" {
  value       = timestamp()
  description = "Timestamp of cluster creation"
}

output "id" {
  value       = var.cluster_name
  description = "The cluster identifier"
}

output "namespace" {
  value       = "root"
  description = "The Vault namespace"
}

output "organization_id" {
  value       = "docker-local"
  description = "The organization identifier"
}

output "region" {
  value       = "local"
  description = "The region or location"
}

output "self_link" {
  value       = ""
  description = "Self link to the cluster"
}

output "state" {
  value       = "RUNNING"
  description = "The state of the cluster"
}

output "vault_private_endpoint_url" {
  value       = ""
  description = "Private endpoint URL (not applicable for Docker)"
}

output "vault_proxy_endpoint_url" {
  value       = ""
  description = "Proxy endpoint URL (not applicable for Docker)"
}

output "vault_public_endpoint_url" {
  value       = "http://localhost:${var.vault_port}"
  description = "Public endpoint URL"
}

output "vault_version" {
  value       = local.vault_version
  description = "The version of Vault deployed"
}

# Docker-specific outputs
output "container_names" {
  value       = docker_container.vault[*].name
  description = "The names of the Vault containers"
}

output "container_ids" {
  value       = docker_container.vault[*].id
  description = "The IDs of the Vault containers"
}

output "vault_addresses" {
  value = [
    for i in range(var.container_count) :
    "http://localhost:${var.vault_port + i}"
  ]
  description = "The addresses of the Vault containers"
}

output "primary_address" {
  value       = "http://localhost:${var.vault_port}"
  description = "The address of the primary Vault container"
}

output "network_id" {
  value       = docker_network.cluster.id
  description = "The ID of the created Docker network"
}

output "network_name" {
  value       = docker_network.cluster.name
  description = "The name of the created Docker network"
}

output "image_name" {
  value       = var.use_local_build ? (length(docker_image.vault_local) > 0 ? docker_image.vault_local[0].name : "none") : (length(docker_image.vault_remote) > 0 ? docker_image.vault_remote[0].name : "none")
  description = "The Docker image being used"
}

output "is_local_build" {
  value       = var.use_local_build
  description = "Whether this is using a local build"
}

output "vault_root_token" {
  value       = local.root_token
  sensitive   = true
  description = "The root token for the Vault cluster"
}
