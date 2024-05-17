# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


output "url" {
  value       = data.enos_artifactory_item.vault.results[0].url
  description = "The artifactory download url for the artifact"
}

output "sha256" {
  value       = data.enos_artifactory_item.vault.results[0].sha256
  description = "The sha256 checksum for the artifact"
}

output "size" {
  value       = data.enos_artifactory_item.vault.results[0].size
  description = "The size in bytes of the artifact"
}

output "name" {
  value       = data.enos_artifactory_item.vault.results[0].name
  description = "The name of the artifact"
}

output "vault_artifactory_release" {
  value = {
    url      = data.enos_artifactory_item.vault.results[0].url
    sha256   = data.enos_artifactory_item.vault.results[0].sha256
    username = var.artifactory_username
    token    = var.artifactory_token
  }
}
