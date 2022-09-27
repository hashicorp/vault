terraform {
  required_providers {
    enos = {
      source  = "app.terraform.io/hashicorp-qti/enos"
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
  path     = var.edition == "oss" ? "vault/*" : "vault-enterprise/*"
  properties = var.edition == "oss" ? tomap({
    "commit"          = var.revision
    "product-name"    = "vault"
    "product-version" = var.vault_product_version
    }) : tomap({
    "commit"          = var.revision
    "product-name"    = "vault-enterprise"
    "product-version" = "${var.vault_product_version}+ent"
  })
}
