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

// KeySummary The representation of KeySummary
type KeySummary struct {

	// The OCID of the compartment that contains the key.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name for the key. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID of the key.
	Id *string `mandatory:"true" json:"id"`

	// The key's current state.
	// Example: `ENABLED`
	LifecycleState KeySummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the key was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2018-04-03T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the vault that contains the key.
	VaultId *string `mandatory:"true" json:"vaultId"`

	// Usage of predefined tag keys. These predefined keys are scoped to namespaces.
	// Example: `{"foo-namespace": {"bar-key": "foo-value"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Simple key-value pair that is applied without any predefined name, type, or scope.
	// Exists for cross-compatibility only.
	// Example: `{"bar-key": "value"}`
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
}

// GetKeySummaryLifecycleStateEnumValues Enumerates the set of values for KeySummaryLifecycleStateEnum
func GetKeySummaryLifecycleStateEnumValues() []KeySummaryLifecycleStateEnum {
	values := make([]KeySummaryLifecycleStateEnum, 0)
	for _, v := range mappingKeySummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
