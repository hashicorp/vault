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

// VaultSummary The representation of VaultSummary
type VaultSummary struct {

	// The OCID of the compartment that contains a particular vault.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The service endpoint to perform cryptographic operations against. Cryptographic operations include
	// Encrypt (https://docs.cloud.oracle.com/api/#/en/key/latest/EncryptedData/Encrypt), Decrypt (https://docs.cloud.oracle.com/api/#/en/key/latest/DecryptedData/Decrypt),
	// and GenerateDataEncryptionKey (https://docs.cloud.oracle.com/api/#/en/key/latest/GeneratedKey/GenerateDataEncryptionKey) operations.
	CryptoEndpoint *string `mandatory:"true" json:"cryptoEndpoint"`

	// A user-friendly name for a vault. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID of a vault.
	Id *string `mandatory:"true" json:"id"`

	// A vault's current lifecycle state.
	// Example: `ACTIVE`
	LifecycleState VaultSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The service endpoint to perform management operations against. Management operations include "Create," "Update," "List," "Get," and "Delete" operations.
	ManagementEndpoint *string `mandatory:"true" json:"managementEndpoint"`

	// The date and time a vault was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2018-04-03T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The type of vault. Each type of vault stores keys with different
	// degrees of isolation and has different options and pricing.
	VaultType VaultSummaryVaultTypeEnum `mandatory:"true" json:"vaultType"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m VaultSummary) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m VaultSummary) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingVaultSummaryLifecycleStateEnum(string(m.LifecycleState)); !ok && m.LifecycleState != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for LifecycleState: %s. Supported values are: %s.", m.LifecycleState, strings.Join(GetVaultSummaryLifecycleStateEnumStringValues(), ",")))
	}
	if _, ok := GetMappingVaultSummaryVaultTypeEnum(string(m.VaultType)); !ok && m.VaultType != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for VaultType: %s. Supported values are: %s.", m.VaultType, strings.Join(GetVaultSummaryVaultTypeEnumStringValues(), ",")))
	}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
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
	VaultSummaryLifecycleStateBackupInProgress   VaultSummaryLifecycleStateEnum = "BACKUP_IN_PROGRESS"
	VaultSummaryLifecycleStateRestoring          VaultSummaryLifecycleStateEnum = "RESTORING"
)

var mappingVaultSummaryLifecycleStateEnum = map[string]VaultSummaryLifecycleStateEnum{
	"CREATING":            VaultSummaryLifecycleStateCreating,
	"ACTIVE":              VaultSummaryLifecycleStateActive,
	"DELETING":            VaultSummaryLifecycleStateDeleting,
	"DELETED":             VaultSummaryLifecycleStateDeleted,
	"PENDING_DELETION":    VaultSummaryLifecycleStatePendingDeletion,
	"SCHEDULING_DELETION": VaultSummaryLifecycleStateSchedulingDeletion,
	"CANCELLING_DELETION": VaultSummaryLifecycleStateCancellingDeletion,
	"UPDATING":            VaultSummaryLifecycleStateUpdating,
	"BACKUP_IN_PROGRESS":  VaultSummaryLifecycleStateBackupInProgress,
	"RESTORING":           VaultSummaryLifecycleStateRestoring,
}

var mappingVaultSummaryLifecycleStateEnumLowerCase = map[string]VaultSummaryLifecycleStateEnum{
	"creating":            VaultSummaryLifecycleStateCreating,
	"active":              VaultSummaryLifecycleStateActive,
	"deleting":            VaultSummaryLifecycleStateDeleting,
	"deleted":             VaultSummaryLifecycleStateDeleted,
	"pending_deletion":    VaultSummaryLifecycleStatePendingDeletion,
	"scheduling_deletion": VaultSummaryLifecycleStateSchedulingDeletion,
	"cancelling_deletion": VaultSummaryLifecycleStateCancellingDeletion,
	"updating":            VaultSummaryLifecycleStateUpdating,
	"backup_in_progress":  VaultSummaryLifecycleStateBackupInProgress,
	"restoring":           VaultSummaryLifecycleStateRestoring,
}

// GetVaultSummaryLifecycleStateEnumValues Enumerates the set of values for VaultSummaryLifecycleStateEnum
func GetVaultSummaryLifecycleStateEnumValues() []VaultSummaryLifecycleStateEnum {
	values := make([]VaultSummaryLifecycleStateEnum, 0)
	for _, v := range mappingVaultSummaryLifecycleStateEnum {
		values = append(values, v)
	}
	return values
}

// GetVaultSummaryLifecycleStateEnumStringValues Enumerates the set of values in String for VaultSummaryLifecycleStateEnum
func GetVaultSummaryLifecycleStateEnumStringValues() []string {
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

// GetMappingVaultSummaryLifecycleStateEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingVaultSummaryLifecycleStateEnum(val string) (VaultSummaryLifecycleStateEnum, bool) {
	enum, ok := mappingVaultSummaryLifecycleStateEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}

// VaultSummaryVaultTypeEnum Enum with underlying type: string
type VaultSummaryVaultTypeEnum string

// Set of constants representing the allowable values for VaultSummaryVaultTypeEnum
const (
	VaultSummaryVaultTypeVirtualPrivate VaultSummaryVaultTypeEnum = "VIRTUAL_PRIVATE"
	VaultSummaryVaultTypeDefault        VaultSummaryVaultTypeEnum = "DEFAULT"
)

var mappingVaultSummaryVaultTypeEnum = map[string]VaultSummaryVaultTypeEnum{
	"VIRTUAL_PRIVATE": VaultSummaryVaultTypeVirtualPrivate,
	"DEFAULT":         VaultSummaryVaultTypeDefault,
}

var mappingVaultSummaryVaultTypeEnumLowerCase = map[string]VaultSummaryVaultTypeEnum{
	"virtual_private": VaultSummaryVaultTypeVirtualPrivate,
	"default":         VaultSummaryVaultTypeDefault,
}

// GetVaultSummaryVaultTypeEnumValues Enumerates the set of values for VaultSummaryVaultTypeEnum
func GetVaultSummaryVaultTypeEnumValues() []VaultSummaryVaultTypeEnum {
	values := make([]VaultSummaryVaultTypeEnum, 0)
	for _, v := range mappingVaultSummaryVaultTypeEnum {
		values = append(values, v)
	}
	return values
}

// GetVaultSummaryVaultTypeEnumStringValues Enumerates the set of values in String for VaultSummaryVaultTypeEnum
func GetVaultSummaryVaultTypeEnumStringValues() []string {
	return []string{
		"VIRTUAL_PRIVATE",
		"DEFAULT",
	}
}

// GetMappingVaultSummaryVaultTypeEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingVaultSummaryVaultTypeEnum(val string) (VaultSummaryVaultTypeEnum, bool) {
	enum, ok := mappingVaultSummaryVaultTypeEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
