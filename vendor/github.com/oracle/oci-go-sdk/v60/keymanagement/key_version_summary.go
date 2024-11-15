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

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m KeyVersionSummary) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingKeyVersionSummaryOriginEnum(string(m.Origin)); !ok && m.Origin != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for Origin: %s. Supported values are: %s.", m.Origin, strings.Join(GetKeyVersionSummaryOriginEnumStringValues(), ",")))
	}

	if _, ok := GetMappingKeyVersionSummaryLifecycleStateEnum(string(m.LifecycleState)); !ok && m.LifecycleState != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for LifecycleState: %s. Supported values are: %s.", m.LifecycleState, strings.Join(GetKeyVersionSummaryLifecycleStateEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
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

var mappingKeyVersionSummaryLifecycleStateEnum = map[string]KeyVersionSummaryLifecycleStateEnum{
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

var mappingKeyVersionSummaryLifecycleStateEnumLowerCase = map[string]KeyVersionSummaryLifecycleStateEnum{
	"creating":            KeyVersionSummaryLifecycleStateCreating,
	"enabling":            KeyVersionSummaryLifecycleStateEnabling,
	"enabled":             KeyVersionSummaryLifecycleStateEnabled,
	"disabling":           KeyVersionSummaryLifecycleStateDisabling,
	"disabled":            KeyVersionSummaryLifecycleStateDisabled,
	"deleting":            KeyVersionSummaryLifecycleStateDeleting,
	"deleted":             KeyVersionSummaryLifecycleStateDeleted,
	"pending_deletion":    KeyVersionSummaryLifecycleStatePendingDeletion,
	"scheduling_deletion": KeyVersionSummaryLifecycleStateSchedulingDeletion,
	"cancelling_deletion": KeyVersionSummaryLifecycleStateCancellingDeletion,
}

// GetKeyVersionSummaryLifecycleStateEnumValues Enumerates the set of values for KeyVersionSummaryLifecycleStateEnum
func GetKeyVersionSummaryLifecycleStateEnumValues() []KeyVersionSummaryLifecycleStateEnum {
	values := make([]KeyVersionSummaryLifecycleStateEnum, 0)
	for _, v := range mappingKeyVersionSummaryLifecycleStateEnum {
		values = append(values, v)
	}
	return values
}

// GetKeyVersionSummaryLifecycleStateEnumStringValues Enumerates the set of values in String for KeyVersionSummaryLifecycleStateEnum
func GetKeyVersionSummaryLifecycleStateEnumStringValues() []string {
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

// GetMappingKeyVersionSummaryLifecycleStateEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingKeyVersionSummaryLifecycleStateEnum(val string) (KeyVersionSummaryLifecycleStateEnum, bool) {
	enum, ok := mappingKeyVersionSummaryLifecycleStateEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}

// KeyVersionSummaryOriginEnum Enum with underlying type: string
type KeyVersionSummaryOriginEnum string

// Set of constants representing the allowable values for KeyVersionSummaryOriginEnum
const (
	KeyVersionSummaryOriginInternal KeyVersionSummaryOriginEnum = "INTERNAL"
	KeyVersionSummaryOriginExternal KeyVersionSummaryOriginEnum = "EXTERNAL"
)

var mappingKeyVersionSummaryOriginEnum = map[string]KeyVersionSummaryOriginEnum{
	"INTERNAL": KeyVersionSummaryOriginInternal,
	"EXTERNAL": KeyVersionSummaryOriginExternal,
}

var mappingKeyVersionSummaryOriginEnumLowerCase = map[string]KeyVersionSummaryOriginEnum{
	"internal": KeyVersionSummaryOriginInternal,
	"external": KeyVersionSummaryOriginExternal,
}

// GetKeyVersionSummaryOriginEnumValues Enumerates the set of values for KeyVersionSummaryOriginEnum
func GetKeyVersionSummaryOriginEnumValues() []KeyVersionSummaryOriginEnum {
	values := make([]KeyVersionSummaryOriginEnum, 0)
	for _, v := range mappingKeyVersionSummaryOriginEnum {
		values = append(values, v)
	}
	return values
}

// GetKeyVersionSummaryOriginEnumStringValues Enumerates the set of values in String for KeyVersionSummaryOriginEnum
func GetKeyVersionSummaryOriginEnumStringValues() []string {
	return []string{
		"INTERNAL",
		"EXTERNAL",
	}
}

// GetMappingKeyVersionSummaryOriginEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingKeyVersionSummaryOriginEnum(val string) (KeyVersionSummaryOriginEnum, bool) {
	enum, ok := mappingKeyVersionSummaryOriginEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
