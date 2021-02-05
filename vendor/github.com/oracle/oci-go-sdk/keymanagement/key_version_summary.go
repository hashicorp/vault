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

// KeyVersionSummary The representation of KeyVersionSummary
type KeyVersionSummary struct {

	// The OCID of the compartment that contains this key version.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the key version.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the master encryption key associated with this key version.
	KeyId *string `mandatory:"true" json:"keyId"`

	// The source of the key material. When this value is INTERNAL, Key Management created the key material. When this value is EXTERNAL, the key material was imported from an external source.
	Origin KeyVersionSummaryOriginEnum `mandatory:"true" json:"origin"`

	// The date and time this key version was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2018-04-03T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the vault that contains this key version.
	VaultId *string `mandatory:"true" json:"vaultId"`

	// The key version's current lifecycle state.
	// Example: `ENABLED`
	LifecycleState KeyVersionSummaryLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// An optional property to indicate when to delete the key version, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2019-04-03T21:10:29.600Z`
	TimeOfDeletion *common.SDKTime `mandatory:"false" json:"timeOfDeletion"`
}

func (m KeyVersionSummary) String() string {
	return common.PointerString(m)
}

// KeyVersionSummaryLifecycleStateEnum Enum with underlying type: string
type KeyVersionSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for KeyVersionSummaryLifecycleStateEnum
const (
	KeyVersionSummaryLifecycleStateCreating           KeyVersionSummaryLifecycleStateEnum = "CREATING"
	KeyVersionSummaryLifecycleStateEnabling           KeyVersionSummaryLifecycleStateEnum = "ENABLING"
	KeyVersionSummaryLifecycleStateEnabled            KeyVersionSummaryLifecycleStateEnum = "ENABLED"
	KeyVersionSummaryLifecycleStateDisabling          KeyVersionSummaryLifecycleStateEnum = "DISABLING"
	KeyVersionSummaryLifecycleStateDisabled           KeyVersionSummaryLifecycleStateEnum = "DISABLED"
	KeyVersionSummaryLifecycleStateDeleting           KeyVersionSummaryLifecycleStateEnum = "DELETING"
	KeyVersionSummaryLifecycleStateDeleted            KeyVersionSummaryLifecycleStateEnum = "DELETED"
	KeyVersionSummaryLifecycleStatePendingDeletion    KeyVersionSummaryLifecycleStateEnum = "PENDING_DELETION"
	KeyVersionSummaryLifecycleStateSchedulingDeletion KeyVersionSummaryLifecycleStateEnum = "SCHEDULING_DELETION"
	KeyVersionSummaryLifecycleStateCancellingDeletion KeyVersionSummaryLifecycleStateEnum = "CANCELLING_DELETION"
)

var mappingKeyVersionSummaryLifecycleState = map[string]KeyVersionSummaryLifecycleStateEnum{
	"CREATING":            KeyVersionSummaryLifecycleStateCreating,
	"ENABLING":            KeyVersionSummaryLifecycleStateEnabling,
	"ENABLED":             KeyVersionSummaryLifecycleStateEnabled,
	"DISABLING":           KeyVersionSummaryLifecycleStateDisabling,
	"DISABLED":            KeyVersionSummaryLifecycleStateDisabled,
	"DELETING":            KeyVersionSummaryLifecycleStateDeleting,
	"DELETED":             KeyVersionSummaryLifecycleStateDeleted,
	"PENDING_DELETION":    KeyVersionSummaryLifecycleStatePendingDeletion,
	"SCHEDULING_DELETION": KeyVersionSummaryLifecycleStateSchedulingDeletion,
	"CANCELLING_DELETION": KeyVersionSummaryLifecycleStateCancellingDeletion,
}

// GetKeyVersionSummaryLifecycleStateEnumValues Enumerates the set of values for KeyVersionSummaryLifecycleStateEnum
func GetKeyVersionSummaryLifecycleStateEnumValues() []KeyVersionSummaryLifecycleStateEnum {
	values := make([]KeyVersionSummaryLifecycleStateEnum, 0)
	for _, v := range mappingKeyVersionSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}

// KeyVersionSummaryOriginEnum Enum with underlying type: string
type KeyVersionSummaryOriginEnum string

// Set of constants representing the allowable values for KeyVersionSummaryOriginEnum
const (
	KeyVersionSummaryOriginInternal KeyVersionSummaryOriginEnum = "INTERNAL"
	KeyVersionSummaryOriginExternal KeyVersionSummaryOriginEnum = "EXTERNAL"
)

var mappingKeyVersionSummaryOrigin = map[string]KeyVersionSummaryOriginEnum{
	"INTERNAL": KeyVersionSummaryOriginInternal,
	"EXTERNAL": KeyVersionSummaryOriginExternal,
}

// GetKeyVersionSummaryOriginEnumValues Enumerates the set of values for KeyVersionSummaryOriginEnum
func GetKeyVersionSummaryOriginEnumValues() []KeyVersionSummaryOriginEnum {
	values := make([]KeyVersionSummaryOriginEnum, 0)
	for _, v := range mappingKeyVersionSummaryOrigin {
		values = append(values, v)
	}
	return values
}
