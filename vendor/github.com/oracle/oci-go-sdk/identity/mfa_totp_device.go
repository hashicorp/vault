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

// MfaTotpDevice Users can enable multi-factor authentication (MFA) for their own user accounts. After MFA is enabled, the
// user is prompted for a time-based one-time password (TOTP) to authenticate before they can sign in to the
// Console. To enable multi-factor authentication, the user must register a mobile device with a TOTP authenticator app
// installed. The registration process creates the `MfaTotpDevice` object. The registration process requires
// interaction with the Console and cannot be completed programmatically. For more information, see
// Managing Multi-Factor Authentication (https://docs.cloud.oracle.com/Content/Identity/Tasks/usingmfa.htm).
type MfaTotpDevice struct {

	// The OCID of the MFA TOTP device.
	Id *string `mandatory:"true" json:"id"`

	// The seed for the MFA TOTP device (Base32 encoded).
	Seed *string `mandatory:"true" json:"seed"`

	// The OCID of the user the MFA TOTP device belongs to.
	UserId *string `mandatory:"true" json:"userId"`

	// Date and time the `MfaTotpDevice` object was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The MFA TOTP device's current state. After creating the MFA TOTP device, make sure its `lifecycleState` changes from
	// CREATING to ACTIVE before using it.
	LifecycleState MfaTotpDeviceLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Flag to indicate if the MFA TOTP device has been activated.
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

func (m MfaTotpDevice) String() string {
	return common.PointerString(m)
}

// MfaTotpDeviceLifecycleStateEnum Enum with underlying type: string
type MfaTotpDeviceLifecycleStateEnum string

// Set of constants representing the allowable values for MfaTotpDeviceLifecycleStateEnum
const (
	MfaTotpDeviceLifecycleStateCreating MfaTotpDeviceLifecycleStateEnum = "CREATING"
	MfaTotpDeviceLifecycleStateActive   MfaTotpDeviceLifecycleStateEnum = "ACTIVE"
	MfaTotpDeviceLifecycleStateInactive MfaTotpDeviceLifecycleStateEnum = "INACTIVE"
	MfaTotpDeviceLifecycleStateDeleting MfaTotpDeviceLifecycleStateEnum = "DELETING"
	MfaTotpDeviceLifecycleStateDeleted  MfaTotpDeviceLifecycleStateEnum = "DELETED"
)

var mappingMfaTotpDeviceLifecycleState = map[string]MfaTotpDeviceLifecycleStateEnum{
	"CREATING": MfaTotpDeviceLifecycleStateCreating,
	"ACTIVE":   MfaTotpDeviceLifecycleStateActive,
	"INACTIVE": MfaTotpDeviceLifecycleStateInactive,
	"DELETING": MfaTotpDeviceLifecycleStateDeleting,
	"DELETED":  MfaTotpDeviceLifecycleStateDeleted,
}

// GetMfaTotpDeviceLifecycleStateEnumValues Enumerates the set of values for MfaTotpDeviceLifecycleStateEnum
func GetMfaTotpDeviceLifecycleStateEnumValues() []MfaTotpDeviceLifecycleStateEnum {
	values := make([]MfaTotpDeviceLifecycleStateEnum, 0)
	for _, v := range mappingMfaTotpDeviceLifecycleState {
		values = append(values, v)
	}
	return values
}
