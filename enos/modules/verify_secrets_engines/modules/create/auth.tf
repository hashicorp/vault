# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

locals {
  // Variables
  auth_userpass_path = "userpass"      # auth/userpass
  user_name          = "testuser"      # auth/userpass/users/testuser
  user_password      = "passtestuser1" # auth/userpass/login/passtestuser1
  user_policy_name   = "reguser"       # sys/policy/reguser

  auth_ldap_path = "ldap" # auth/ldap

  // Response data
  user_login_data = jsondecode(enos_remote_exec.auth_login_testuser.stdout)
  sys_auth_data   = jsondecode(enos_remote_exec.read_sys_auth.stdout).data

  // Output
  auth_output = {
    sys = local.sys_auth_data
    userpass = {
      path = local.auth_userpass_path
      user = {
        name        = local.user_name
        password    = local.user_password
        policy_name = local.user_policy_name
        login       = local.user_login_data
      }
    }
  }
}

output "auth" {
  value = local.auth_output
}

# Enable userpass auth
resource "enos_remote_exec" "auth_enable_userpass" {
  environment = {
    AUTH_METHOD       = "userpass"
    AUTH_PATH         = local.auth_userpass_path
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/auth-enable.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Get the sys/auth data after enabling our auth method
resource "enos_remote_exec" "read_sys_auth" {
  depends_on = [
    enos_remote_exec.auth_enable_userpass,
  ]
  environment = {
    REQPATH           = "sys/auth"
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/read.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Create a default policy for our users that allows them to read and list.
resource "enos_remote_exec" "policy_read_reguser" {
  environment = {
    POLICY_NAME       = local.user_policy_name
    POLICY_CONFIG     = <<-EOF
      path "*" {
        capabilities = ["read", "list"]
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

# Create our user
resource "enos_remote_exec" "auth_create_testuser" {
  depends_on = [
    enos_remote_exec.auth_enable_userpass,
    enos_remote_exec.policy_read_reguser,
  ]

  environment = {
    AUTH_PATH         = local.auth_userpass_path
    PASSWORD          = local.user_password
    POLICIES          = local.user_policy_name
    USERNAME          = local.user_name
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/auth-userpass-write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

resource "enos_remote_exec" "auth_login_testuser" {
  depends_on = [
    // Don't try to login until created our user and added it to the kv_writers group
    enos_remote_exec.auth_create_testuser,
    enos_remote_exec.identity_group_kv_writers,
  ]

  environment = {
    AUTH_PATH         = local.auth_userpass_path
    PASSWORD          = local.user_password
    USERNAME          = local.user_name
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/auth-userpass-login.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Enable ldap auth
resource "enos_remote_exec" "auth_enable_ldap" {
  environment = {
    AUTH_METHOD       = "ldap"
    AUTH_PATH         = local.auth_ldap_path
    VAULT_ADDR        = var.vault_addr
    VAULT_TOKEN       = var.vault_root_token
    VAULT_INSTALL_DIR = var.vault_install_dir
  }

  scripts = [abspath("${path.module}/../../scripts/auth-enable.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Write the initial ldap config
# This is a one time write to the leader node.
resource "enos_remote_exec" "auth_write_ldap_config" {
  depends_on = [
    enos_remote_exec.auth_enable_ldap
  ]

  environment = {
    AUTH_PATH         = local.auth_ldap_path
    GROUPATTR         = "memberOf"
    GROUPDN           = "CN=Users,DC=corp,DC=example,DC=net"
    INSECURE_TLS      = "true"
    POLICIES          = local.auth_ldap_path
    UPNDOMAIN         = "corp.example.net"
    URL               = "ldaps://ldap.example.com"
    USERATTR          = "sAMAccountName"
    USERDN            = "CN=Users,DC=corp,DC=example,DC=net"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../scripts/auth-ldap-write.sh")]

  transport = {
    ssh = {
      host = var.leader_host.public_ip
    }
  }
}

# Update the ldap config. Choose a random node each time to ensure that writes
# to all nodes are forwarded correctly and behave as we expect.
resource "random_integer" "auth_update_ldap_config_idx" {
  min = 0
  max = length(var.hosts) - 1
}

resource "enos_remote_exec" "auth_update_ldap_config" {
  depends_on = [
    enos_remote_exec.auth_write_ldap_config
  ]

  environment = {
    AUTH_PATH         = local.auth_ldap_path
    GROUPATTR         = "memberOf"
    GROUPDN           = "CN=Users,DC=corp,DC=example,DC=net"
    INSECURE_TLS      = "true"
    POLICIES          = local.auth_ldap_path
    UPNDOMAIN         = "corp.example.net"
    URL               = "ldaps://ldap2.example.com"
    USERATTR          = "sAMAccountName"
    USERDN            = "CN=Users,DC=corp,DC=example,DC=net"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../scripts/auth-ldap-write.sh")]

  transport = {
    ssh = {
      host = var.hosts[random_integer.auth_update_ldap_config_idx.result].public_ip
    }
  }
}
