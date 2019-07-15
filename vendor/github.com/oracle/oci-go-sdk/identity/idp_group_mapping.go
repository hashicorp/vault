// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Identity and Access Management Service API
//
// APIs for managing users, groups, compartments, and policies.
//

package identity

import (
	"github.com/oracle/oci-go-sdk/common"
)

// IdpGroupMapping A mapping between a single group defined by the identity provider (IdP) you're federating with
// and a single IAM Service Group in Oracle Cloud Infrastructure.
// For more information about group mappings and what they're for, see
// Identity Providers and Federation (https://docs.cloud.oracle.com/Content/Identity/Concepts/federation.htm).
// A given IdP group can be mapped to zero, one, or multiple IAM Service groups, and vice versa.
// But each `IdPGroupMapping` object is between only a single IdP group and IAM Service group.
// Each `IdPGroupMapping` object has its own OCID.
// **Note:** Any users who are in more than 50 IdP groups cannot be authenticated to use the Oracle
// Cloud Infrastructure Console.
type IdpGroupMapping struct {

	// The OCID of the `IdpGroupMapping`.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the `IdentityProvider` this mapping belongs to.
	IdpId *string `mandatory:"true" json:"idpId"`

	// The name of the IdP group that is mapped to the IAM Service group.
	IdpGroupName *string `mandatory:"true" json:"idpGroupName"`

	// The OCID of the IAM Service group that is mapped to the IdP group.
	GroupId *string `mandatory:"true" json:"groupId"`

	// The OCID of the tenancy containing the `IdentityProvider`.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// Date and time the mapping was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The mapping's current state.  After creating a mapping object, make sure its `lifecycleState` changes
	// from CREATING to ACTIVE before using it.
	LifecycleState IdpGroupMappingLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`
}

func (m IdpGroupMapping) String() string {
	return common.PointerString(m)
}

// IdpGroupMappingLifecycleStateEnum Enum with underlying type: string
type IdpGroupMappingLifecycleStateEnum string

// Set of constants representing the allowable values for IdpGroupMappingLifecycleStateEnum
const (
	IdpGroupMappingLifecycleStateCreating IdpGroupMappingLifecycleStateEnum = "CREATING"
	IdpGroupMappingLifecycleStateActive   IdpGroupMappingLifecycleStateEnum = "ACTIVE"
	IdpGroupMappingLifecycleStateInactive IdpGroupMappingLifecycleStateEnum = "INACTIVE"
	IdpGroupMappingLifecycleStateDeleting IdpGroupMappingLifecycleStateEnum = "DELETING"
	IdpGroupMappingLifecycleStateDeleted  IdpGroupMappingLifecycleStateEnum = "DELETED"
)

var mappingIdpGroupMappingLifecycleState = map[string]IdpGroupMappingLifecycleStateEnum{
	"CREATING": IdpGroupMappingLifecycleStateCreating,
	"ACTIVE":   IdpGroupMappingLifecycleStateActive,
	"INACTIVE": IdpGroupMappingLifecycleStateInactive,
	"DELETING": IdpGroupMappingLifecycleStateDeleting,
	"DELETED":  IdpGroupMappingLifecycleStateDeleted,
}

// GetIdpGroupMappingLifecycleStateEnumValues Enumerates the set of values for IdpGroupMappingLifecycleStateEnum
func GetIdpGroupMappingLifecycleStateEnumValues() []IdpGroupMappingLifecycleStateEnum {
	values := make([]IdpGroupMappingLifecycleStateEnum, 0)
	for _, v := range mappingIdpGroupMappingLifecycleState {
		values = append(values, v)
	}
	return values
}
