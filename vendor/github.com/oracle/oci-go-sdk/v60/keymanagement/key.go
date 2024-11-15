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

// Key The representation of Key
type Key struct {

	// The OCID of the compartment that contains this master encryption key.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the key version used in cryptographic operations. During key rotation, the service might be
	// in a transitional state where this or a newer key version are used intermittently. The `currentKeyVersion`
	// property is updated when the service is guaranteed to use the new key version for all subsequent encryption operations.
	CurrentKeyVersion *string `mandatory:"true" json:"currentKeyVersion"`

	// A user-friendly name for the key. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID of the key.
	Id *string `mandatory:"true" json:"id"`

	KeyShape *KeyShape `mandatory:"true" json:"keyShape"`

	// The key's current lifecycle state.
	// Example: `ENABLED`
	LifecycleState KeyLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the key was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2018-04-03T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the vault that contains this key.
	VaultId *string `mandatory:"true" json:"vaultId"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The key's protection mode indicates how the key persists and where cryptographic operations that use the key are performed.
	// A protection mode of `HSM` means that the key persists on a hardware security module (HSM) and all cryptographic operations are performed inside
	// the HSM. A protection mode of `SOFTWARE` means that the key persists on the server, protected by the vault's RSA wrapping key which persists
	// on the HSM. All cryptographic operations that use a key with a protection mode of `SOFTWARE` are performed on the server. By default,
	// a key's protection mode is set to `HSM`. You can't change a key's protection mode after the key is created or imported.
	ProtectionMode KeyProtectionModeEnum `mandatory:"false" json:"protectionMode,omitempty"`

	// An optional property indicating when to delete the key, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2019-04-03T21:10:29.600Z`
	TimeOfDeletion *common.SDKTime `mandatory:"false" json:"timeOfDeletion"`

	// The OCID of the key from which this key was restored.
	RestoredFromKeyId *string `mandatory:"false" json:"restoredFromKeyId"`

	ReplicaDetails *KeyReplicaDetails `mandatory:"false" json:"replicaDetails"`

	IsPrimary *bool `mandatory:"false" json:"isPrimary"`
}

func (m Key) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m Key) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingKeyLifecycleStateEnum(string(m.LifecycleState)); !ok && m.LifecycleState != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for LifecycleState: %s. Supported values are: %s.", m.LifecycleState, strings.Join(GetKeyLifecycleStateEnumStringValues(), ",")))
	}

	if _, ok := GetMappingKeyProtectionModeEnum(string(m.ProtectionMode)); !ok && m.ProtectionMode != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for ProtectionMode: %s. Supported values are: %s.", m.ProtectionMode, strings.Join(GetKeyProtectionModeEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// KeyProtectionModeEnum Enum with underlying type: string
type KeyProtectionModeEnum string

// Set of constants representing the allowable values for KeyProtectionModeEnum
const (
	KeyProtectionModeHsm      KeyProtectionModeEnum = "HSM"
	KeyProtectionModeSoftware KeyProtectionModeEnum = "SOFTWARE"
)

var mappingKeyProtectionModeEnum = map[string]KeyProtectionModeEnum{
	"HSM":      KeyProtectionModeHsm,
	"SOFTWARE": KeyProtectionModeSoftware,
}

var mappingKeyProtectionModeEnumLowerCase = map[string]KeyProtectionModeEnum{
	"hsm":      KeyProtectionModeHsm,
	"software": KeyProtectionModeSoftware,
}

// GetKeyProtectionModeEnumValues Enumerates the set of values for KeyProtectionModeEnum
func GetKeyProtectionModeEnumValues() []KeyProtectionModeEnum {
	values := make([]KeyProtectionModeEnum, 0)
	for _, v := range mappingKeyProtectionModeEnum {
		values = append(values, v)
	}
	return values
}

// GetKeyProtectionModeEnumStringValues Enumerates the set of values in String for KeyProtectionModeEnum
func GetKeyProtectionModeEnumStringValues() []string {
	return []string{
		"HSM",
		"SOFTWARE",
	}
}

// GetMappingKeyProtectionModeEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingKeyProtectionModeEnum(val string) (KeyProtectionModeEnum, bool) {
	enum, ok := mappingKeyProtectionModeEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}

// KeyLifecycleStateEnum Enum with underlying type: string
type KeyLifecycleStateEnum string

// Set of constants representing the allowable values for KeyLifecycleStateEnum
const (
	KeyLifecycleStateCreating           KeyLifecycleStateEnum = "CREATING"
	KeyLifecycleStateEnabling           KeyLifecycleStateEnum = "ENABLING"
	KeyLifecycleStateEnabled            KeyLifecycleStateEnum = "ENABLED"
	KeyLifecycleStateDisabling          KeyLifecycleStateEnum = "DISABLING"
	KeyLifecycleStateDisabled           KeyLifecycleStateEnum = "DISABLED"
	KeyLifecycleStateDeleting           KeyLifecycleStateEnum = "DELETING"
	KeyLifecycleStateDeleted            KeyLifecycleStateEnum = "DELETED"
	KeyLifecycleStatePendingDeletion    KeyLifecycleStateEnum = "PENDING_DELETION"
	KeyLifecycleStateSchedulingDeletion KeyLifecycleStateEnum = "SCHEDULING_DELETION"
	KeyLifecycleStateCancellingDeletion KeyLifecycleStateEnum = "CANCELLING_DELETION"
	KeyLifecycleStateUpdating           KeyLifecycleStateEnum = "UPDATING"
	KeyLifecycleStateBackupInProgress   KeyLifecycleStateEnum = "BACKUP_IN_PROGRESS"
	KeyLifecycleStateRestoring          KeyLifecycleStateEnum = "RESTORING"
)

var mappingKeyLifecycleStateEnum = map[string]KeyLifecycleStateEnum{
	"CREATING":            KeyLifecycleStateCreating,
	"ENABLING":            KeyLifecycleStateEnabling,
	"ENABLED":             KeyLifecycleStateEnabled,
	"DISABLING":           KeyLifecycleStateDisabling,
	"DISABLED":            KeyLifecycleStateDisabled,
	"DELETING":            KeyLifecycleStateDeleting,
	"DELETED":             KeyLifecycleStateDeleted,
	"PENDING_DELETION":    KeyLifecycleStatePendingDeletion,
	"SCHEDULING_DELETION": KeyLifecycleStateSchedulingDeletion,
	"CANCELLING_DELETION": KeyLifecycleStateCancellingDeletion,
	"UPDATING":            KeyLifecycleStateUpdating,
	"BACKUP_IN_PROGRESS":  KeyLifecycleStateBackupInProgress,
	"RESTORING":           KeyLifecycleStateRestoring,
}

var mappingKeyLifecycleStateEnumLowerCase = map[string]KeyLifecycleStateEnum{
	"creating":            KeyLifecycleStateCreating,
	"enabling":            KeyLifecycleStateEnabling,
	"enabled":             KeyLifecycleStateEnabled,
	"disabling":           KeyLifecycleStateDisabling,
	"disabled":            KeyLifecycleStateDisabled,
	"deleting":            KeyLifecycleStateDeleting,
	"deleted":             KeyLifecycleStateDeleted,
	"pending_deletion":    KeyLifecycleStatePendingDeletion,
	"scheduling_deletion": KeyLifecycleStateSchedulingDeletion,
	"cancelling_deletion": KeyLifecycleStateCancellingDeletion,
	"updating":            KeyLifecycleStateUpdating,
	"backup_in_progress":  KeyLifecycleStateBackupInProgress,
	"restoring":           KeyLifecycleStateRestoring,
}

// GetKeyLifecycleStateEnumValues Enumerates the set of values for KeyLifecycleStateEnum
func GetKeyLifecycleStateEnumValues() []KeyLifecycleStateEnum {
	values := make([]KeyLifecycleStateEnum, 0)
	for _, v := range mappingKeyLifecycleStateEnum {
		values = append(values, v)
	}
	return values
}

// GetKeyLifecycleStateEnumStringValues Enumerates the set of values in String for KeyLifecycleStateEnum
func GetKeyLifecycleStateEnumStringValues() []string {
	return []string{
		"CREATING",
		"ENABLING",
		"ENABLED",
		"DISABLING",
		"DISABLED",
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

// GetMappingKeyLifecycleStateEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingKeyLifecycleStateEnum(val string) (KeyLifecycleStateEnum, bool) {
	enum, ok := mappingKeyLifecycleStateEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
