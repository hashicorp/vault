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

// AuthToken An `AuthToken` is an Oracle-generated token string that you can use to authenticate with third-party APIs
// that do not support Oracle Cloud Infrastructure's signature-based authentication. For example, use an `AuthToken`
// to authenticate with a Swift client with the Object Storage Service.
// The auth token is associated with the user's Console login. Auth tokens never expire. A user can have up to two
// auth tokens at a time.
// **Note:** The token is always an Oracle-generated string; you can't change it to a string of your choice.
// For more information, see Managing User Credentials (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingcredentials.htm).
type AuthToken struct {

	// The auth token. The value is available only in the response for `CreateAuthToken`, and not
	// for `ListAuthTokens` or `UpdateAuthToken`.
	Token *string `mandatory:"false" json:"token"`

	// The OCID of the auth token.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the user the auth token belongs to.
	UserId *string `mandatory:"false" json:"userId"`

	// The description you assign to the auth token. Does not have to be unique, and it's changeable.
	Description *string `mandatory:"false" json:"description"`

	// Date and time the `AuthToken` object was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// Date and time when this auth token will expire, in the format defined by RFC3339.
	// Null if it never expires.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeExpires *common.SDKTime `mandatory:"false" json:"timeExpires"`

	// The token's current state. After creating an auth token, make sure its `lifecycleState` changes from
	// CREATING to ACTIVE before using it.
	LifecycleState AuthTokenLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`
}

func (m AuthToken) String() string {
	return common.PointerString(m)
}

// AuthTokenLifecycleStateEnum Enum with underlying type: string
type AuthTokenLifecycleStateEnum string

// Set of constants representing the allowable values for AuthTokenLifecycleStateEnum
const (
	AuthTokenLifecycleStateCreating AuthTokenLifecycleStateEnum = "CREATING"
	AuthTokenLifecycleStateActive   AuthTokenLifecycleStateEnum = "ACTIVE"
	AuthTokenLifecycleStateInactive AuthTokenLifecycleStateEnum = "INACTIVE"
	AuthTokenLifecycleStateDeleting AuthTokenLifecycleStateEnum = "DELETING"
	AuthTokenLifecycleStateDeleted  AuthTokenLifecycleStateEnum = "DELETED"
)

var mappingAuthTokenLifecycleState = map[string]AuthTokenLifecycleStateEnum{
	"CREATING": AuthTokenLifecycleStateCreating,
	"ACTIVE":   AuthTokenLifecycleStateActive,
	"INACTIVE": AuthTokenLifecycleStateInactive,
	"DELETING": AuthTokenLifecycleStateDeleting,
	"DELETED":  AuthTokenLifecycleStateDeleted,
}

// GetAuthTokenLifecycleStateEnumValues Enumerates the set of values for AuthTokenLifecycleStateEnum
func GetAuthTokenLifecycleStateEnumValues() []AuthTokenLifecycleStateEnum {
	values := make([]AuthTokenLifecycleStateEnum, 0)
	for _, v := range mappingAuthTokenLifecycleState {
		values = append(values, v)
	}
	return values
}
