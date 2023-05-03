# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "enos" "default" {}

provider "helm" "default" {
  kubernetes {
    config_path = var.kubeconfig_path != null ? var.kubeconfig_path : abspath(joinpath(path.root, "kubeconfig"))
  }
}

provider "kubernetes" "default" {
  config_path = var.kubeconfig_path != null ? var.kubeconfig_path : abspath(joinpath(path.root, "kubeconfig"))
}

provider "azurerm" "default" {
  features {}

  skip_provider_registration = true
}

provider "azuread" "default" {}


