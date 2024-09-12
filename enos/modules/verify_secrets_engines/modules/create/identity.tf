# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  identity_entity_metadata = {
    "organization" = "vault",
    "team"         = "qt",
  }
  group_name_oidc_readers     = "oidc_token_readers"            // identity/group/name/oidc_token_readers
  oidc_config_issuer_url      = "https://enos.example.com:1234" // identity/oidc/config
  oidc_key_algorithms         = ["RS256", "RS384", "RS512", "ES256", "ES384", "ES512", "EdDSA"]
  oidc_key_algorithm          = local.oidc_key_algorithms[random_integer.oidc_key_algorithm_idx.result]
  oidc_key_name               = "reguser" // identity/oidc/key/reguser
  oidc_key_rotation_period    = 86400     // 24h
  oidc_key_verification_ttl   = 21600     // 6h
  oidc_role_name              = "reguser" // identity/oidc/role/reguser
  oidc_role_ttl               = 3600      // 1h
  oidc_client_id              = "reguser" // optional client ID but required if we want to scope a key and role together without a *
  oidc_token_read_policy_name = "oidc_token_reader"

  // Response data
  oidc_token_data               = jsondecode(enos_remote_exec.oidc_token.stdout).data
  group_oidc_token_readers_data = jsondecode(enos_remote_exec.identity_group_oidc_token_readers.stdout).data
  initial_oidc_token_data       = jsondecode(enos_remote_exec.initial_oidc_token.stdout).data
  user_entity_data              = jsondecode(enos_remote_exec.identity_entity_testuser.stdout).data
  user_entity_alias_data        = jsondecode(enos_remote_exec.identity_entity_alias_testuser.stdout).data

  // Output
  identity_output = {
    oidc = {
      reader_group_name    = local.group_name_oidc_readers
      reader_policy_name   = local.oidc_token_read_policy_name
      issuer_url           = local.oidc_config_issuer_url
      key_algorithm        = local.oidc_key_algorithm
      key_name             = local.oidc_key_name
      key_rotation_period  = local.oidc_key_rotation_period
      key_verification_ttl = local.oidc_key_verification_ttl
      role_name            = local.oidc_role_name
      role_ttl             = local.oidc_role_ttl
      client_id            = local.oidc_client_id
    }
    identity_entity_metadata = local.identity_entity_metadata
    data = {
      entity                   = local.user_entity_data
      entity_alias             = local.user_entity_alias_data
      oidc_token               = local.oidc_token_data
      group_oidc_token_readers = local.group_oidc_token_readers_data
    }
  }
}

output "identity" {
  value = local.identity_output
}

// Get a random index for our algorithms so that we can randomly rotate through the various algorithms
resource "random_integer" "oidc_key_algorithm_idx" {
  min = 0
  max = length(local.oidc_key_algorithms) - 1
}

// Create identity entity for our user
resource "enos_remote_exec" "identity_entity_testuser" {
  depends_on = [
    enos_remote_exec.auth_create_testuser,
  ]

  environment = {
    REQPATH = "identity/entity"
    PAYLOAD = jsonencode({
      name     = local.user_name,
      metadata = local.identity_entity_metadata,
      policies = [local.user_policy_name],
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Create identity entity alias for our user
resource "enos_remote_exec" "identity_entity_alias_testuser" {
  environment = {
    REQPATH = "identity/entity-alias"
    PAYLOAD = jsonencode({
      name           = local.user_name,
      canonical_id   = local.user_entity_data.id
      mount_accessor = local.sys_auth_data["${local.auth_userpass_path}/"].accessor
      policies       = [local.user_policy_name],
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Configure our the oidc token backend
resource "enos_remote_exec" "oidc_config" {
  environment = {
    REQPATH = "identity/oidc/config"
    PAYLOAD = jsonencode({
      issuer = local.oidc_config_issuer_url,
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Create a named key that can sign OIDC identity token
resource "enos_remote_exec" "oidc_key" {
  environment = {
    REQPATH = "identity/oidc/key/${local.oidc_key_name}"
    PAYLOAD = jsonencode({
      allowed_client_ids = [local.oidc_client_id],
      algorithm          = local.oidc_key_algorithm,
      rotation_period    = local.oidc_key_rotation_period,
      verification_ttl   = local.oidc_key_verification_ttl,
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Create a role with custom template and that uses the named key
resource "enos_remote_exec" "oidc_role" {
  depends_on = [
    enos_remote_exec.oidc_key,
  ]

  environment = {
    REQPATH = "identity/oidc/role/${local.oidc_role_name}"
    PAYLOAD = jsonencode({
      client_id = local.oidc_client_id,
      key       = local.oidc_key_name,
      ttl       = local.oidc_role_ttl
      template = base64encode(<<-EOF
        {
          "team": {{identity.entity.metadata.team}},
          "organization": {{identity.entity.metadata.organization}},
          "groups": {{identity.entity.groups.names}}
        }
      EOF
      ),
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Create a group policy that allows "reading" a new signed OIDC token
resource "enos_remote_exec" "policy_write_oidc_token" {
  depends_on = [
    enos_remote_exec.secrets_enable_kv_secret,
  ]
  environment = {
    POLICY_NAME       = local.oidc_token_read_policy_name
    POLICY_CONFIG     = <<-EOF
      path "identity/oidc/token/*" {
        capabilities = ["read"]
      }
    EOF
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/policy-write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Create oidc_token_readers group and add our testuser to it
resource "enos_remote_exec" "identity_group_oidc_token_readers" {
  environment = {
    REQPATH = "identity/group"
    PAYLOAD = jsonencode({
      member_entity_ids = [local.user_entity_data.id],
      name              = local.group_name_oidc_readers,
      policies          = [local.oidc_token_read_policy_name],
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Generate a signed ID token with our test user
resource "enos_remote_exec" "initial_oidc_token" {
  depends_on = [
    enos_remote_exec.oidc_role,
  ]

  environment = {
    REQPATH           = "identity/oidc/token/${local.oidc_role_name}"
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = local.user_login_data.auth.client_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/read.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Introspect the signed ID and verify it
resource "enos_remote_exec" "oidc_introspect_initial_token" {
  environment = {
    ASSERT_ACTIVE = true // Our token should be "active"
    PAYLOAD = jsonencode({
      token     = local.initial_oidc_token_data.token,
      client_id = local.initial_oidc_token_data.client_id
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/identity-oidc-introspect-token.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Rotate the key with a zero TTL to force expiration
resource "enos_remote_exec" "oidc_key_rotate" {
  depends_on = [
    enos_remote_exec.oidc_introspect_initial_token,
  ]

  environment = {
    REQPATH = "identity/oidc/key/${local.oidc_key_name}/rotate"
    PAYLOAD = jsonencode({
      verification_ttl = 0,
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/write-payload.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Introspect it again to make sure it's no longer active
resource "enos_remote_exec" "oidc_introspect_initial_token_post_rotate" {
  depends_on = [
    enos_remote_exec.oidc_key_rotate,
  ]

  environment = {
    ASSERT_ACTIVE = false // Our token should not be "active"
    PAYLOAD = jsonencode({
      token     = local.initial_oidc_token_data.token,
      client_id = local.initial_oidc_token_data.client_id
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/identity-oidc-introspect-token.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Generate a new token that we can use later
resource "enos_remote_exec" "oidc_token" {
  depends_on = [
    enos_remote_exec.oidc_introspect_initial_token_post_rotate,
  ]

  environment = {
    REQPATH           = "identity/oidc/token/${local.oidc_role_name}"
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = local.user_login_data.auth.client_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/read.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

// Introspect the new token to ensure it's active before we export it for user later via outputs
resource "enos_remote_exec" "oidc_introspect_token" {
  environment = {
    ASSERT_ACTIVE = true // Our token should be "active"
    PAYLOAD = jsonencode({
      token     = local.oidc_token_data.token,
      client_id = local.oidc_token_data.client_id
    })
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/identity-oidc-introspect-token.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}
