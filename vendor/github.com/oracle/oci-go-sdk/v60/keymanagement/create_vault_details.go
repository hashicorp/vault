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

// CreateVaultDetails The representation of CreateVaultDetails
type CreateVaultDetails struct {

	// The OCID of the compartment where you want to create this vault.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name for the vault. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The type of vault to create. Each type of vault stores the key with different degrees of isolation and has different options and pricing.
	VaultType CreateVaultDetailsVaultTypeEnum `mandatory:"true" json:"vaultType"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m CreateVaultDetails) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m CreateVaultDetails) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingCreateVaultDetailsVaultTypeEnum(string(m.VaultType)); !ok && m.VaultType != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for VaultType: %s. Supported values are: %s.", m.VaultType, strings.Join(GetCreateVaultDetailsVaultTypeEnumStringValues(), ",")))
	}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// CreateVaultDetailsVaultTypeEnum Enum with underlying type: string
type CreateVaultDetailsVaultTypeEnum string

// Set of constants representing the allowable values for CreateVaultDetailsVaultTypeEnum
const (
	CreateVaultDetailsVaultTypeVirtualPrivate CreateVaultDetailsVaultTypeEnum = "VIRTUAL_PRIVATE"
	CreateVaultDetailsVaultTypeDefault        CreateVaultDetailsVaultTypeEnum = "DEFAULT"
)

var mappingCreateVaultDetailsVaultTypeEnum = map[string]CreateVaultDetailsVaultTypeEnum{
	"VIRTUAL_PRIVATE": CreateVaultDetailsVaultTypeVirtualPrivate,
	"DEFAULT":         CreateVaultDetailsVaultTypeDefault,
}

var mappingCreateVaultDetailsVaultTypeEnumLowerCase = map[string]CreateVaultDetailsVaultTypeEnum{
	"virtual_private": CreateVaultDetailsVaultTypeVirtualPrivate,
	"default":         CreateVaultDetailsVaultTypeDefault,
}

// GetCreateVaultDetailsVaultTypeEnumValues Enumerates the set of values for CreateVaultDetailsVaultTypeEnum
func GetCreateVaultDetailsVaultTypeEnumValues() []CreateVaultDetailsVaultTypeEnum {
	values := make([]CreateVaultDetailsVaultTypeEnum, 0)
	for _, v := range mappingCreateVaultDetailsVaultTypeEnum {
		values = append(values, v)
	}
	return values
}

// GetCreateVaultDetailsVaultTypeEnumStringValues Enumerates the set of values in String for CreateVaultDetailsVaultTypeEnum
func GetCreateVaultDetailsVaultTypeEnumStringValues() []string {
	return []string{
		"VIRTUAL_PRIVATE",
		"DEFAULT",
	}
}

// GetMappingCreateVaultDetailsVaultTypeEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingCreateVaultDetailsVaultTypeEnum(val string) (CreateVaultDetailsVaultTypeEnum, bool) {
	enum, ok := mappingCreateVaultDetailsVaultTypeEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
