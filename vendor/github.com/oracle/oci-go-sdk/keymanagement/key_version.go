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

// KeyVersion The representation of KeyVersion
type KeyVersion struct {

	// The OCID of the compartment that contains this key version.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the key version.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the key associated with this key version.
	KeyId *string `mandatory:"true" json:"keyId"`

	// The date and time this key version was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: "2018-04-03T21:10:29.600Z"
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the vault that contains this key version.
	VaultId *string `mandatory:"true" json:"vaultId"`

	// The key version's current lifecycle state.
	// Example: `ENABLED`
	LifecycleState KeyVersionLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The source of the key material. When this value is `INTERNAL`, Key Management
	// created the key material. When this value is `EXTERNAL`, the key material
	// was imported from an external source.
	Origin KeyVersionOriginEnum `mandatory:"false" json:"origin,omitempty"`

	// An optional property indicating when to delete the key version, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2019-04-03T21:10:29.600Z`
	TimeOfDeletion *common.SDKTime `mandatory:"false" json:"timeOfDeletion"`

	// The OCID of the key version from which this key version was restored.
	RestoredFromKeyVersionId *string `mandatory:"false" json:"restoredFromKeyVersionId"`
}

func (m KeyVersion) String() string {
	return common.PointerString(m)
}

// KeyVersionLifecycleStateEnum Enum with underlying type: string
type KeyVersionLifecycleStateEnum string

// Set of constants representing the allowable values for KeyVersionLifecycleStateEnum
const (
	KeyVersionLifecycleStateCreating           KeyVersionLifecycleStateEnum = "CREATING"
	KeyVersionLifecycleStateEnabling           KeyVersionLifecycleStateEnum = "ENABLING"
	KeyVersionLifecycleStateEnabled            KeyVersionLifecycleStateEnum = "ENABLED"
	KeyVersionLifecycleStateDisabling          KeyVersionLifecycleStateEnum = "DISABLING"
	KeyVersionLifecycleStateDisabled           KeyVersionLifecycleStateEnum = "DISABLED"
	KeyVersionLifecycleStateDeleting           KeyVersionLifecycleStateEnum = "DELETING"
	KeyVersionLifecycleStateDeleted            KeyVersionLifecycleStateEnum = "DELETED"
	KeyVersionLifecycleStatePendingDeletion    KeyVersionLifecycleStateEnum = "PENDING_DELETION"
	KeyVersionLifecycleStateSchedulingDeletion KeyVersionLifecycleStateEnum = "SCHEDULING_DELETION"
	KeyVersionLifecycleStateCancellingDeletion KeyVersionLifecycleStateEnum = "CANCELLING_DELETION"
)

var mappingKeyVersionLifecycleState = map[string]KeyVersionLifecycleStateEnum{
	"CREATING":            KeyVersionLifecycleStateCreating,
	"ENABLING":            KeyVersionLifecycleStateEnabling,
	"ENABLED":             KeyVersionLifecycleStateEnabled,
	"DISABLING":           KeyVersionLifecycleStateDisabling,
	"DISABLED":            KeyVersionLifecycleStateDisabled,
	"DELETING":            KeyVersionLifecycleStateDeleting,
	"DELETED":             KeyVersionLifecycleStateDeleted,
	"PENDING_DELETION":    KeyVersionLifecycleStatePendingDeletion,
	"SCHEDULING_DELETION": KeyVersionLifecycleStateSchedulingDeletion,
	"CANCELLING_DELETION": KeyVersionLifecycleStateCancellingDeletion,
}

// GetKeyVersionLifecycleStateEnumValues Enumerates the set of values for KeyVersionLifecycleStateEnum
func GetKeyVersionLifecycleStateEnumValues() []KeyVersionLifecycleStateEnum {
	values := make([]KeyVersionLifecycleStateEnum, 0)
	for _, v := range mappingKeyVersionLifecycleState {
		values = append(values, v)
	}
	return values
}

// KeyVersionOriginEnum Enum with underlying type: string
type KeyVersionOriginEnum string

// Set of constants representing the allowable values for KeyVersionOriginEnum
const (
	KeyVersionOriginInternal KeyVersionOriginEnum = "INTERNAL"
	KeyVersionOriginExternal KeyVersionOriginEnum = "EXTERNAL"
)

var mappingKeyVersionOrigin = map[string]KeyVersionOriginEnum{
	"INTERNAL": KeyVersionOriginInternal,
	"EXTERNAL": KeyVersionOriginExternal,
}

// GetKeyVersionOriginEnumValues Enumerates the set of values for KeyVersionOriginEnum
func GetKeyVersionOriginEnumValues() []KeyVersionOriginEnum {
	values := make([]KeyVersionOriginEnum, 0)
	for _, v := range mappingKeyVersionOrigin {
		values = append(values, v)
	}
	return values
}
