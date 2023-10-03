# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.4.3"
    }
  }
}

resource "random_pet" "cluster_name" {}

resource "enos_local_kind_cluster" "this" {
  name            = random_pet.cluster_name.id
  kubeconfig_path = var.kubeconfig_path
}

variable "kubeconfig_path" {
  type = string
}

output "cluster_name" {
  value = random_pet.cluster_name.id
}

output "kubeconfig_base64" {
  value = enos_local_kind_cluster.this.kubeconfig_base64
}

output "context_name" {
  value = enos_local_kind_cluster.this.context_name
}

output "host" {
  value = enos_local_kind_cluster.this.endpoint
}

output "client_certificate" {
  value = enos_local_kind_cluster.this.client_certificate
}

output "client_key" {
  value = enos_local_kind_cluster.this.client_key
}

output "cluster_ca_certificate" {
  value = enos_local_kind_cluster.this.cluster_ca_certificate
}
