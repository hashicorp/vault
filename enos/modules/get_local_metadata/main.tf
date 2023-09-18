# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

resource "enos_local_exec" "get_build_date" {
  scripts = [abspath("${path.module}/scripts/build_date.sh")]
}

resource "enos_local_exec" "get_revision" {
  inline = ["git rev-parse HEAD"]
}

resource "enos_local_exec" "get_version" {
  inline = ["${abspath("${path.module}/scripts/version.sh")} version"]
}

resource "enos_local_exec" "get_version_base" {
  inline = ["${abspath("${path.module}/scripts/version.sh")} version-base"]
}

resource "enos_local_exec" "get_version_pre" {
  inline = ["${abspath("${path.module}/scripts/version.sh")} version-pre"]
}

resource "enos_local_exec" "get_version_meta" {
  inline = ["${abspath("${path.module}/scripts/version.sh")} version-meta"]
}

output "build_date" {
  value = trimspace(enos_local_exec.get_build_date.stdout)
}

output "revision" {
  value = trimspace(enos_local_exec.get_revision.stdout)
}

output "version" {
  value = trimspace(enos_local_exec.get_version.stdout)
}

output "version_base" {
  value = trimspace(enos_local_exec.get_version_base.stdout)
}

output "version_pre" {
  value = trimspace(enos_local_exec.get_version_pre.stdout)
}

output "version_meta" {
  value = trimspace(enos_local_exec.get_version_meta.stdout)
}
