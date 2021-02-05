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

// WrappingKey The representation of WrappingKey
type WrappingKey struct {

	// The OCID of the compartment that contains this key.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the key.
	Id *string `mandatory:"true" json:"id"`

	// The key's current lifecycle state.
	// Example: `ENABLED`
	LifecycleState WrappingKeyLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The public key in PEM format to encrypt the key material before importing it with ImportKey/ImportKeyVersion.
	PublicKey *string `mandatory:"true" json:"publicKey"`

	// The date and time the key was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2018-04-03T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the vault that contains this key.
	VaultId *string `mandatory:"true" json:"vaultId"`
}

func (m WrappingKey) String() string {
	return common.PointerString(m)
}

// WrappingKeyLifecycleStateEnum Enum with underlying type: string
type WrappingKeyLifecycleStateEnum string

// Set of constants representing the allowable values for WrappingKeyLifecycleStateEnum
const (
	WrappingKeyLifecycleStateCreating           WrappingKeyLifecycleStateEnum = "CREATING"
	WrappingKeyLifecycleStateEnabling           WrappingKeyLifecycleStateEnum = "ENABLING"
	WrappingKeyLifecycleStateEnabled            WrappingKeyLifecycleStateEnum = "ENABLED"
	WrappingKeyLifecycleStateDisabling          WrappingKeyLifecycleStateEnum = "DISABLING"
	WrappingKeyLifecycleStateDisabled           WrappingKeyLifecycleStateEnum = "DISABLED"
	WrappingKeyLifecycleStateDeleting           WrappingKeyLifecycleStateEnum = "DELETING"
	WrappingKeyLifecycleStateDeleted            WrappingKeyLifecycleStateEnum = "DELETED"
	WrappingKeyLifecycleStatePendingDeletion    WrappingKeyLifecycleStateEnum = "PENDING_DELETION"
	WrappingKeyLifecycleStateSchedulingDeletion WrappingKeyLifecycleStateEnum = "SCHEDULING_DELETION"
	WrappingKeyLifecycleStateCancellingDeletion WrappingKeyLifecycleStateEnum = "CANCELLING_DELETION"
	WrappingKeyLifecycleStateUpdating           WrappingKeyLifecycleStateEnum = "UPDATING"
	WrappingKeyLifecycleStateBackupInProgress   WrappingKeyLifecycleStateEnum = "BACKUP_IN_PROGRESS"
	WrappingKeyLifecycleStateRestoring          WrappingKeyLifecycleStateEnum = "RESTORING"
)

var mappingWrappingKeyLifecycleState = map[string]WrappingKeyLifecycleStateEnum{
	"CREATING":            WrappingKeyLifecycleStateCreating,
	"ENABLING":            WrappingKeyLifecycleStateEnabling,
	"ENABLED":             WrappingKeyLifecycleStateEnabled,
	"DISABLING":           WrappingKeyLifecycleStateDisabling,
	"DISABLED":            WrappingKeyLifecycleStateDisabled,
	"DELETING":            WrappingKeyLifecycleStateDeleting,
	"DELETED":             WrappingKeyLifecycleStateDeleted,
	"PENDING_DELETION":    WrappingKeyLifecycleStatePendingDeletion,
	"SCHEDULING_DELETION": WrappingKeyLifecycleStateSchedulingDeletion,
	"CANCELLING_DELETION": WrappingKeyLifecycleStateCancellingDeletion,
	"UPDATING":            WrappingKeyLifecycleStateUpdating,
	"BACKUP_IN_PROGRESS":  WrappingKeyLifecycleStateBackupInProgress,
	"RESTORING":           WrappingKeyLifecycleStateRestoring,
}

// GetWrappingKeyLifecycleStateEnumValues Enumerates the set of values for WrappingKeyLifecycleStateEnum
func GetWrappingKeyLifecycleStateEnumValues() []WrappingKeyLifecycleStateEnum {
	values := make([]WrappingKeyLifecycleStateEnum, 0)
	for _, v := range mappingWrappingKeyLifecycleState {
		values = append(values, v)
	}
	return values
}
