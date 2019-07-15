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

// UserGroupMembership An object that represents the membership of a user in a group. When you add a user to a group, the result is a
// `UserGroupMembership` with its own OCID. To remove a user from a group, you delete the `UserGroupMembership` object.
type UserGroupMembership struct {

	// The OCID of the membership.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the tenancy containing the user, group, and membership object.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the group.
	GroupId *string `mandatory:"true" json:"groupId"`

	// The OCID of the user.
	UserId *string `mandatory:"true" json:"userId"`

	// Date and time the membership was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The membership's current state.  After creating a membership object, make sure its `lifecycleState` changes
	// from CREATING to ACTIVE before using it.
	LifecycleState UserGroupMembershipLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`
}

func (m UserGroupMembership) String() string {
	return common.PointerString(m)
}

// UserGroupMembershipLifecycleStateEnum Enum with underlying type: string
type UserGroupMembershipLifecycleStateEnum string

// Set of constants representing the allowable values for UserGroupMembershipLifecycleStateEnum
const (
	UserGroupMembershipLifecycleStateCreating UserGroupMembershipLifecycleStateEnum = "CREATING"
	UserGroupMembershipLifecycleStateActive   UserGroupMembershipLifecycleStateEnum = "ACTIVE"
	UserGroupMembershipLifecycleStateInactive UserGroupMembershipLifecycleStateEnum = "INACTIVE"
	UserGroupMembershipLifecycleStateDeleting UserGroupMembershipLifecycleStateEnum = "DELETING"
	UserGroupMembershipLifecycleStateDeleted  UserGroupMembershipLifecycleStateEnum = "DELETED"
)

var mappingUserGroupMembershipLifecycleState = map[string]UserGroupMembershipLifecycleStateEnum{
	"CREATING": UserGroupMembershipLifecycleStateCreating,
	"ACTIVE":   UserGroupMembershipLifecycleStateActive,
	"INACTIVE": UserGroupMembershipLifecycleStateInactive,
	"DELETING": UserGroupMembershipLifecycleStateDeleting,
	"DELETED":  UserGroupMembershipLifecycleStateDeleted,
}

// GetUserGroupMembershipLifecycleStateEnumValues Enumerates the set of values for UserGroupMembershipLifecycleStateEnum
func GetUserGroupMembershipLifecycleStateEnumValues() []UserGroupMembershipLifecycleStateEnum {
	values := make([]UserGroupMembershipLifecycleStateEnum, 0)
	for _, v := range mappingUserGroupMembershipLifecycleState {
		values = append(values, v)
	}
	return values
}
