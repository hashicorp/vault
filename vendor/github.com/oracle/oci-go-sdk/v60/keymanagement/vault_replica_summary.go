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

// VaultReplicaSummary Summary of vault replicas
type VaultReplicaSummary struct {

	// The vault replica's crypto endpoint
	CryptoEndpoint *string `mandatory:"false" json:"cryptoEndpoint"`

	// The vault replica's management endpoint
	ManagementEndpoint *string `mandatory:"false" json:"managementEndpoint"`

	// Region to which vault is replicated to
	Region *string `mandatory:"false" json:"region"`

	Status VaultReplicaSummaryStatusEnum `mandatory:"false" json:"status,omitempty"`
}

func (m VaultReplicaSummary) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m VaultReplicaSummary) ValidateEnumValue() (bool, error) {
	errMessage := []string{}

	if _, ok := GetMappingVaultReplicaSummaryStatusEnum(string(m.Status)); !ok && m.Status != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for Status: %s. Supported values are: %s.", m.Status, strings.Join(GetVaultReplicaSummaryStatusEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// VaultReplicaSummaryStatusEnum Enum with underlying type: string
type VaultReplicaSummaryStatusEnum string

// Set of constants representing the allowable values for VaultReplicaSummaryStatusEnum
const (
	VaultReplicaSummaryStatusCreating VaultReplicaSummaryStatusEnum = "CREATING"
	VaultReplicaSummaryStatusCreated  VaultReplicaSummaryStatusEnum = "CREATED"
	VaultReplicaSummaryStatusDeleting VaultReplicaSummaryStatusEnum = "DELETING"
	VaultReplicaSummaryStatusDeleted  VaultReplicaSummaryStatusEnum = "DELETED"
)

var mappingVaultReplicaSummaryStatusEnum = map[string]VaultReplicaSummaryStatusEnum{
	"CREATING": VaultReplicaSummaryStatusCreating,
	"CREATED":  VaultReplicaSummaryStatusCreated,
	"DELETING": VaultReplicaSummaryStatusDeleting,
	"DELETED":  VaultReplicaSummaryStatusDeleted,
}

var mappingVaultReplicaSummaryStatusEnumLowerCase = map[string]VaultReplicaSummaryStatusEnum{
	"creating": VaultReplicaSummaryStatusCreating,
	"created":  VaultReplicaSummaryStatusCreated,
	"deleting": VaultReplicaSummaryStatusDeleting,
	"deleted":  VaultReplicaSummaryStatusDeleted,
}

// GetVaultReplicaSummaryStatusEnumValues Enumerates the set of values for VaultReplicaSummaryStatusEnum
func GetVaultReplicaSummaryStatusEnumValues() []VaultReplicaSummaryStatusEnum {
	values := make([]VaultReplicaSummaryStatusEnum, 0)
	for _, v := range mappingVaultReplicaSummaryStatusEnum {
		values = append(values, v)
	}
	return values
}

// GetVaultReplicaSummaryStatusEnumStringValues Enumerates the set of values in String for VaultReplicaSummaryStatusEnum
func GetVaultReplicaSummaryStatusEnumStringValues() []string {
	return []string{
		"CREATING",
		"CREATED",
		"DELETING",
		"DELETED",
	}
}

// GetMappingVaultReplicaSummaryStatusEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingVaultReplicaSummaryStatusEnum(val string) (VaultReplicaSummaryStatusEnum, bool) {
	enum, ok := mappingVaultReplicaSummaryStatusEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
