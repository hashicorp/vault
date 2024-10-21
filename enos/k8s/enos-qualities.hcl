# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

quality "vault_artifact_container_alpine" {
  description = "The candidate binary packaged as an Alpine package is used for testing"
}

quality "vault_artifact_container_ubi" {
  description = "The candidate binary packaged as an UBI package is used for testing"
}

quality "vault_artifact_container_tags" {
  description = "The candidate binary has the expected tags"
}
