# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

terraform "kind" {
  required_version = ">= 1.2.0"
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }

    helm = {
      source = "hashicorp/helm"
    }
  }
}

terraform "azure" {
  required_version = ">= 1.2.0"
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
    helm = {
      source = "hashicorp/helm"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
    azuread = {
      source = "hashicorp/azuread"
    }
  }
}

terraform_cli "default" {
  plugin_cache_dir = var.terraform_plugin_cache_dir != null ? abspath(var.terraform_plugin_cache_dir) : null

  credentials "app.terraform.io" {
    token = var.tfc_api_token
  }
}
