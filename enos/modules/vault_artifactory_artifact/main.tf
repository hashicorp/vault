# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source  = "registry.terraform.io/hashicorp-forge/enos"
      version = ">= 0.2.3"
    }
  }
}

data "enos_artifactory_item" "vault" {
  username = var.artifactory_username
  token    = var.artifactory_token
  name     = local.artifact_name
  host     = var.artifactory_host
  repo     = var.artifactory_repo
  path     = var.edition == "ce" ? "vault/*" : "vault-enterprise/*"
  properties = tomap({
    "commit"          = var.revision
    "product-name"    = var.edition == "ce" ? "vault" : "vault-enterprise"
    "product-version" = local.artifact_version
  })
}
