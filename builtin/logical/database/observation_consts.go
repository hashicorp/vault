// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package database

const (
	// Connection related observations:

	ObservationTypeDatabaseConfigWrite     = "database/connection/config/write"
	ObservationTypeDatabaseConfigDelete    = "database/connection/config/delete"
	ObservationTypeDatabaseConfigRead      = "database/connection/config/read"
	ObservationTypeDatabaseConnectionReset = "database/connection/reset"

	// Reload related observations:
	// Note: the following three observations mean that, for any reload, Vault will emit
	// n+1 observations, where n is the number of connections that we attempt to reload.
	// database/plugin/reload will not contain a connection_name, as it will be a summary
	// database/connection/reload/* will be per-plugin, and will contain a connection_name

	// ObservationTypeDatabaseReloadPlugin is emitted when a plugin reload is issued
	ObservationTypeDatabaseReloadPlugin = "database/plugin/reload"
	// ObservationTypeDatabaseReloadSuccess is emitted for each connection successfully reloaded
	ObservationTypeDatabaseReloadSuccess = "database/connection/reload/success"
	// ObservationTypeDatabaseReloadFail is emitted for each connection unsuccessfully reloaded
	ObservationTypeDatabaseReloadFail = "database/connection/reload/fail"

	// Role related observations

	ObservationTypeDatabaseRoleCreate = "database/role/create"
	ObservationTypeDatabaseRoleUpdate = "database/role/update"
	ObservationTypeDatabaseRoleRead   = "database/role/read"
	// ObservationTypeDatabaseRoleDelete is emitted whenever a role is deleted.
	// Note that this observation does not include a connection_name, to avoid doing an
	// additional storage read.
	ObservationTypeDatabaseRoleDelete = "database/role/delete"

	ObservationTypeDatabaseCredentialCreateSuccess = "database/credential/create/success"
	ObservationTypeDatabaseCredentialCreateFail    = "database/credential/create/fail"
	ObservationTypeDatabaseCredentialRenew         = "database/credential/renew"
	ObservationTypeDatabaseCredentialRevoke        = "database/credential/revoke"

	// Rotate-root observations

	ObservationTypeDatabaseRotateRootSuccess = "database/rotate-root/success"
	ObservationTypeDatabaseRotateRootFailure = "database/rotate-root/fail"

	// Static role observations

	ObservationTypeDatabaseRotateStaticRoleSuccess = "database/static-role/rotate/success"
	ObservationTypeDatabaseRotateStaticRoleFailure = "database/static-role/rotate/fail"
	ObservationTypeDatabaseStaticRoleCreate        = "database/static-role/create"
	ObservationTypeDatabaseStaticRoleUpdate        = "database/static-role/update"
	ObservationTypeDatabaseStaticRoleRead          = "database/static-role/read"
	// ObservationTypeDatabaseStaticRoleDelete is emitted whenever a static role is deleted.
	// Note that this observation does not include a connection_name, to avoid doing an
	// additional storage read.
	ObservationTypeDatabaseStaticRoleDelete = "database/static-role/delete"

	ObservationTypeDatabaseStaticCredentialRead = "database/static-credential/read"
)
