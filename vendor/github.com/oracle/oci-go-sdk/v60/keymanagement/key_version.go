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

	// The public key in PEM format. (This value pertains only to RSA and ECDSA keys.)
	PublicKey *string `mandatory:"false" json:"publicKey"`

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

	ReplicaDetails *KeyVersionReplicaDetails `mandatory:"false" json:"replicaDetails"`

	IsPrimary *bool `mandatory:"false" json:"isPrimary"`
}

func (m KeyVersion) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m KeyVersion) ValidateEnumValue() (bool, error) {
	errMessage := []string{}

	if _, ok := GetMappingKeyVersionLifecycleStateEnum(string(m.LifecycleState)); !ok && m.LifecycleState != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for LifecycleState: %s. Supported values are: %s.", m.LifecycleState, strings.Join(GetKeyVersionLifecycleStateEnumStringValues(), ",")))
	}
	if _, ok := GetMappingKeyVersionOriginEnum(string(m.Origin)); !ok && m.Origin != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for Origin: %s. Supported values are: %s.", m.Origin, strings.Join(GetKeyVersionOriginEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
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

var mappingKeyVersionLifecycleStateEnum = map[string]KeyVersionLifecycleStateEnum{
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

var mappingKeyVersionLifecycleStateEnumLowerCase = map[string]KeyVersionLifecycleStateEnum{
	"creating":            KeyVersionLifecycleStateCreating,
	"enabling":            KeyVersionLifecycleStateEnabling,
	"enabled":             KeyVersionLifecycleStateEnabled,
	"disabling":           KeyVersionLifecycleStateDisabling,
	"disabled":            KeyVersionLifecycleStateDisabled,
	"deleting":            KeyVersionLifecycleStateDeleting,
	"deleted":             KeyVersionLifecycleStateDeleted,
	"pending_deletion":    KeyVersionLifecycleStatePendingDeletion,
	"scheduling_deletion": KeyVersionLifecycleStateSchedulingDeletion,
	"cancelling_deletion": KeyVersionLifecycleStateCancellingDeletion,
}

// GetKeyVersionLifecycleStateEnumValues Enumerates the set of values for KeyVersionLifecycleStateEnum
func GetKeyVersionLifecycleStateEnumValues() []KeyVersionLifecycleStateEnum {
	values := make([]KeyVersionLifecycleStateEnum, 0)
	for _, v := range mappingKeyVersionLifecycleStateEnum {
		values = append(values, v)
	}
	return values
}

// GetKeyVersionLifecycleStateEnumStringValues Enumerates the set of values in String for KeyVersionLifecycleStateEnum
func GetKeyVersionLifecycleStateEnumStringValues() []string {
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
	}
}

// GetMappingKeyVersionLifecycleStateEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingKeyVersionLifecycleStateEnum(val string) (KeyVersionLifecycleStateEnum, bool) {
	enum, ok := mappingKeyVersionLifecycleStateEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}

// KeyVersionOriginEnum Enum with underlying type: string
type KeyVersionOriginEnum string

// Set of constants representing the allowable values for KeyVersionOriginEnum
const (
	KeyVersionOriginInternal KeyVersionOriginEnum = "INTERNAL"
	KeyVersionOriginExternal KeyVersionOriginEnum = "EXTERNAL"
)

var mappingKeyVersionOriginEnum = map[string]KeyVersionOriginEnum{
	"INTERNAL": KeyVersionOriginInternal,
	"EXTERNAL": KeyVersionOriginExternal,
}

var mappingKeyVersionOriginEnumLowerCase = map[string]KeyVersionOriginEnum{
	"internal": KeyVersionOriginInternal,
	"external": KeyVersionOriginExternal,
}

// GetKeyVersionOriginEnumValues Enumerates the set of values for KeyVersionOriginEnum
func GetKeyVersionOriginEnumValues() []KeyVersionOriginEnum {
	values := make([]KeyVersionOriginEnum, 0)
	for _, v := range mappingKeyVersionOriginEnum {
		values = append(values, v)
	}
	return values
}

// GetKeyVersionOriginEnumStringValues Enumerates the set of values in String for KeyVersionOriginEnum
func GetKeyVersionOriginEnumStringValues() []string {
	return []string{
		"INTERNAL",
		"EXTERNAL",
	}
}

// GetMappingKeyVersionOriginEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingKeyVersionOriginEnum(val string) (KeyVersionOriginEnum, bool) {
	enum, ok := mappingKeyVersionOriginEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
