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

// CreateVaultDetails The representation of CreateVaultDetails
type CreateVaultDetails struct {

	// The OCID of the compartment where you want to create this vault.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// A user-friendly name for the vault. It does not have to be unique, and it is changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The type of vault to create. Each type of vault stores the key with different degrees of isolation and has different options and pricing.
	VaultType CreateVaultDetailsVaultTypeEnum `mandatory:"true" json:"vaultType"`

	// Usage of predefined tag keys. These predefined keys are scoped to namespaces.
	// Example: `{"foo-namespace": {"bar-key": "foo-value"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Simple key-value pair that is applied without any predefined name, type, or scope.
	// Exists for cross-compatibility only.
	// Example: `{"bar-key": "value"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m CreateVaultDetails) String() string {
	return common.PointerString(m)
}

// CreateVaultDetailsVaultTypeEnum Enum with underlying type: string
type CreateVaultDetailsVaultTypeEnum string

// Set of constants representing the allowable values for CreateVaultDetailsVaultTypeEnum
const (
	CreateVaultDetailsVaultTypeVirtualPrivate CreateVaultDetailsVaultTypeEnum = "VIRTUAL_PRIVATE"
)

var mappingCreateVaultDetailsVaultType = map[string]CreateVaultDetailsVaultTypeEnum{
	"VIRTUAL_PRIVATE": CreateVaultDetailsVaultTypeVirtualPrivate,
}

// GetCreateVaultDetailsVaultTypeEnumValues Enumerates the set of values for CreateVaultDetailsVaultTypeEnum
func GetCreateVaultDetailsVaultTypeEnumValues() []CreateVaultDetailsVaultTypeEnum {
	values := make([]CreateVaultDetailsVaultTypeEnum, 0)
	for _, v := range mappingCreateVaultDetailsVaultType {
		values = append(values, v)
	}
	return values
}
