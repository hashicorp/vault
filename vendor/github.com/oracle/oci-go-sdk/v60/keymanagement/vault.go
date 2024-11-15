// Copyright (c) 2016, 2018, 2022, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Vault Service Key Management API
//
// API for managing and performing operations with keys and vaults. (For the API for managing secrets, see the Vault Service
// Secret Management API. For the API for retrieving secrets, see the Vault Service Secret Retrieval API.)
//

package keymanagement

import (
	"fmt"
	"github.com/oracle/oci-go-sdk/v60/common"
	"strings"
)

// Vault The representation of Vault
type Vault struct {

	// The OCID of the compartment that contains this vault.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The service endpoint to perform cryptographic operations against. Cryptographic operations include
	// Encrypt (https://docs.cloud.oracle.com/api/#/en/key/latest/EncryptedData/Encrypt), Decrypt (https://docs.cloud.oracle.com/api/#/en/key/latest/DecryptedData/Decrypt),
	// and GenerateDataEncryptionKey (https://docs.cloud.oracle.com/api/#/en/key/latest/GeneratedKey/GenerateDataEncryptionKey) operations.
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

	ReplicaDetails *VaultReplicaDetails `mandatory:"false" json:"replicaDetails"`

	IsPrimary *bool `mandatory:"false" json:"isPrimary"`
}

func (m Vault) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m Vault) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingVaultLifecycleStateEnum(string(m.LifecycleState)); !ok && m.LifecycleState != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for LifecycleState: %s. Supported values are: %s.", m.LifecycleState, strings.Join(GetVaultLifecycleStateEnumStringValues(), ",")))
	}
	if _, ok := GetMappingVaultVaultTypeEnum(string(m.VaultType)); !ok && m.VaultType != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for VaultType: %s. Supported values are: %s.", m.VaultType, strings.Join(GetVaultVaultTypeEnumStringValues(), ",")))
	}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
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

var mappingVaultLifecycleStateEnum = map[string]VaultLifecycleStateEnum{
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

var mappingVaultLifecycleStateEnumLowerCase = map[string]VaultLifecycleStateEnum{
	"creating":            VaultLifecycleStateCreating,
	"active":              VaultLifecycleStateActive,
	"deleting":            VaultLifecycleStateDeleting,
	"deleted":             VaultLifecycleStateDeleted,
	"pending_deletion":    VaultLifecycleStatePendingDeletion,
	"scheduling_deletion": VaultLifecycleStateSchedulingDeletion,
	"cancelling_deletion": VaultLifecycleStateCancellingDeletion,
	"updating":            VaultLifecycleStateUpdating,
	"backup_in_progress":  VaultLifecycleStateBackupInProgress,
	"restoring":           VaultLifecycleStateRestoring,
}

// GetVaultLifecycleStateEnumValues Enumerates the set of values for VaultLifecycleStateEnum
func GetVaultLifecycleStateEnumValues() []VaultLifecycleStateEnum {
	values := make([]VaultLifecycleStateEnum, 0)
	for _, v := range mappingVaultLifecycleStateEnum {
		values = append(values, v)
	}
	return values
}

// GetVaultLifecycleStateEnumStringValues Enumerates the set of values in String for VaultLifecycleStateEnum
func GetVaultLifecycleStateEnumStringValues() []string {
	return []string{
		"CREATING",
		"ACTIVE",
		"DELETING",
		"DELETED",
		"PENDING_DELETION",
		"SCHEDULING_DELETION",
		"CANCELLING_DELETION",
		"UPDATING",
		"BACKUP_IN_PROGRESS",
		"RESTORING",
	}
}

// GetMappingVaultLifecycleStateEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingVaultLifecycleStateEnum(val string) (VaultLifecycleStateEnum, bool) {
	enum, ok := mappingVaultLifecycleStateEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}

// VaultVaultTypeEnum Enum with underlying type: string
type VaultVaultTypeEnum string

// Set of constants representing the allowable values for VaultVaultTypeEnum
const (
	VaultVaultTypeVirtualPrivate VaultVaultTypeEnum = "VIRTUAL_PRIVATE"
	VaultVaultTypeDefault        VaultVaultTypeEnum = "DEFAULT"
)

var mappingVaultVaultTypeEnum = map[string]VaultVaultTypeEnum{
	"VIRTUAL_PRIVATE": VaultVaultTypeVirtualPrivate,
	"DEFAULT":         VaultVaultTypeDefault,
}

var mappingVaultVaultTypeEnumLowerCase = map[string]VaultVaultTypeEnum{
	"virtual_private": VaultVaultTypeVirtualPrivate,
	"default":         VaultVaultTypeDefault,
}

// GetVaultVaultTypeEnumValues Enumerates the set of values for VaultVaultTypeEnum
func GetVaultVaultTypeEnumValues() []VaultVaultTypeEnum {
	values := make([]VaultVaultTypeEnum, 0)
	for _, v := range mappingVaultVaultTypeEnum {
		values = append(values, v)
	}
	return values
}

// GetVaultVaultTypeEnumStringValues Enumerates the set of values in String for VaultVaultTypeEnum
func GetVaultVaultTypeEnumStringValues() []string {
	return []string{
		"VIRTUAL_PRIVATE",
		"DEFAULT",
	}
}

// GetMappingVaultVaultTypeEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingVaultVaultTypeEnum(val string) (VaultVaultTypeEnum, bool) {
	enum, ok := mappingVaultVaultTypeEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
