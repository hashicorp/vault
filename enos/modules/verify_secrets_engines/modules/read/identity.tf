# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

// Read our testuser identity entity and verify that it matches our expected alias, groups, policy,
// and metadata.
resource "enos_remote_exec" "identity_verify_entity" {
  for_each = var.hosts

  environment = {
    ENTITY_ALIAS_ID = var.create_state.identity.data.entity_alias.id
    ENTITY_GROUP_IDS = jsonencode([
      var.create_state.kv.data.identity_group_kv_writers.id,
      var.create_state.identity.data.group_oidc_token_readers.id,
    ])
    ENTITY_METADATA   = jsonencode(var.create_state.identity.identity_entity_metadata)
    ENTITY_NAME       = var.create_state.identity.data.entity.name
    ENTITY_POLICIES   = jsonencode([var.create_state.auth.userpass.user.policy_name])
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = local.user_login_data.auth.client_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/identity-verify-entity.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}

// Read our OIDC key and role and verify that they have the correct configuration, TTLs, and algorithms.
resource "enos_remote_exec" "identity_verify_oidc" {
  for_each = var.hosts

  environment = {
    OIDC_ISSUER_URL           = var.create_state.identity.oidc.issuer_url
    OIDC_KEY_NAME             = var.create_state.identity.oidc.key_name
    OIDC_KEY_ROTATION_PERIOD  = var.create_state.identity.oidc.key_rotation_period
    OIDC_KEY_VERIFICATION_TTL = var.create_state.identity.oidc.key_verification_ttl
    OIDC_KEY_ALGORITHM        = var.create_state.identity.oidc.key_algorithm
    OIDC_ROLE_NAME            = var.create_state.identity.oidc.role_name
    OIDC_ROLE_TTL             = var.create_state.identity.oidc.role_ttl
    VAULT_ADDR                = var.vault_addr
    VAULT_TOKEN               = local.user_login_data.auth.client_token
    VAULT_INSTALL_DIR         = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/identity-verify-oidc.sh")]

  transport = {
    ssh = {
      host = each.value.public_ip
    }
  }
}
