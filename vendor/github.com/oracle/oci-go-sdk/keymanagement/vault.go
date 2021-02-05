// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Key Management Service API
//
// API for managing and performing operations with keys and vaults.
//

package keymanagement

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Vault The representation of Vault
type Vault struct {

	// The OCID of the compartment that contains this vault.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The service endpoint to perform cryptographic operations against. Cryptographic operations include
	// Encrypt (https://docs.cloud.oracle.com/api/#/en/key/release/EncryptedData/Encrypt), Decrypt (https://docs.cloud.oracle.com/api/#/en/key/release/DecryptedData/Decrypt),
	// and GenerateDataEncryptionKey (https://docs.cloud.oracle.com/api/#/en/key/release/GeneratedKey/GenerateDataEncryptionKey) operations.
	CryptoEndpoint *string `mandatory:"true" json:"cryptoEndpoint"`

	// A user-friendly name for the vault. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID of the vault.
	Id *string `mandatory:"true" json:"id"`

	// The vault's current lifecycle state.
	// Example: `DELETED`
	LifecycleState VaultLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The service endpoint to perform management operations against. Management operations include "Create," "Update," "List," "Get," and "Delete" operations.
	ManagementEndpoint *string `mandatory:"true" json:"managementEndpoint"`

	// The date and time this vault was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2018-04-03T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The type of vault. Each type of vault stores the key with different
	// degrees of isolation and has different options and pricing.
	VaultType VaultVaultTypeEnum `mandatory:"true" json:"vaultType"`

	// The OCID of the vault's wrapping key.
	WrappingkeyId *string `mandatory:"true" json:"wrappingkeyId"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// An optional property to indicate when to delete the vault, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2018-04-03T21:10:29.600Z`
	TimeOfDeletion *common.SDKTime `mandatory:"false" json:"timeOfDeletion"`

	// The OCID of the vault from which this vault was restored, if it was restored from a backup file.
	// If you restore a vault to the same region, the vault retains the same OCID that it had when you
	// backed up the vault.
	RestoredFromVaultId *string `mandatory:"false" json:"restoredFromVaultId"`
}

func (m Vault) String() string {
	return common.PointerString(m)
}

// VaultLifecycleStateEnum Enum with underlying type: string
type VaultLifecycleStateEnum string

// Set of constants representing the allowable values for VaultLifecycleStateEnum
const (
	VaultLifecycleStateCreating           VaultLifecycleStateEnum = "CREATING"
	VaultLifecycleStateActive             VaultLifecycleStateEnum = "ACTIVE"
	VaultLifecycleStateDeleting           VaultLifecycleStateEnum = "DELETING"
	VaultLifecycleStateDeleted            VaultLifecycleStateEnum = "DELETED"
	VaultLifecycleStatePendingDeletion    VaultLifecycleStateEnum = "PENDING_DELETION"
	VaultLifecycleStateSchedulingDeletion VaultLifecycleStateEnum = "SCHEDULING_DELETION"
	VaultLifecycleStateCancellingDeletion VaultLifecycleStateEnum = "CANCELLING_DELETION"
	VaultLifecycleStateUpdating           VaultLifecycleStateEnum = "UPDATING"
	VaultLifecycleStateBackupInProgress   VaultLifecycleStateEnum = "BACKUP_IN_PROGRESS"
	VaultLifecycleStateRestoring          VaultLifecycleStateEnum = "RESTORING"
)

var mappingVaultLifecycleState = map[string]VaultLifecycleStateEnum{
	"CREATING":            VaultLifecycleStateCreating,
	"ACTIVE":              VaultLifecycleStateActive,
	"DELETING":            VaultLifecycleStateDeleting,
	"DELETED":             VaultLifecycleStateDeleted,
	"PENDING_DELETION":    VaultLifecycleStatePendingDeletion,
	"SCHEDULING_DELETION": VaultLifecycleStateSchedulingDeletion,
	"CANCELLING_DELETION": VaultLifecycleStateCancellingDeletion,
	"UPDATING":            VaultLifecycleStateUpdating,
	"BACKUP_IN_PROGRESS":  VaultLifecycleStateBackupInProgress,
	"RESTORING":           VaultLifecycleStateRestoring,
}

// GetVaultLifecycleStateEnumValues Enumerates the set of values for VaultLifecycleStateEnum
func GetVaultLifecycleStateEnumValues() []VaultLifecycleStateEnum {
	values := make([]VaultLifecycleStateEnum, 0)
	for _, v := range mappingVaultLifecycleState {
		values = append(values, v)
	}
	return values
}

// VaultVaultTypeEnum Enum with underlying type: string
type VaultVaultTypeEnum string

// Set of constants representing the allowable values for VaultVaultTypeEnum
const (
	VaultVaultTypeVirtualPrivate VaultVaultTypeEnum = "VIRTUAL_PRIVATE"
	VaultVaultTypeDefault        VaultVaultTypeEnum = "DEFAULT"
)

var mappingVaultVaultType = map[string]VaultVaultTypeEnum{
	"VIRTUAL_PRIVATE": VaultVaultTypeVirtualPrivate,
	"DEFAULT":         VaultVaultTypeDefault,
}

// GetVaultVaultTypeEnumValues Enumerates the set of values for VaultVaultTypeEnum
func GetVaultVaultTypeEnumValues() []VaultVaultTypeEnum {
	values := make([]VaultVaultTypeEnum, 0)
	for _, v := range mappingVaultVaultType {
		values = append(values, v)
	}
	return values
}
