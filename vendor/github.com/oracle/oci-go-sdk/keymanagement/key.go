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

	// An optional property indicating when to delete the key, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2019-04-03T21:10:29.600Z`
	TimeOfDeletion *common.SDKTime `mandatory:"false" json:"timeOfDeletion"`

	// The OCID of the key from which this key was restored.
	RestoredFromKeyId *string `mandatory:"false" json:"restoredFromKeyId"`
}

func (m Key) String() string {
	return common.PointerString(m)
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

var mappingKeyLifecycleState = map[string]KeyLifecycleStateEnum{
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

// GetKeyLifecycleStateEnumValues Enumerates the set of values for KeyLifecycleStateEnum
func GetKeyLifecycleStateEnumValues() []KeyLifecycleStateEnum {
	values := make([]KeyLifecycleStateEnum, 0)
	for _, v := range mappingKeyLifecycleState {
		values = append(values, v)
	}
	return values
}
