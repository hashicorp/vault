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

// KeySummary The representation of KeySummary
type KeySummary struct {

	// The OCID of the compartment that contains the key.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name for the key. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID of the key.
	Id *string `mandatory:"true" json:"id"`

	// The key's current lifecycle state.
	// Example: `ENABLED`
	LifecycleState KeySummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the key was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2018-04-03T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the vault that contains the key.
	VaultId *string `mandatory:"true" json:"vaultId"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m KeySummary) String() string {
	return common.PointerString(m)
}

// KeySummaryLifecycleStateEnum Enum with underlying type: string
type KeySummaryLifecycleStateEnum string

// Set of constants representing the allowable values for KeySummaryLifecycleStateEnum
const (
	KeySummaryLifecycleStateCreating           KeySummaryLifecycleStateEnum = "CREATING"
	KeySummaryLifecycleStateEnabling           KeySummaryLifecycleStateEnum = "ENABLING"
	KeySummaryLifecycleStateEnabled            KeySummaryLifecycleStateEnum = "ENABLED"
	KeySummaryLifecycleStateDisabling          KeySummaryLifecycleStateEnum = "DISABLING"
	KeySummaryLifecycleStateDisabled           KeySummaryLifecycleStateEnum = "DISABLED"
	KeySummaryLifecycleStateDeleting           KeySummaryLifecycleStateEnum = "DELETING"
	KeySummaryLifecycleStateDeleted            KeySummaryLifecycleStateEnum = "DELETED"
	KeySummaryLifecycleStatePendingDeletion    KeySummaryLifecycleStateEnum = "PENDING_DELETION"
	KeySummaryLifecycleStateSchedulingDeletion KeySummaryLifecycleStateEnum = "SCHEDULING_DELETION"
	KeySummaryLifecycleStateCancellingDeletion KeySummaryLifecycleStateEnum = "CANCELLING_DELETION"
	KeySummaryLifecycleStateUpdating           KeySummaryLifecycleStateEnum = "UPDATING"
	KeySummaryLifecycleStateBackupInProgress   KeySummaryLifecycleStateEnum = "BACKUP_IN_PROGRESS"
	KeySummaryLifecycleStateRestoring          KeySummaryLifecycleStateEnum = "RESTORING"
)

var mappingKeySummaryLifecycleState = map[string]KeySummaryLifecycleStateEnum{
	"CREATING":            KeySummaryLifecycleStateCreating,
	"ENABLING":            KeySummaryLifecycleStateEnabling,
	"ENABLED":             KeySummaryLifecycleStateEnabled,
	"DISABLING":           KeySummaryLifecycleStateDisabling,
	"DISABLED":            KeySummaryLifecycleStateDisabled,
	"DELETING":            KeySummaryLifecycleStateDeleting,
	"DELETED":             KeySummaryLifecycleStateDeleted,
	"PENDING_DELETION":    KeySummaryLifecycleStatePendingDeletion,
	"SCHEDULING_DELETION": KeySummaryLifecycleStateSchedulingDeletion,
	"CANCELLING_DELETION": KeySummaryLifecycleStateCancellingDeletion,
	"UPDATING":            KeySummaryLifecycleStateUpdating,
	"BACKUP_IN_PROGRESS":  KeySummaryLifecycleStateBackupInProgress,
	"RESTORING":           KeySummaryLifecycleStateRestoring,
}

// GetKeySummaryLifecycleStateEnumValues Enumerates the set of values for KeySummaryLifecycleStateEnum
func GetKeySummaryLifecycleStateEnumValues() []KeySummaryLifecycleStateEnum {
	values := make([]KeySummaryLifecycleStateEnum, 0)
	for _, v := range mappingKeySummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
