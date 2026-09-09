# Copyright IBM Corp. 2016, 2025
# SPDX-License-Identifier: BUSL-1.1

terraform {
  required_providers {
    enos = {
      source = "registry.terraform.io/hashicorp-forge/enos"
    }
  }
}
locals {
  dynamic_role_name = "dynamic-role"
}

variable "hosts" {
  type = map(object({
    ipv6       = string
    private_ip = string
    public_ip  = string
  }))
  description = "The Vault cluster instances that were created"
}

variable "create_state" {
  description = "The state of the secrets engines from the 'create' module"
}

variable "vault_addr" {
  type        = string
  description = "The local vault API listen address"
}

variable "vault_install_dir" {
  type        = string
  description = "The directory where the Vault binary will be installed"
}

variable "vault_root_token" {
  type        = string
  description = "The Vault root token"
  default     = null
}

variable "credential_ttl_buffer" {
  description = "Buffer (seconds) to wait after LDAP credential TTL expiry"
  type        = number
  default     = 80
}

variable "default_ttl" {
  description = "Default time-to-live (in seconds) for issued LDAP credentials."
  type        = number
  default     = 60
}

variable "max_ttl" {
  description = "Maximum time-to-live (in seconds) allowed for issued LDAP credentials"
  type        = number
  default     = 60
}

variable "enable_secrets_verification" {
  type        = bool
  description = "Enable LDAP secrets engine verification (dynamic credentials)"
  default     = true
}

variable "enable_rollback_verification" {
  type        = bool
  description = "Enable LDAP secrets engine rollback verification"
  default     = true
}

variable "enable_dynamic_credentials_verification" {
  type        = bool
  description = "Enable comprehensive LDAP dynamic credentials verification tests"
  default     = true
}

variable "enable_rotation_verification" {
  type        = bool
  description = "Enable LDAP root rotation verification"
  default     = true
}

variable "enable_auth_verification" {
  type        = bool
  description = "Enable LDAP authentication verification"
  default     = true
}


variable "enable_dynamic_role_verification" {
  type        = bool
  description = "Enable LDAP secrets engine dynamic role"
  default     = true
}

variable "vault_audit_log_path" {
  type        = string
  description = "The file path for the audit device"
}

variable "enable_password_policy_verification" {
  type        = bool
  description = "Enable LDAP authentication verification"
  default     = true
}

locals {
  strong_password_policy = "strong-policy"
}

variable "enable_static_role_verification" {
  type        = bool
  description = "Enable LDAP secrets engine static role verification"
  default     = true
}

locals {
  lease_checks_count = try(var.create_state.ldap.data.checkout_custom.lease_id, "") != "" ? 1 : 0
}

resource "enos_remote_exec" "ldap_verify_auth" {
  count = var.enable_auth_verification ? 1 : 0
  environment = {
    MOUNT             = "${var.create_state.ldap.ldap_mount}"
    LDAP_SERVER       = "${var.create_state.ldap.host.private_ip}"
    LDAP_PORT         = "${var.create_state.ldap.port}"
    LDAP_USERNAME     = "${var.create_state.ldap.username}"
    LDAP_ADMIN_PW     = "${var.create_state.ldap.pw}"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }
  scripts = [abspath("${path.module}/../../../scripts/ldap/verify-auth.sh")]
  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Configure and verify LDAP secrets engine 
resource "enos_remote_exec" "ldap_verify_secrets" {
  count = var.enable_secrets_verification ? 1 : 0

  environment = {
    MOUNT                 = var.create_state.ldap.ldap_mount
    LDAP_SERVER           = var.create_state.ldap.host.private_ip
    LDAP_PORT             = var.create_state.ldap.port
    LDAP_USERNAME         = var.create_state.ldap.username
    LDAP_ADMIN_PW         = var.create_state.ldap.pw
    VAULT_ADDR            = var.vault_addr
    VAULT_INSTALL_DIR     = var.vault_install_dir
    VAULT_TOKEN           = var.vault_root_token
    CREDENTIAL_TTL_BUFFER = tostring(var.credential_ttl_buffer)
    DEFAULT_TTL           = tostring(var.default_ttl)
    MAX_TTL               = tostring(var.max_ttl)
    STRONG_POLICY         = local.strong_password_policy
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap/verify-secrets.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Verify LDAP root rotation 
resource "enos_remote_exec" "ldap_verify_rotation" {
  count = var.enable_rotation_verification ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_secrets
  ]
  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    LDAP_SERVER       = var.create_state.ldap.host.private_ip
    LDAP_PORT         = var.create_state.ldap.port
    LDAP_USERNAME     = var.create_state.ldap.username
    LDAP_ADMIN_PW     = var.create_state.ldap.pw
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }
  scripts = [abspath("${path.module}/../../../scripts/ldap/verify-rotation.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Configure and verify LDAP secrets engine rollback behavior
resource "enos_remote_exec" "ldap_verify_rollback" {
  count = var.enable_rollback_verification ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
    enos_remote_exec.ldap_verify_rotation
  ]

  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    LDAP_SERVER       = var.create_state.ldap.host.private_ip
    LDAP_PORT         = var.create_state.ldap.port
    LDAP_USERNAME     = var.create_state.ldap.username
    LDAP_ADMIN_PW     = var.create_state.ldap.pw
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [
    abspath("${path.module}/../../../scripts/ldap/secrets-rollback/secrets-rollback-invalid-config.sh"),
    abspath("${path.module}/../../../scripts/ldap/secrets-rollback/secrets-rollback-creds-mismatch.sh"),
    abspath("${path.module}/../../../scripts/ldap/secrets-rollback/secrets-rollback-transactional.sh"),
  ]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# LDAP dynamic credentials verification tests - run together in sequence
resource "enos_remote_exec" "ldap_verify_dynamic_credentials_suite" {
  count = var.enable_dynamic_credentials_verification ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
    enos_remote_exec.ldap_verify_rotation
  ]

  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    LDAP_SERVER       = var.create_state.ldap.host.private_ip
    LDAP_PORT         = var.create_state.ldap.port
    LDAP_USERNAME     = var.create_state.ldap.username
    LDAP_ADMIN_PW     = var.create_state.ldap.pw
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
    DEFAULT_TTL       = tostring(var.default_ttl)
    MAX_TTL           = tostring(var.max_ttl)
  }

  scripts = [
    abspath("${path.module}/../../../scripts/ldap/vault_dynamic_credentials/verify-ttl-limits.sh"),
    abspath("${path.module}/../../../scripts/ldap/vault_dynamic_credentials/verify-credential-renewal.sh"),
    abspath("${path.module}/../../../scripts/ldap/vault_dynamic_credentials/verify-manual-revocation.sh"),
    abspath("${path.module}/../../../scripts/ldap/vault_dynamic_credentials/verify-rollback.sh"),
    abspath("${path.module}/../../../scripts/ldap/vault_dynamic_credentials/verify-auto-cleanup.sh"),
    abspath("${path.module}/../../../scripts/ldap/vault_dynamic_credentials/verify-error-handling.sh")
  ]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Read Library configuration
# Test Case: Read Library configuration - Read the library set details
resource "enos_remote_exec" "ldap_library_set_read" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
  ]

  environment = {
    REQPATH           = "${var.create_state.ldap.ldap_mount}/library/test-set"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/read.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# List all library sets
# Test Case #5: List all library sets - List all the service account library sets
resource "enos_remote_exec" "ldap_library_list_all" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
  ]

  environment = {
    REQPATH           = "${var.create_state.ldap.ldap_mount}/library"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/list.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# List library set by name
# Test Case #6: List library sets by account name - List account details for the given service account set
resource "enos_remote_exec" "ldap_library_list_set" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
  ]

  environment = {
    REQPATH           = "${var.create_state.ldap.ldap_mount}/library/test-set"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/list.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# List library sets by account name
# Test Case #7: List library sets by account name - List account details for the given service account
resource "enos_remote_exec" "ldap_library_list_by_account" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
  ]

  environment = {
    # Using the service account name from test case #1 (uid=fizz)
    REQPATH           = "${var.create_state.ldap.ldap_mount}/library/fizz"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/list.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

resource "enos_remote_exec" "ldap_verify_audit_trail" {
  count = var.enable_dynamic_credentials_verification && var.vault_audit_log_path != null && var.vault_audit_log_path != "" ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_secrets
  ]

  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    LDAP_SERVER       = var.create_state.ldap.host.private_ip
    LDAP_PORT         = var.create_state.ldap.port
    LDAP_USERNAME     = var.create_state.ldap.username
    LDAP_ADMIN_PW     = var.create_state.ldap.pw
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
    VAULT_AUDIT_LOG   = var.vault_audit_log_path
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap/vault_dynamic_credentials/verify-audit-trail.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Renew Check-out Lease
# Test Case #10: Renew Check-out Lease - Renew the lease for a checked-out account
resource "enos_remote_exec" "ldap_library_checkout_lease_renew" {
  count = local.lease_checks_count

  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
  ]

  environment = {
    # LEASE_ID will be provided via create_state.ldap from the create module after checkout
    LEASE_ID          = var.create_state.ldap.data.checkout_custom.lease_id
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap-lease-renew.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}
# Configure and verify LDAP secrets engine rollback behavior
resource "enos_remote_exec" "verify_dynamic_role" {
  count = var.enable_dynamic_role_verification ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
    enos_remote_exec.ldap_verify_rotation
  ]
  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    LDAP_SERVER       = var.create_state.ldap.host.private_ip
    LDAP_PORT         = var.create_state.ldap.port
    LDAP_USERNAME     = var.create_state.ldap.username
    LDAP_ADMIN_PW     = var.create_state.ldap.pw
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
    ROLE_NAME         = local.dynamic_role_name
    DEFAULT_TTL       = tostring(var.default_ttl)
    MAX_TTL           = tostring(var.max_ttl)
  }

  scripts = [
    abspath("${path.module}/../../../scripts/ldap/Dynamic-roles/dynamic-roles.sh"),
    abspath("${path.module}/../../../scripts/ldap/Dynamic-roles/dynamic-roles-validation.sh"),
    abspath("${path.module}/../../../scripts/ldap/Dynamic-roles/dynamic-roles-listing.sh"),
    abspath("${path.module}/../../../scripts/ldap/Dynamic-roles/dynamic-roles-audit.sh"),
    abspath("${path.module}/../../../scripts/ldap/Dynamic-roles/dynamic-roles-rollback.sh"),
    abspath("${path.module}/../../../scripts/ldap/Dynamic-roles/dynamic-roles-deletion.sh")
  ]
  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}


# Self Check-in (Automatic on Revoke)
# Test Case #12: Self Check-in (Automatic on Revoke) - Return account when lease expires (revoke)
resource "enos_remote_exec" "ldap_library_checkout_lease_revoke" {
  count = local.lease_checks_count

  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
  ]

  environment = {
    # LEASE_ID will be provided via create_state.ldap from the create module after checkout
    LEASE_ID          = var.create_state.ldap.data.checkout_custom.lease_id
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap-lease-revoke.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Check Library Status
# Test Case #14: Check Library Status - See which accounts are available/checked-out
resource "enos_remote_exec" "ldap_library_status_read" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
  ]

  environment = {
    REQPATH           = "${var.create_state.ldap.ldap_mount}/library/test-set/status"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/read.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# View Check-out Details
# Test Case #15: View Check-out Details - Track which accounts are available/checked out
resource "enos_remote_exec" "ldap_library_checkout_details" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
  ]

  environment = {
    REQPATH           = "${var.create_state.ldap.ldap_mount}/library/test-set/status"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/read.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Check-in Specific Accounts
# Test Case #11: Check-in Specific Accounts - Explicitly check in accounts by name
resource "enos_remote_exec" "ldap_library_checkin_specific" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
    enos_remote_exec.ldap_library_checkout_details,
  ]

  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    SET_NAME          = "test-set"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap-library-checkin-specific.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Force Check-in (Admin)
# Test Case #13: Force Check-in (Admin) - Admin force check-in using /manage/ endpoint
resource "enos_remote_exec" "ldap_library_force_checkin" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
    enos_remote_exec.ldap_library_checkout_details,
    enos_remote_exec.ldap_library_checkin_specific,
  ]

  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    SET_NAME          = "test-set"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap-library-force-checkin.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Password Rotation on Check-in
# Test Case #16: Password Rotation on Check-in - Verify password is rotated when account is checked back in
resource "enos_remote_exec" "ldap_library_password_rotation" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
    enos_remote_exec.ldap_library_checkout_details,
    enos_remote_exec.ldap_library_checkin_specific,
    enos_remote_exec.ldap_library_force_checkin,
  ]

  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    SET_NAME          = "test-set"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap-library-password-rotation.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Password Retrieval on Check-out
# Test Case #17: Password Retrieval on Check-out - Verify password is returned when account is checked out
resource "enos_remote_exec" "ldap_library_password_checkout" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
    enos_remote_exec.ldap_library_checkout_details,
    enos_remote_exec.ldap_library_password_rotation,
    enos_remote_exec.ldap_library_checkin_specific,
    enos_remote_exec.ldap_library_force_checkin,
  ]

  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    SET_NAME          = "test-set"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap-library-password-checkout.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Audit Trail for All Operations
# Test Case #18: Audit Trail for All Operations - Verify Vault Core logs LDAP operations
resource "enos_remote_exec" "ldap_library_verify_audit_trail" {
  count = var.vault_audit_log_path != null ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
    enos_remote_exec.ldap_library_checkout_details,
    enos_remote_exec.ldap_library_password_rotation,
    enos_remote_exec.ldap_library_checkout_lease_renew,
    enos_remote_exec.ldap_library_checkout_lease_revoke,
  ]

  environment = {
    MOUNT          = var.create_state.ldap.ldap_mount
    AUDIT_LOG_PATH = var.vault_audit_log_path
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap/verify-audit-trail.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Optional Check-In Enforcement
# Test Case #19: Optional Check-In Enforcement - Configure whether check-in is required
resource "enos_remote_exec" "ldap_library_enforcement_config" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
  ]

  environment = {
    MOUNT                        = var.create_state.ldap.ldap_mount
    SET_NAME                     = "test-set-enforcement"
    SERVICE_ACCOUNT_NAMES        = "alice,carol"
    TTL                          = "10h"
    MAX_TTL                      = "20h"
    DISABLE_CHECK_IN_ENFORCEMENT = "true"
    VAULT_ADDR                   = var.vault_addr
    VAULT_INSTALL_DIR            = var.vault_install_dir
    VAULT_TOKEN                  = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap-library-enforcement-config.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# TTL Configuration
# Test Case #20: TTL Configuration - Configure TTL per library
resource "enos_remote_exec" "ldap_library_ttl_config" {
  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
  ]

  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    SET_NAME          = "test-set"
    TTL               = "1h"
    MAX_TTL           = "2h"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }

  scripts = [abspath("${path.module}/../../../scripts/ldap-library-ttl-config.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

# Verify Static role
resource "enos_remote_exec" "ldap_static_roles" {
  count = var.enable_static_role_verification ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_secrets,
    enos_remote_exec.ldap_verify_dynamic_credentials_suite,
    enos_remote_exec.ldap_verify_audit_trail
  ]
  environment = {
    MOUNT             = "${var.create_state.ldap.ldap_mount}"
    LDAP_SERVER       = "${var.create_state.ldap.host.private_ip}"
    LDAP_PORT         = "${var.create_state.ldap.port}"
    LDAP_USERNAME     = "${var.create_state.ldap.username}"
    LDAP_ADMIN_PW     = "${var.create_state.ldap.pw}"
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
  }
  scripts = [abspath("${path.module}/../../../scripts/ldap/static-roles.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

resource "enos_remote_exec" "ldap_verify_password_policy" {
  count = var.enable_password_policy_verification ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_rollback,
    enos_remote_exec.ldap_static_roles,
    enos_remote_exec.ldap_verify_dynamic_credentials_suite
  ]
  environment = {
    MOUNT             = var.create_state.ldap.ldap_mount
    LDAP_SERVER       = var.create_state.ldap.host.private_ip
    LDAP_PORT         = var.create_state.ldap.port
    LDAP_USERNAME     = var.create_state.ldap.username
    LDAP_ADMIN_PW     = var.create_state.ldap.pw
    VAULT_ADDR        = var.vault_addr
    VAULT_INSTALL_DIR = var.vault_install_dir
    VAULT_TOKEN       = var.vault_root_token
    STRONG_POLICY     = local.strong_password_policy
  }
  scripts = [abspath("${path.module}/../../../scripts/ldap/verify-password-policy.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}

resource "enos_remote_exec" "verify_audit_log" {
  count = var.enable_password_policy_verification && var.vault_audit_log_path != null ? 1 : 0

  depends_on = [
    enos_remote_exec.ldap_verify_rotation,
    enos_remote_exec.ldap_verify_password_policy
  ]
  environment = {
    VAULT_AUDIT_LOG = var.vault_audit_log_path
  }
  scripts = [abspath("${path.module}/../../../scripts/ldap/audit-verify.sh")]

  transport = {
    ssh = {
      host = var.hosts[0].public_ip
    }
  }
}
