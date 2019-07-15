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

// CustomerSecretKey A `CustomerSecretKey` is an Oracle-provided key for using the Object Storage Service's
// Amazon S3 compatible API (https://docs.cloud.oracle.com/Content/Object/Tasks/s3compatibleapi.htm).
// A user can have up to two secret keys at a time.
// **Note:** The secret key is always an Oracle-generated string; you can't change it to a string of your choice.
// For more information, see Managing User Credentials (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingcredentials.htm).
type CustomerSecretKey struct {

	// The secret key.
	Key *string `mandatory:"false" json:"key"`

	// The OCID of the secret key.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the user the password belongs to.
	UserId *string `mandatory:"false" json:"userId"`

	// The display name you assign to the secret key. Does not have to be unique, and it's changeable.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Date and time the `CustomerSecretKey` object was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// Date and time when this password will expire, in the format defined by RFC3339.
	// Null if it never expires.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeExpires *common.SDKTime `mandatory:"false" json:"timeExpires"`

	// The secret key's current state. After creating a secret key, make sure its `lifecycleState` changes from
	// CREATING to ACTIVE before using it.
	LifecycleState CustomerSecretKeyLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`
}

func (m CustomerSecretKey) String() string {
	return common.PointerString(m)
}

// CustomerSecretKeyLifecycleStateEnum Enum with underlying type: string
type CustomerSecretKeyLifecycleStateEnum string

// Set of constants representing the allowable values for CustomerSecretKeyLifecycleStateEnum
const (
	CustomerSecretKeyLifecycleStateCreating CustomerSecretKeyLifecycleStateEnum = "CREATING"
	CustomerSecretKeyLifecycleStateActive   CustomerSecretKeyLifecycleStateEnum = "ACTIVE"
	CustomerSecretKeyLifecycleStateInactive CustomerSecretKeyLifecycleStateEnum = "INACTIVE"
	CustomerSecretKeyLifecycleStateDeleting CustomerSecretKeyLifecycleStateEnum = "DELETING"
	CustomerSecretKeyLifecycleStateDeleted  CustomerSecretKeyLifecycleStateEnum = "DELETED"
)

var mappingCustomerSecretKeyLifecycleState = map[string]CustomerSecretKeyLifecycleStateEnum{
	"CREATING": CustomerSecretKeyLifecycleStateCreating,
	"ACTIVE":   CustomerSecretKeyLifecycleStateActive,
	"INACTIVE": CustomerSecretKeyLifecycleStateInactive,
	"DELETING": CustomerSecretKeyLifecycleStateDeleting,
	"DELETED":  CustomerSecretKeyLifecycleStateDeleted,
}

// GetCustomerSecretKeyLifecycleStateEnumValues Enumerates the set of values for CustomerSecretKeyLifecycleStateEnum
func GetCustomerSecretKeyLifecycleStateEnumValues() []CustomerSecretKeyLifecycleStateEnum {
	values := make([]CustomerSecretKeyLifecycleStateEnum, 0)
	for _, v := range mappingCustomerSecretKeyLifecycleState {
		values = append(values, v)
	}
	return values
}
