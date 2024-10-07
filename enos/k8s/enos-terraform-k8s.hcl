# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform "k8s" {
  required_version = ">= 1.2.0"

  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }

    helm = {
      source = "hashicorp/helm"
    }
  }
}

terraform_cli "default" {
  plugin_cache_dir = var.terraform_plugin_cache_dir != null ? abspath(var.terraform_plugin_cache_dir) : null
}
