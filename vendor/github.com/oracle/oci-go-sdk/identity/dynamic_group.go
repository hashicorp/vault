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

// DynamicGroup A dynamic group defines a matching rule. Every bare metal or virtual machine instance is deployed with an instance certificate.
// The certificate contains metadata about the instance. This includes the instance OCID and the compartment OCID, along
// with a few other optional properties. When an API call is made using this instance certificate as the authenticator,
// the certificate can be matched to one or multiple dynamic groups. The instance can then get access to the API
// based on the permissions granted in policies written for the dynamic groups.
// This works like regular user/group membership. But in that case, the membership is a static relationship, whereas
// in a dynamic group, the membership of an instance certificate to a dynamic group is determined during runtime.
// For more information, see Managing Dynamic Groups (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingdynamicgroups.htm).
type DynamicGroup struct {

	// The OCID of the group.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the tenancy containing the group.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The name you assign to the group during creation. The name must be unique across all groups in
	// the tenancy and cannot be changed.
	Name *string `mandatory:"true" json:"name"`

	// The description you assign to the group. Does not have to be unique, and it's changeable.
	Description *string `mandatory:"true" json:"description"`

	// A rule string that defines which instance certificates will be matched.
	// For syntax, see Managing Dynamic Groups (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingdynamicgroups.htm).
	MatchingRule *string `mandatory:"true" json:"matchingRule"`

	// Date and time the group was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The group's current state. After creating a group, make sure its `lifecycleState` changes from CREATING to
	// ACTIVE before using it.
	LifecycleState DynamicGroupLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m DynamicGroup) String() string {
	return common.PointerString(m)
}

// DynamicGroupLifecycleStateEnum Enum with underlying type: string
type DynamicGroupLifecycleStateEnum string

// Set of constants representing the allowable values for DynamicGroupLifecycleStateEnum
const (
	DynamicGroupLifecycleStateCreating DynamicGroupLifecycleStateEnum = "CREATING"
	DynamicGroupLifecycleStateActive   DynamicGroupLifecycleStateEnum = "ACTIVE"
	DynamicGroupLifecycleStateInactive DynamicGroupLifecycleStateEnum = "INACTIVE"
	DynamicGroupLifecycleStateDeleting DynamicGroupLifecycleStateEnum = "DELETING"
	DynamicGroupLifecycleStateDeleted  DynamicGroupLifecycleStateEnum = "DELETED"
)

var mappingDynamicGroupLifecycleState = map[string]DynamicGroupLifecycleStateEnum{
	"CREATING": DynamicGroupLifecycleStateCreating,
	"ACTIVE":   DynamicGroupLifecycleStateActive,
	"INACTIVE": DynamicGroupLifecycleStateInactive,
	"DELETING": DynamicGroupLifecycleStateDeleting,
	"DELETED":  DynamicGroupLifecycleStateDeleted,
}

// GetDynamicGroupLifecycleStateEnumValues Enumerates the set of values for DynamicGroupLifecycleStateEnum
func GetDynamicGroupLifecycleStateEnumValues() []DynamicGroupLifecycleStateEnum {
	values := make([]DynamicGroupLifecycleStateEnum, 0)
	for _, v := range mappingDynamicGroupLifecycleState {
		values = append(values, v)
	}
	return values
}
