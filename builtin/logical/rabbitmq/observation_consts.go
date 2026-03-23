// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package rabbitmq

const (
	ObservationTypeRabbitMQConnectionConfigWrite   = "rabbitmq/connection/config/write"
	ObservationTypeRabbitMQLeaseConfigWrite        = "rabbitmq/lease/config/write"
	ObservationTypeRabbitMQLeaseConfigRead         = "rabbitmq/lease/config/read"
	ObservationTypeRabbitMQRoleWrite               = "rabbitmq/role/write"
	ObservationTypeRabbitMQRoleRead                = "rabbitmq/role/read"
	ObservationTypeRabbitMQRoleDelete              = "rabbitmq/role/delete"
	ObservationTypeRabbitMQCredentialCreateSuccess = "rabbitmq/credential/create/success"
	ObservationTypeRabbitMQCredentialCreateFail    = "rabbitmq/credential/create/fail"
	ObservationTypeRabbitMQCredentialRenew         = "rabbitmq/credential/renew"
	ObservationTypeRabbitMQCredentialRevoke        = "rabbitmq/credential/revoke"
)
