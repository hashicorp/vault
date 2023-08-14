# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform "k8s" {
  required_version = ">= 1.2.0"

  required_providers {
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
      version = "< 0.4.0"
    }

    helm = {
      source = "hashicorp/helm"
    }
  }
}

terraform_cli "default" {
  plugin_cache_dir = var.terraform_plugin_cache_dir != null ? abspath(var.terraform_plugin_cache_dir) : null

  credentials "app.terraform.io" {
    token = var.tfc_api_token
  }
}
