// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

terraform_cli "default" {
  plugin_cache_dir = var.terraform_plugin_cache_dir != null ? abspath(var.terraform_plugin_cache_dir) : null
}

terraform_cli "dev" {
  plugin_cache_dir = var.terraform_plugin_cache_dir != null ? abspath(var.terraform_plugin_cache_dir) : null

  provider_installation {
    dev_overrides = {
      "registry.terraform.io/hashicorp-forge/enos" = try(abspath("../../terraform-provider-enos/dist"), null)
    }
    direct {}
  }
}

terraform "default" {
  required_version = ">= 1.2.0"

  required_providers {
    aws = {
      source = "hashicorp/aws"
    }

    docker = {
      source = "kreuzwerker/docker"
    }

    enos = {
      source  = "registry.terraform.io/hashicorp-forge/enos"
      version = ">= 0.4.0"
    }
  }
}
