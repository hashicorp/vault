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

// CustomerSecretKeySummary As the name suggests, a `CustomerSecretKeySummary` object contains information about a `CustomerSecretKey`.
// A `CustomerSecretKey` is an Oracle-provided key for using the Object Storage Service's Amazon S3 compatible API.
type CustomerSecretKeySummary struct {

	// The OCID of the secret key.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the user the password belongs to.
	UserId *string `mandatory:"false" json:"userId"`

	// The displayName you assign to the secret key. Does not have to be unique, and it's changeable.
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
	LifecycleState CustomerSecretKeySummaryLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`
}

func (m CustomerSecretKeySummary) String() string {
	return common.PointerString(m)
}

// CustomerSecretKeySummaryLifecycleStateEnum Enum with underlying type: string
type CustomerSecretKeySummaryLifecycleStateEnum string

// Set of constants representing the allowable values for CustomerSecretKeySummaryLifecycleStateEnum
const (
	CustomerSecretKeySummaryLifecycleStateCreating CustomerSecretKeySummaryLifecycleStateEnum = "CREATING"
	CustomerSecretKeySummaryLifecycleStateActive   CustomerSecretKeySummaryLifecycleStateEnum = "ACTIVE"
	CustomerSecretKeySummaryLifecycleStateInactive CustomerSecretKeySummaryLifecycleStateEnum = "INACTIVE"
	CustomerSecretKeySummaryLifecycleStateDeleting CustomerSecretKeySummaryLifecycleStateEnum = "DELETING"
	CustomerSecretKeySummaryLifecycleStateDeleted  CustomerSecretKeySummaryLifecycleStateEnum = "DELETED"
)

var mappingCustomerSecretKeySummaryLifecycleState = map[string]CustomerSecretKeySummaryLifecycleStateEnum{
	"CREATING": CustomerSecretKeySummaryLifecycleStateCreating,
	"ACTIVE":   CustomerSecretKeySummaryLifecycleStateActive,
	"INACTIVE": CustomerSecretKeySummaryLifecycleStateInactive,
	"DELETING": CustomerSecretKeySummaryLifecycleStateDeleting,
	"DELETED":  CustomerSecretKeySummaryLifecycleStateDeleted,
}

// GetCustomerSecretKeySummaryLifecycleStateEnumValues Enumerates the set of values for CustomerSecretKeySummaryLifecycleStateEnum
func GetCustomerSecretKeySummaryLifecycleStateEnumValues() []CustomerSecretKeySummaryLifecycleStateEnum {
	values := make([]CustomerSecretKeySummaryLifecycleStateEnum, 0)
	for _, v := range mappingCustomerSecretKeySummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
