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

// ReplicaDetails Details of replication status
type ReplicaDetails struct {

	// The replica region
	Region *string `mandatory:"false" json:"region"`

	// Replication status associated with a replicationId
	Status ReplicaDetailsStatusEnum `mandatory:"false" json:"status,omitempty"`
}

func (m ReplicaDetails) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m ReplicaDetails) ValidateEnumValue() (bool, error) {
	errMessage := []string{}

	if _, ok := GetMappingReplicaDetailsStatusEnum(string(m.Status)); !ok && m.Status != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for Status: %s. Supported values are: %s.", m.Status, strings.Join(GetReplicaDetailsStatusEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// ReplicaDetailsStatusEnum Enum with underlying type: string
type ReplicaDetailsStatusEnum string

// Set of constants representing the allowable values for ReplicaDetailsStatusEnum
const (
	ReplicaDetailsStatusReplicating ReplicaDetailsStatusEnum = "REPLICATING"
	ReplicaDetailsStatusReplicated  ReplicaDetailsStatusEnum = "REPLICATED"
)

var mappingReplicaDetailsStatusEnum = map[string]ReplicaDetailsStatusEnum{
	"REPLICATING": ReplicaDetailsStatusReplicating,
	"REPLICATED":  ReplicaDetailsStatusReplicated,
}

var mappingReplicaDetailsStatusEnumLowerCase = map[string]ReplicaDetailsStatusEnum{
	"replicating": ReplicaDetailsStatusReplicating,
	"replicated":  ReplicaDetailsStatusReplicated,
}

// GetReplicaDetailsStatusEnumValues Enumerates the set of values for ReplicaDetailsStatusEnum
func GetReplicaDetailsStatusEnumValues() []ReplicaDetailsStatusEnum {
	values := make([]ReplicaDetailsStatusEnum, 0)
	for _, v := range mappingReplicaDetailsStatusEnum {
		values = append(values, v)
	}
	return values
}

// GetReplicaDetailsStatusEnumStringValues Enumerates the set of values in String for ReplicaDetailsStatusEnum
func GetReplicaDetailsStatusEnumStringValues() []string {
	return []string{
		"REPLICATING",
		"REPLICATED",
	}
}

// GetMappingReplicaDetailsStatusEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingReplicaDetailsStatusEnum(val string) (ReplicaDetailsStatusEnum, bool) {
	enum, ok := mappingReplicaDetailsStatusEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
