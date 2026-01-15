// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package ssh

const (
	// ObservationTypeSSHConfigZeroAddressRead - Metadata: role_names ([]string)
	ObservationTypeSSHConfigZeroAddressRead = "ssh/config/zero-address/read"
	// ObservationTypeSSHConfigZeroAddressWrite - Metadata: role_names ([]string)
	ObservationTypeSSHConfigZeroAddressWrite = "ssh/config/zero-address/write"
	// ObservationTypeSSHConfigZeroAddressDelete - Metadata: none
	ObservationTypeSSHConfigZeroAddressDelete = "ssh/config/zero-address/delete"

	// ObservationTypeSSHRoleRead - Metadata: role_name, key_type, and for CA roles: ttl, max_ttl,
	// allow_user_certificates, allow_host_certificates, allow_bare_domains, allow_subdomains,
	// allow_user_key_ids, allowed_users_template, allowed_domains_template, default_user_template,
	// default_extensions_template, algorithm_signer, not_before_duration, allow_empty_principals
	ObservationTypeSSHRoleRead = "ssh/role/read"
	// ObservationTypeSSHRoleWrite - Metadata: role_name, key_type, and for CA roles: ttl, max_ttl,
	// allow_user_certificates, allow_host_certificates, allow_bare_domains, allow_subdomains,
	// allow_user_key_ids, allowed_users_template, allowed_domains_template, default_user_template,
	// default_extensions_template, algorithm_signer, not_before_duration, allow_empty_principals
	ObservationTypeSSHRoleWrite = "ssh/role/write"
	// ObservationTypeSSHRoleDelete - Metadata: role_name
	ObservationTypeSSHRoleDelete = "ssh/role/delete"

	// ObservationTypeSSHOTPCreate - Metadata: role_name, key_type, and for CA roles: ttl, max_ttl,
	// allow_user_certificates, allow_host_certificates, allow_bare_domains, allow_subdomains,
	// allow_user_key_ids, allowed_users_template, allowed_domains_template, default_user_template,
	// default_extensions_template, algorithm_signer, not_before_duration, allow_empty_principals
	ObservationTypeSSHOTPCreate = "ssh/otp/create"
	// ObservationTypeSSHOTPRevoke - Metadata: none
	ObservationTypeSSHOTPRevoke = "ssh/otp/revoke"
	// ObservationTypeSSHOTPVerify - Metadata: role_name
	ObservationTypeSSHOTPVerify = "ssh/otp/verify"

	// ObservationTypeSSHLookup - Metadata: role_names ([]string)
	ObservationTypeSSHLookup = "ssh/lookup"
)
