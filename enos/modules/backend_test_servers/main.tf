# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

locals {
  // Variables
  open_ldap_name = "openldap" # identity/group/name/kv_writers
}

# Enable kv secrets engine
resource "enos_remote_exec" "create_server_enviroment" {
  for_each = var.hosts

  environment = {
    HOST = "testing"
  }

  scripts = [abspath("${path.module}/scripts/test_server_setup.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}