# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

variable "cluster_name" {
  type        = string
  description = "The name of the cluster to load the image into"
}

variable "image" {
  type        = string
  description = "The image name for the image to load, i.e. hashicorp/vault"
}

variable "tag" {
  type        = string
  description = "The tag for the image to load, i.e. 1.12.0-dev"
}

variable "archive" {
  type        = string
  description = "The path to the image archive to load"
  default     = null
}

resource "enos_local_kind_load_image" "vault" {
  cluster_name = var.cluster_name
  image        = var.image
  tag          = var.tag
  archive      = var.archive
}

output "tag" {
  value       = var.tag
  description = "The tag of the docker image to load without the tag, i.e. 1.10.0"
}

output "image" {
  value       = var.image
  description = "The tag of the docker image to load without the tag, i.e. vault"
}

output "repository" {
  value       = enos_local_kind_load_image.vault.loaded_images.repository
  description = "The name of the image's repository, i.e. hashicorp/vault"
}
