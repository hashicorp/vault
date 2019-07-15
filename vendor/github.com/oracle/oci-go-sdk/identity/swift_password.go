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

// SwiftPassword **Deprecated. Use AuthToken instead.**
// Swift is the OpenStack object storage service. A `SwiftPassword` is an Oracle-provided password for using a
// Swift client with the Object Storage Service. This password is associated with
// the user's Console login. Swift passwords never expire. A user can have up to two Swift passwords at a time.
// **Note:** The password is always an Oracle-generated string; you can't change it to a string of your choice.
// For more information, see Managing User Credentials (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingcredentials.htm).
type SwiftPassword struct {

	// The Swift password. The value is available only in the response for `CreateSwiftPassword`, and not
	// for `ListSwiftPasswords` or `UpdateSwiftPassword`.
	Password *string `mandatory:"false" json:"password"`

	// The OCID of the Swift password.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the user the password belongs to.
	UserId *string `mandatory:"false" json:"userId"`

	// The description you assign to the Swift password. Does not have to be unique, and it's changeable.
	Description *string `mandatory:"false" json:"description"`

	// Date and time the `SwiftPassword` object was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// Date and time when this password will expire, in the format defined by RFC3339.
	// Null if it never expires.
	// Example: `2016-08-25T21:10:29.600Z`
	ExpiresOn *common.SDKTime `mandatory:"false" json:"expiresOn"`

	// The password's current state. After creating a password, make sure its `lifecycleState` changes from
	// CREATING to ACTIVE before using it.
	LifecycleState SwiftPasswordLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`
}

func (m SwiftPassword) String() string {
	return common.PointerString(m)
}

// SwiftPasswordLifecycleStateEnum Enum with underlying type: string
type SwiftPasswordLifecycleStateEnum string

// Set of constants representing the allowable values for SwiftPasswordLifecycleStateEnum
const (
	SwiftPasswordLifecycleStateCreating SwiftPasswordLifecycleStateEnum = "CREATING"
	SwiftPasswordLifecycleStateActive   SwiftPasswordLifecycleStateEnum = "ACTIVE"
	SwiftPasswordLifecycleStateInactive SwiftPasswordLifecycleStateEnum = "INACTIVE"
	SwiftPasswordLifecycleStateDeleting SwiftPasswordLifecycleStateEnum = "DELETING"
	SwiftPasswordLifecycleStateDeleted  SwiftPasswordLifecycleStateEnum = "DELETED"
)

var mappingSwiftPasswordLifecycleState = map[string]SwiftPasswordLifecycleStateEnum{
	"CREATING": SwiftPasswordLifecycleStateCreating,
	"ACTIVE":   SwiftPasswordLifecycleStateActive,
	"INACTIVE": SwiftPasswordLifecycleStateInactive,
	"DELETING": SwiftPasswordLifecycleStateDeleting,
	"DELETED":  SwiftPasswordLifecycleStateDeleted,
}

// GetSwiftPasswordLifecycleStateEnumValues Enumerates the set of values for SwiftPasswordLifecycleStateEnum
func GetSwiftPasswordLifecycleStateEnumValues() []SwiftPasswordLifecycleStateEnum {
	values := make([]SwiftPasswordLifecycleStateEnum, 0)
	for _, v := range mappingSwiftPasswordLifecycleState {
		values = append(values, v)
	}
	return values
}
