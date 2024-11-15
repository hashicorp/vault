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

// CreateKeyDetails The representation of CreateKeyDetails
type CreateKeyDetails struct {

	// The OCID of the compartment where you want to create the master encryption key.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name for the key. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	KeyShape *KeyShape `mandatory:"true" json:"keyShape"`

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
	ProtectionMode CreateKeyDetailsProtectionModeEnum `mandatory:"false" json:"protectionMode,omitempty"`
}

func (m CreateKeyDetails) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m CreateKeyDetails) ValidateEnumValue() (bool, error) {
	errMessage := []string{}

	if _, ok := GetMappingCreateKeyDetailsProtectionModeEnum(string(m.ProtectionMode)); !ok && m.ProtectionMode != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for ProtectionMode: %s. Supported values are: %s.", m.ProtectionMode, strings.Join(GetCreateKeyDetailsProtectionModeEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// CreateKeyDetailsProtectionModeEnum Enum with underlying type: string
type CreateKeyDetailsProtectionModeEnum string

// Set of constants representing the allowable values for CreateKeyDetailsProtectionModeEnum
const (
	CreateKeyDetailsProtectionModeHsm      CreateKeyDetailsProtectionModeEnum = "HSM"
	CreateKeyDetailsProtectionModeSoftware CreateKeyDetailsProtectionModeEnum = "SOFTWARE"
)

var mappingCreateKeyDetailsProtectionModeEnum = map[string]CreateKeyDetailsProtectionModeEnum{
	"HSM":      CreateKeyDetailsProtectionModeHsm,
	"SOFTWARE": CreateKeyDetailsProtectionModeSoftware,
}

var mappingCreateKeyDetailsProtectionModeEnumLowerCase = map[string]CreateKeyDetailsProtectionModeEnum{
	"hsm":      CreateKeyDetailsProtectionModeHsm,
	"software": CreateKeyDetailsProtectionModeSoftware,
}

// GetCreateKeyDetailsProtectionModeEnumValues Enumerates the set of values for CreateKeyDetailsProtectionModeEnum
func GetCreateKeyDetailsProtectionModeEnumValues() []CreateKeyDetailsProtectionModeEnum {
	values := make([]CreateKeyDetailsProtectionModeEnum, 0)
	for _, v := range mappingCreateKeyDetailsProtectionModeEnum {
		values = append(values, v)
	}
	return values
}

// GetCreateKeyDetailsProtectionModeEnumStringValues Enumerates the set of values in String for CreateKeyDetailsProtectionModeEnum
func GetCreateKeyDetailsProtectionModeEnumStringValues() []string {
	return []string{
		"HSM",
		"SOFTWARE",
	}
}

// GetMappingCreateKeyDetailsProtectionModeEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingCreateKeyDetailsProtectionModeEnum(val string) (CreateKeyDetailsProtectionModeEnum, bool) {
	enum, ok := mappingCreateKeyDetailsProtectionModeEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
