terraform {
  required_providers {
    enos = {
      source = "app.terraform.io/hashicorp-qti/enos"
    }
  }
}

locals {
  # base test environment excludes the filter argument
  ui_test_environment_base = {
    VAULT_ADDR        = "http://${var.vault_addr}:8200"
    VAULT_TOKEN       = var.vault_root_token
    VAULT_UNSEAL_KEYS = jsonencode(slice(var.vault_unseal_keys, 0, var.vault_recovery_threshold))
  }
  ui_test_environment = var.ui_test_filter == null ? local.ui_test_environment_base : merge(local.ui_test_environment_base, {
    TEST_FILTER = var.ui_test_filter
  })
  # The environment variables need to be double escaped since the process of rendering them to the
  # outputs eats the escaping. Therefore double escaping ensures that the values are rendered as
  # properly escaped json, i.e. "[\"value\"]" suitable to be parsed as json.
  escaped_ui_test_environment = [
    for key, value in local.ui_test_environment : "export ${key}='${value}'"
  ]
}

resource "enos_local_exec" "test_ui" {
  count       = var.ui_run_tests ? 1 : 0
  environment = local.ui_test_environment
  scripts     = ["${path.module}/scripts/test_ui.sh"]
}
