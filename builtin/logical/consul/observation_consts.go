// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package consul

const (
	ObservationTypeConsulConfigAccessWrite = "consul/config/access/write"
	ObservationTypeConsulConfigAccessRead  = "consul/config/access/read"

	ObservationTypeConsulRoleWrite  = "consul/role/write"
	ObservationTypeConsulRoleRead   = "consul/role/read"
	ObservationTypeConsulRoleDelete = "consul/role/delete"

	ObservationTypeConsulTokenRead = "consul/token/read"

	ObservationTypeConsulCredentialRenew  = "consul/credential/renew"
	ObservationTypeConsulCredentialRevoke = "consul/credential/revoke"
)
