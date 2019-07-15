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

// UiPasswordInformation Information about the UIPassword, which is a text password that enables a user to sign in to the Console,
// the user interface for interacting with Oracle Cloud Infrastructure.
// For more information about user credentials, see User Credentials (https://docs.cloud.oracle.com/Content/Identity/Concepts/usercredentials.htm).
type UiPasswordInformation struct {

	// The OCID of the user.
	UserId *string `mandatory:"false" json:"userId"`

	// Date and time the password was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The password's current state. After creating a password, make sure its `lifecycleState` changes from
	// CREATING to ACTIVE before using it.
	LifecycleState UiPasswordInformationLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`
}

func (m UiPasswordInformation) String() string {
	return common.PointerString(m)
}

// UiPasswordInformationLifecycleStateEnum Enum with underlying type: string
type UiPasswordInformationLifecycleStateEnum string

// Set of constants representing the allowable values for UiPasswordInformationLifecycleStateEnum
const (
	UiPasswordInformationLifecycleStateCreating UiPasswordInformationLifecycleStateEnum = "CREATING"
	UiPasswordInformationLifecycleStateActive   UiPasswordInformationLifecycleStateEnum = "ACTIVE"
	UiPasswordInformationLifecycleStateInactive UiPasswordInformationLifecycleStateEnum = "INACTIVE"
	UiPasswordInformationLifecycleStateDeleting UiPasswordInformationLifecycleStateEnum = "DELETING"
	UiPasswordInformationLifecycleStateDeleted  UiPasswordInformationLifecycleStateEnum = "DELETED"
)

var mappingUiPasswordInformationLifecycleState = map[string]UiPasswordInformationLifecycleStateEnum{
	"CREATING": UiPasswordInformationLifecycleStateCreating,
	"ACTIVE":   UiPasswordInformationLifecycleStateActive,
	"INACTIVE": UiPasswordInformationLifecycleStateInactive,
	"DELETING": UiPasswordInformationLifecycleStateDeleting,
	"DELETED":  UiPasswordInformationLifecycleStateDeleted,
}

// GetUiPasswordInformationLifecycleStateEnumValues Enumerates the set of values for UiPasswordInformationLifecycleStateEnum
func GetUiPasswordInformationLifecycleStateEnumValues() []UiPasswordInformationLifecycleStateEnum {
	values := make([]UiPasswordInformationLifecycleStateEnum, 0)
	for _, v := range mappingUiPasswordInformationLifecycleState {
		values = append(values, v)
	}
	return values
}
