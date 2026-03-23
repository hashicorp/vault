# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}

# Generate matrix.json for gotestsum from the test list
resource "local_file" "test_matrix" {
  filename = "/tmp/vault_test_matrix_${random_string.test_id.result}.json"
  content = jsonencode({
    include = length(var.test_names) > 0 ? [
      for test in var.test_names : {
        test = test
      }
    ] : []
  })
}

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "enos_local_exec" "run_blackbox_test" {
  scripts = [abspath("${path.module}/scripts/run-test.sh")]
  environment = merge({
    VAULT_TOKEN        = var.vault_root_token
    VAULT_ADDR         = var.vault_addr != null ? var.vault_addr : "http://${var.leader_public_ip}:8200"
    VAULT_TEST_PACKAGE = var.test_package
    VAULT_TEST_MATRIX  = length(var.test_names) > 0 ? local_file.test_matrix.filename : ""
    VAULT_EDITION      = var.vault_edition
    # PATH and Go-related environment variables are inherited from the calling process
    }, var.vault_namespace != null ? {
    VAULT_NAMESPACE = var.vault_namespace
    } : {}, local.ldap_environment
  )
  depends_on = [local_file.test_matrix]
}

# Local variables for LDAP environment setup
locals {
  # Extract LDAP configuration safely, defaulting to empty map if not available
  ldap_config = try(var.integration_host_state.ldap, {})

  # Convert domain (e.g., "enos.com") to DN format (e.g., "dc=enos,dc=com")
  domain_dn = try(local.ldap_config.domain, "") != "" ? join(",", [for part in split(".", local.ldap_config.domain) : "dc=${part}"]) : ""

  # Set up LDAP environment variables when LDAP integration is available
  # LDAP_SERVER uses private_ip for Vault operations (runs on Vault leader host)
  # LDAP_SERVER_PUBLIC uses public_ip for test setup operations (runs from GitHub runner)
  ldap_environment = try(local.ldap_config.domain, "") != "" ? {
    LDAP_SERVER        = local.ldap_config.host.private_ip
    LDAP_SERVER_PUBLIC = local.ldap_config.host.public_ip
    LDAP_PORT          = tostring(local.ldap_config.port)
    LDAP_BIND_DN       = "cn=admin,${local.domain_dn}"
    LDAP_BIND_PASS     = local.ldap_config.admin_pw
  } : {}
}

# Extract information from the script output
locals {
  json_file_path = try(
    regex("JSON_RESULTS_FILE=(.+)", enos_local_exec.run_blackbox_test.stdout)[0],
    ""
  )
  test_status = try(
    regex("TEST_STATUS=(.+)", enos_local_exec.run_blackbox_test.stdout)[0],
    "UNKNOWN"
  )
  test_exit_code = try(
    tonumber(regex("TEST_EXIT_CODE=(.+)", enos_local_exec.run_blackbox_test.stdout)[0]),
    null
  )
}
