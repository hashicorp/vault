// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package nomad

const (
	// Role observation types
	// These contain role_name, token_type, global as metadata

	ObservationTypeNomadRoleWrite = "nomad/role/write"
	ObservationTypeNomadRoleRead  = "nomad/role/read"

	// Role delete observation type
	// This contains role_name as metadata

	ObservationTypeNomadRoleDelete = "nomad/role/delete"

	// Config observation types
	// These contain no metadata

	ObservationTypeNomadConfigAccessRead   = "nomad/config/access/read"
	ObservationTypeNomadConfigAccessWrite  = "nomad/config/access/write"
	ObservationTypeNomadConfigAccessDelete = "nomad/config/access/delete"

	// Lease config observation types
	// These contain ttl and max_ttl as metadata

	ObservationTypeNomadConfigLeaseRead   = "nomad/config/lease/read"
	ObservationTypeNomadConfigLeaseWrite  = "nomad/config/lease/write"
	ObservationTypeNomadConfigLeaseDelete = "nomad/config/lease/delete"

	// Credential success observation type
	// This contains role_name, token_type, global, ttl, max_ttl, accessor_id as metadata

	ObservationTypeNomadCredentialCreateSuccess = "nomad/credential/create/success"

	// Credential fail observation type
	// This contains role_name, token_type, global, ttl, max_ttl as metadata

	ObservationTypeNomadCredentialCreateFail = "nomad/credential/create/fail"

	// Credential revoke observation type
	// This contains accessor_id as metadata

	ObservationTypeNomadCredentialRevoke = "nomad/credential/revoke"

	// Credential renew observation type
	// This contains ttl, max_ttl, accessor_id as metadata

	ObservationTypeNomadCredentialRenew = "nomad/credential/renew"
)
