# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0"
    }
  }
}

variable "network_name" {
  type        = string
  description = "The name of the Docker network to create"
  default     = "vault_cluster"
}

resource "docker_network" "cluster" {
  name = var.network_name
}

output "network_id" {
  value       = docker_network.cluster.id
  description = "The ID of the created Docker network"
}

output "network_name" {
  value       = docker_network.cluster.name
  description = "The name of the created Docker network"
}
