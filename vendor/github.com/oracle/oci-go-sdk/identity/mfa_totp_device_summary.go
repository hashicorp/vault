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

// MfaTotpDeviceSummary As the name suggests, a `MfaTotpDeviceSummary` object contains information about a `MfaTotpDevice`.
type MfaTotpDeviceSummary struct {

	// The OCID of the MFA TOTP Device.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the user the MFA TOTP device belongs to.
	UserId *string `mandatory:"true" json:"userId"`

	// Date and time the `MfaTotpDevice` object was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The MFA TOTP device's current state.
	LifecycleState MfaTotpDeviceSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Flag to indicate if the MFA TOTP device has been activated
	IsActivated *bool `mandatory:"true" json:"isActivated"`

	// Date and time when this MFA TOTP device will expire, in the format defined by RFC3339.
	// Null if it never expires.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeExpires *common.SDKTime `mandatory:"false" json:"timeExpires"`

	// The detailed status of INACTIVE lifecycleState.
	// Allowed values are:
	//  - 1 - SUSPENDED
	//  - 2 - DISABLED
	//  - 4 - BLOCKED
	//  - 8 - LOCKED
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`
}

func (m MfaTotpDeviceSummary) String() string {
	return common.PointerString(m)
}

// MfaTotpDeviceSummaryLifecycleStateEnum Enum with underlying type: string
type MfaTotpDeviceSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for MfaTotpDeviceSummaryLifecycleStateEnum
const (
	MfaTotpDeviceSummaryLifecycleStateCreating MfaTotpDeviceSummaryLifecycleStateEnum = "CREATING"
	MfaTotpDeviceSummaryLifecycleStateActive   MfaTotpDeviceSummaryLifecycleStateEnum = "ACTIVE"
	MfaTotpDeviceSummaryLifecycleStateInactive MfaTotpDeviceSummaryLifecycleStateEnum = "INACTIVE"
	MfaTotpDeviceSummaryLifecycleStateDeleting MfaTotpDeviceSummaryLifecycleStateEnum = "DELETING"
	MfaTotpDeviceSummaryLifecycleStateDeleted  MfaTotpDeviceSummaryLifecycleStateEnum = "DELETED"
)

var mappingMfaTotpDeviceSummaryLifecycleState = map[string]MfaTotpDeviceSummaryLifecycleStateEnum{
	"CREATING": MfaTotpDeviceSummaryLifecycleStateCreating,
	"ACTIVE":   MfaTotpDeviceSummaryLifecycleStateActive,
	"INACTIVE": MfaTotpDeviceSummaryLifecycleStateInactive,
	"DELETING": MfaTotpDeviceSummaryLifecycleStateDeleting,
	"DELETED":  MfaTotpDeviceSummaryLifecycleStateDeleted,
}

// GetMfaTotpDeviceSummaryLifecycleStateEnumValues Enumerates the set of values for MfaTotpDeviceSummaryLifecycleStateEnum
func GetMfaTotpDeviceSummaryLifecycleStateEnumValues() []MfaTotpDeviceSummaryLifecycleStateEnum {
	values := make([]MfaTotpDeviceSummaryLifecycleStateEnum, 0)
	for _, v := range mappingMfaTotpDeviceSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
