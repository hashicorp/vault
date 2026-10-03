// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
)

// authTuneRequestFields returns the request fields for auth tuning endpoints.
// Used in:
// - POST /sys/auth/{path}/tune
// - POST /sys/mounts/auth/{path}/tune
func authTuneRequestFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"allowed_response_headers": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["allowed_response_headers"][0]),
		},
		"audit_non_hmac_request_keys": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_request_keys"][0]),
		},
		"audit_non_hmac_response_keys": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_response_keys"][0]),
		},
		"default_lease_ttl": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["tune_default_lease_ttl"][0]),
		},
		"description": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
		},
		"identity_token_key": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["identity_token_key"][0]),
			Required:    false,
		},
		"listing_visibility": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["listing_visibility"][0]),
		},
		"max_lease_ttl": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["tune_max_lease_ttl"][0]),
		},
		"options": {
			Type:        framework.TypeKVPairs,
			Description: strings.TrimSpace(sysHelp["tune_mount_options"][0]),
		},
		"passthrough_request_headers": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["passthrough_request_headers"][0]),
		},
		"path": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["auth_tune"][0]),
		},
		"plugin_version": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
		},
		"token_type": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["token_type"][0]),
		},
		"trim_request_trailing_slashes": {
			Type:     framework.TypeBool,
			Required: false,
		},
		"user_lockout_config": {
			Type:        framework.TypeMap,
			Description: strings.TrimSpace(sysHelp["tune_user_lockout_config"][0]),
		},
	}

	entAddAuthTuneRequestFields(fields)

	return fields
}

// authTuneResponseFields returns the response fields for auth tuning endpoints.
// Used in:
// - GET /sys/auth/{path}/tune
// - GET /sys/mounts/auth/{path}/tune
func authTuneResponseFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"allowed_managed_keys": {
			Type:     framework.TypeCommaStringSlice,
			Required: false,
		},
		"allowed_response_headers": {
			Type:     framework.TypeCommaStringSlice,
			Required: false,
		},
		"audit_non_hmac_request_keys": {
			Type:     framework.TypeCommaStringSlice,
			Required: false,
		},
		"audit_non_hmac_response_keys": {
			Type:     framework.TypeCommaStringSlice,
			Required: false,
		},
		"default_lease_ttl": {
			Type:     framework.TypeInt,
			Required: true,
		},
		"description": {
			Type:     framework.TypeString,
			Required: true,
		},
		"external_entropy_access": {
			Type:     framework.TypeBool,
			Required: false,
		},
		"force_no_cache": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"identity_token_key": {
			Type:     framework.TypeString,
			Required: false,
		},
		"listing_visibility": {
			Type:     framework.TypeString,
			Required: false,
		},
		"max_lease_ttl": {
			Type:     framework.TypeInt,
			Required: true,
		},
		"options": {
			Type:     framework.TypeMap,
			Required: false,
		},
		"passthrough_request_headers": {
			Type:     framework.TypeCommaStringSlice,
			Required: false,
		},
		"plugin_version": {
			Type:     framework.TypeString,
			Required: false,
		},
		"token_type": {
			Type:     framework.TypeString,
			Required: false,
		},
		"trim_request_trailing_slashes": {
			Type:     framework.TypeBool,
			Required: false,
		},
		"user_lockout_counter_reset_duration": {
			Type:     framework.TypeInt64,
			Required: false,
		},
		"user_lockout_disable": {
			Type:     framework.TypeBool,
			Required: false,
		},
		"user_lockout_duration": {
			Type:     framework.TypeInt64,
			Required: false,
		},
		"user_lockout_threshold": {
			Type:     framework.TypeInt64, // uint64
			Required: false,
		},
	}

	entAddAuthTuneResponseFields(fields)

	return fields
}

// authRequestFields returns the request fields for auth engine mount/unmount operations.
// Used in:
// - POST /sys/auth/{path} (mount auth engine)
// - DELETE /sys/auth/{path} (unmount auth engine)
func authRequestFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"config": {
			Type:        framework.TypeMap,
			Description: strings.TrimSpace(sysHelp["auth_config"][0]),
		},
		"description": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
		},
		"external_entropy_access": {
			Type:        framework.TypeBool,
			Default:     false,
			Description: strings.TrimSpace(sysHelp["external_entropy_access"][0]),
		},
		"local": {
			Type:        framework.TypeBool,
			Default:     false,
			Description: strings.TrimSpace(sysHelp["mount_local"][0]),
		},
		"options": {
			Type:        framework.TypeKVPairs,
			Description: strings.TrimSpace(sysHelp["auth_options"][0]),
		},
		"path": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["auth_path"][0]),
		},
		"plugin_name": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["auth_plugin"][0]),
		},
		"plugin_version": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
		},
		"seal_wrap": {
			Type:        framework.TypeBool,
			Default:     false,
			Description: strings.TrimSpace(sysHelp["seal_wrap"][0]),
		},
		"type": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["auth_type"][0]),
		},
	}

	entAddAuthRequestFields(fields)

	return fields
}

// authResponseFields returns the response fields for auth engine read operations.
// Used in:
// - GET /sys/auth/{path} (read auth engine configuration)
func authResponseFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"accessor": {
			Type:     framework.TypeString,
			Required: true,
		},
		"config": {
			Type:     framework.TypeMap,
			Required: true,
		},
		"deprecation_status": {
			Type:     framework.TypeString,
			Required: false,
		},
		"description": {
			Type:     framework.TypeString,
			Required: true,
		},
		"external_entropy_access": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"local": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"options": {
			Type:     framework.TypeMap,
			Required: true,
		},
		"plugin_version": {
			Type:     framework.TypeString,
			Required: true,
		},
		"running_plugin_version": {
			Type:     framework.TypeString,
			Required: true,
		},
		"running_sha256": {
			Type:     framework.TypeString,
			Required: true,
		},
		"seal_wrap": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"type": {
			Type:     framework.TypeString,
			Required: true,
		},
		"uuid": {
			Type:     framework.TypeString,
			Required: true,
		},
	}

	entAddAuthResponseFields(fields)

	return fields
}

// secretsTuneRequestFields returns the request fields for secrets engine tuning operations.
// Used in:
// - POST /sys/mounts/{path}/tune
func secretsTuneRequestFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"allowed_managed_keys": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["tune_allowed_managed_keys"][0]),
		},
		"allowed_response_headers": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["allowed_response_headers"][0]),
		},
		"audit_non_hmac_request_keys": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_request_keys"][0]),
		},
		"audit_non_hmac_response_keys": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["tune_audit_non_hmac_response_keys"][0]),
		},
		"default_lease_ttl": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["tune_default_lease_ttl"][0]),
		},
		"delegated_auth_accessors": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["allowed_delegated_auth_accessors"][0]),
		},
		"description": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
		},
		"identity_token_key": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["identity_token_key"][0]),
		},
		"listing_visibility": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["listing_visibility"][0]),
		},
		"max_lease_ttl": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["tune_max_lease_ttl"][0]),
		},
		"options": {
			Type:        framework.TypeKVPairs,
			Description: strings.TrimSpace(sysHelp["tune_mount_options"][0]),
		},
		"passthrough_request_headers": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["passthrough_request_headers"][0]),
		},
		"path": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["mount_path"][0]),
		},
		"plugin_version": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
		},
		"token_type": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["token_type"][0]),
		},
		"trim_request_trailing_slashes": {
			Type:        framework.TypeBool,
			Description: strings.TrimSpace(sysHelp["trim_request_trailing_slashes"][0]),
		},
		"user_lockout_config": {
			Type:        framework.TypeMap,
			Description: strings.TrimSpace(sysHelp["tune_user_lockout_config"][0]),
		},
	}

	entAddSecretsTuneRequestFields(fields)

	return fields
}

// secretsTuneResponseFields returns the response fields for secrets engine tuning operations.
// Used in:
// - GET /sys/mounts/{path}/tune
func secretsTuneResponseFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"allowed_managed_keys": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["tune_allowed_managed_keys"][0]),
			Required:    false,
		},
		"allowed_response_headers": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["allowed_response_headers"][0]),
			Required:    false,
		},
		"audit_non_hmac_request_keys": {
			Type:     framework.TypeCommaStringSlice,
			Required: false,
		},
		"audit_non_hmac_response_keys": {
			Type:     framework.TypeCommaStringSlice,
			Required: false,
		},
		"default_lease_ttl": {
			Type:        framework.TypeInt,
			Description: strings.TrimSpace(sysHelp["tune_default_lease_ttl"][0]),
			Required:    true,
		},
		"delegated_auth_accessors": {
			Type:        framework.TypeCommaStringSlice,
			Description: strings.TrimSpace(sysHelp["delegated_auth_accessors"][0]),
			Required:    false,
		},
		"description": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["auth_desc"][0]),
			Required:    true,
		},
		"external_entropy_access": {
			Type:     framework.TypeBool,
			Required: false,
		},
		"force_no_cache": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"identity_token_key": {
			Type:     framework.TypeString,
			Required: false,
		},
		"listing_visibility": {
			Type:     framework.TypeString,
			Required: false,
		},
		"max_lease_ttl": {
			Type:        framework.TypeInt,
			Description: strings.TrimSpace(sysHelp["tune_max_lease_ttl"][0]),
			Required:    true,
		},
		"options": {
			Type:        framework.TypeKVPairs,
			Description: strings.TrimSpace(sysHelp["tune_mount_options"][0]),
			Required:    false,
		},
		"passthrough_request_headers": {
			Type:     framework.TypeCommaStringSlice,
			Required: false,
		},
		"plugin_version": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
			Required:    false,
		},
		"token_type": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["token_type"][0]),
			Required:    false,
		},
		"trim_request_trailing_slashes": {
			Type:     framework.TypeBool,
			Required: false,
		},
		"user_lockout_counter_reset_duration": {
			Type:     framework.TypeInt64,
			Required: false,
		},
		"user_lockout_disable": {
			Type:     framework.TypeBool,
			Required: false,
		},
		"user_lockout_duration": {
			Type:     framework.TypeInt64,
			Required: false,
		},
		"user_lockout_threshold": {
			Type:     framework.TypeInt64, // TODO this is actuall a Uint64 do we need a new type?
			Required: false,
		},
	}

	entAddSecretsTuneResponseFields(fields)

	return fields
}

// secretsRequestFields returns the request fields for secrets engine mount/unmount operations.
// Used in:
// - POST /sys/mounts/{path} (mount secrets engine)
// - DELETE /sys/mounts/{path} (unmount secrets engine)
func secretsRequestFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"config": {
			Type:        framework.TypeMap,
			Description: strings.TrimSpace(sysHelp["mount_config"][0]),
		},
		"description": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["mount_desc"][0]),
		},
		"external_entropy_access": {
			Type:        framework.TypeBool,
			Default:     false,
			Description: strings.TrimSpace(sysHelp["external_entropy_access"][0]),
		},
		"local": {
			Type:        framework.TypeBool,
			Default:     false,
			Description: strings.TrimSpace(sysHelp["mount_local"][0]),
		},
		"options": {
			Type:        framework.TypeKVPairs,
			Description: strings.TrimSpace(sysHelp["mount_options"][0]),
		},
		"path": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["mount_path"][0]),
		},
		"plugin_name": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["mount_plugin_name"][0]),
		},
		"plugin_version": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
		},
		"seal_wrap": {
			Type:        framework.TypeBool,
			Default:     false,
			Description: strings.TrimSpace(sysHelp["seal_wrap"][0]),
		},
		"type": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["mount_type"][0]),
		},
	}

	entAddSecretsRequestFields(fields)

	return fields
}

// secretsResponseFields returns the response fields for secrets engine read operations.
// Used in:
// - GET /sys/mounts/{path} (read secrets engine configuration)
func secretsResponseFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"accessor": {
			Type:     framework.TypeString,
			Required: true,
		},
		"config": {
			Type:        framework.TypeMap,
			Description: strings.TrimSpace(sysHelp["mount_config"][0]),
			Required:    true,
		},
		"deprecation_status": {
			Type:     framework.TypeString,
			Required: false,
		},
		"description": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["mount_desc"][0]),
			Required:    true,
		},
		"external_entropy_access": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"local": {
			Type:        framework.TypeBool,
			Default:     false,
			Description: strings.TrimSpace(sysHelp["mount_local"][0]),
			Required:    true,
		},
		"options": {
			Type:        framework.TypeKVPairs,
			Description: strings.TrimSpace(sysHelp["mount_options"][0]),
			Required:    true,
		},
		"plugin_version": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["plugin-catalog_version"][0]),
			Required:    true,
		},
		"running_plugin_version": {
			Type:     framework.TypeString,
			Required: true,
		},
		"running_sha256": {
			Type:     framework.TypeString,
			Required: true,
		},
		"seal_wrap": {
			Type:        framework.TypeBool,
			Default:     false,
			Description: strings.TrimSpace(sysHelp["seal_wrap"][0]),
			Required:    true,
		},
		"type": {
			Type:        framework.TypeString,
			Description: strings.TrimSpace(sysHelp["mount_type"][0]),
			Required:    true,
		},
		"uuid": {
			Type:     framework.TypeString,
			Required: true,
		},
	}

	entAddSecretsResponseFields(fields)

	return fields
}

// internalUIMountsPathResponse returns the response fields for the internal UI mounts path.
// Used in:
// - GET /sys/internal/ui/mounts/{path}
func internalUIMountsPathResponseFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"accessor": {
			Type:     framework.TypeString,
			Required: true,
		},
		"config": {
			Type:     framework.TypeMap,
			Required: true,
		},
		"description": {
			Type:     framework.TypeString,
			Required: true,
		},
		"external_entropy_access": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"local": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"options": {
			Type:     framework.TypeMap,
			Required: true,
		},
		"path": {
			Type:     framework.TypeString,
			Required: true,
		},
		"plugin_version": {
			Type:     framework.TypeString,
			Required: true,
		},
		"running_plugin_version": {
			Type:     framework.TypeString,
			Required: true,
		},
		"running_sha256": {
			Type:     framework.TypeString,
			Required: true,
		},
		"seal_wrap": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"type": {
			Type:     framework.TypeString,
			Required: true,
		},
		"uuid": {
			Type:     framework.TypeString,
			Required: true,
		},
	}

	return fields
}
