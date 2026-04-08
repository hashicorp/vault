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
locals {
  test_names = var.test_names != null ? var.test_names : []
}

resource "local_file" "test_matrix" {
  filename = "/tmp/vault_test_matrix_${random_string.test_id.result}.json"
  content = jsonencode({
    include = [
      for test in local.test_names : {
        test = test
      }
    ]
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
    VAULT_TEST_MATRIX  = length(local.test_names) > 0 ? local_file.test_matrix.filename : ""
    VAULT_EDITION      = var.vault_edition
    # PATH and Go-related environment variables are inherited from the calling process
    }, var.vault_namespace != null ? {
    VAULT_NAMESPACE = var.vault_namespace
    } : {}, local.ldap_environment, local.postgres_environment
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
  ldap_environment = try(local.ldap_config.domain, "") != "" ? {
    LDAP_URL_PRIVATE = "ldap://${local.ldap_config.host.private_ip}:${local.ldap_config.port}"
    LDAP_URL_PUBLIC  = "ldap://${local.ldap_config.host.public_ip}:${local.ldap_config.port}"
    LDAP_BIND_DN     = "cn=admin,${local.domain_dn}"
    LDAP_BIND_PASS   = local.ldap_config.admin_pw
  } : {}

  # Extract PostgreSQL configuration safely, defaulting to empty map if not available
  postgres_config = try(var.integration_host_state.postgres, {})

  # Set up PostgreSQL environment variables when PostgreSQL integration is available
  postgres_environment = try(local.postgres_config.host.private_ip, "") != "" ? {
    PG_URL            = "postgres://${local.postgres_config.username}:${local.postgres_config.password}@${local.postgres_config.host.private_ip}:${local.postgres_config.port}/${local.postgres_config.database}?sslmode=disable"
    POSTGRES_USER     = local.postgres_config.username
    POSTGRES_PASSWORD = local.postgres_config.password
    POSTGRES_DB       = local.postgres_config.database
    PGHOST            = local.postgres_config.host.private_ip
    PGPORT            = local.postgres_config.port
    PGUSER            = local.postgres_config.username
    PGPASSWORD        = local.postgres_config.password
    PGDATABASE        = local.postgres_config.database
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
