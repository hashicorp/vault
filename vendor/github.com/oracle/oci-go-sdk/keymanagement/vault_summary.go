// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Key Management Service API
//
// API for managing and performing operations with keys and vaults.
//

package keymanagement

import (
	"github.com/oracle/oci-go-sdk/common"
)

// VaultSummary The representation of VaultSummary
type VaultSummary struct {

	// The OCID of the compartment that contains a particular vault.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The service endpoint to perform cryptographic operations against. Cryptographic operations include 'Encrypt,' 'Decrypt,' and 'GenerateDataEncryptionKey' operations.
	CryptoEndpoint *string `mandatory:"true" json:"cryptoEndpoint"`

	// A user-friendly name for a vault. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID of a vault.
	Id *string `mandatory:"true" json:"id"`

	// A vault's current state.
	// Example: `ACTIVE`
	LifecycleState VaultSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The service endpoint to perform management operations against. Management operations include 'Create,' 'Update,' 'List,' 'Get,' and 'Delete' operations.
	ManagementEndpoint *string `mandatory:"true" json:"managementEndpoint"`

	// The date and time a vault was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2018-04-03T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The type of vault. Each type of vault stores keys with different degrees of isolation and has different options and pricing.
	VaultType VaultSummaryVaultTypeEnum `mandatory:"true" json:"vaultType"`

	// Usage of predefined tag keys. These predefined keys are scoped to namespaces.
	// Example: `{"foo-namespace": {"bar-key": "foo-value"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Simple key-value pair that is applied without any predefined name, type, or scope.
	// Exists for cross-compatibility only.
	// Example: `{"bar-key": "value"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m VaultSummary) String() string {
	return common.PointerString(m)
}

// VaultSummaryLifecycleStateEnum Enum with underlying type: string
type VaultSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for VaultSummaryLifecycleStateEnum
const (
	VaultSummaryLifecycleStateCreating           VaultSummaryLifecycleStateEnum = "CREATING"
	VaultSummaryLifecycleStateActive             VaultSummaryLifecycleStateEnum = "ACTIVE"
	VaultSummaryLifecycleStateDeleting           VaultSummaryLifecycleStateEnum = "DELETING"
	VaultSummaryLifecycleStateDeleted            VaultSummaryLifecycleStateEnum = "DELETED"
	VaultSummaryLifecycleStatePendingDeletion    VaultSummaryLifecycleStateEnum = "PENDING_DELETION"
	VaultSummaryLifecycleStateSchedulingDeletion VaultSummaryLifecycleStateEnum = "SCHEDULING_DELETION"
	VaultSummaryLifecycleStateCancellingDeletion VaultSummaryLifecycleStateEnum = "CANCELLING_DELETION"
	VaultSummaryLifecycleStateUpdating           VaultSummaryLifecycleStateEnum = "UPDATING"
)

var mappingVaultSummaryLifecycleState = map[string]VaultSummaryLifecycleStateEnum{
	"CREATING":            VaultSummaryLifecycleStateCreating,
	"ACTIVE":              VaultSummaryLifecycleStateActive,
	"DELETING":            VaultSummaryLifecycleStateDeleting,
	"DELETED":             VaultSummaryLifecycleStateDeleted,
	"PENDING_DELETION":    VaultSummaryLifecycleStatePendingDeletion,
	"SCHEDULING_DELETION": VaultSummaryLifecycleStateSchedulingDeletion,
	"CANCELLING_DELETION": VaultSummaryLifecycleStateCancellingDeletion,
	"UPDATING":            VaultSummaryLifecycleStateUpdating,
}

// GetVaultSummaryLifecycleStateEnumValues Enumerates the set of values for VaultSummaryLifecycleStateEnum
func GetVaultSummaryLifecycleStateEnumValues() []VaultSummaryLifecycleStateEnum {
	values := make([]VaultSummaryLifecycleStateEnum, 0)
	for _, v := range mappingVaultSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}

// VaultSummaryVaultTypeEnum Enum with underlying type: string
type VaultSummaryVaultTypeEnum string

// Set of constants representing the allowable values for VaultSummaryVaultTypeEnum
const (
	VaultSummaryVaultTypeVirtualPrivate VaultSummaryVaultTypeEnum = "VIRTUAL_PRIVATE"
)

var mappingVaultSummaryVaultType = map[string]VaultSummaryVaultTypeEnum{
	"VIRTUAL_PRIVATE": VaultSummaryVaultTypeVirtualPrivate,
}

// GetVaultSummaryVaultTypeEnumValues Enumerates the set of values for VaultSummaryVaultTypeEnum
func GetVaultSummaryVaultTypeEnumValues() []VaultSummaryVaultTypeEnum {
	values := make([]VaultSummaryVaultTypeEnum, 0)
	for _, v := range mappingVaultSummaryVaultType {
		values = append(values, v)
	}
	return values
}
