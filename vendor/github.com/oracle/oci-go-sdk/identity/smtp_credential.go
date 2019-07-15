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

// SmtpCredential Simple Mail Transfer Protocol (SMTP) credentials are needed to send email through Email Delivery.
// The SMTP credentials are used for SMTP authentication with the service. The credentials never expire.
// A user can have up to 2 SMTP credentials at a time.
// **Note:** The credential set is always an Oracle-generated SMTP user name and password pair;
// you cannot designate the SMTP user name or the SMTP password.
// For more information, see Managing User Credentials (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingcredentials.htm#SMTP).
type SmtpCredential struct {

	// The SMTP user name.
	Username *string `mandatory:"false" json:"username"`

	// The SMTP password.
	Password *string `mandatory:"false" json:"password"`

	// The OCID of the SMTP credential.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the user the SMTP credential belongs to.
	UserId *string `mandatory:"false" json:"userId"`

	// The description you assign to the SMTP credential. Does not have to be unique, and it's changeable.
	Description *string `mandatory:"false" json:"description"`

	// Date and time the `SmtpCredential` object was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// Date and time when this credential will expire, in the format defined by RFC3339.
	// Null if it never expires.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeExpires *common.SDKTime `mandatory:"false" json:"timeExpires"`

	// The credential's current state. After creating a SMTP credential, make sure its `lifecycleState` changes from
	// CREATING to ACTIVE before using it.
	LifecycleState SmtpCredentialLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`
}

func (m SmtpCredential) String() string {
	return common.PointerString(m)
}

// SmtpCredentialLifecycleStateEnum Enum with underlying type: string
type SmtpCredentialLifecycleStateEnum string

// Set of constants representing the allowable values for SmtpCredentialLifecycleStateEnum
const (
	SmtpCredentialLifecycleStateCreating SmtpCredentialLifecycleStateEnum = "CREATING"
	SmtpCredentialLifecycleStateActive   SmtpCredentialLifecycleStateEnum = "ACTIVE"
	SmtpCredentialLifecycleStateInactive SmtpCredentialLifecycleStateEnum = "INACTIVE"
	SmtpCredentialLifecycleStateDeleting SmtpCredentialLifecycleStateEnum = "DELETING"
	SmtpCredentialLifecycleStateDeleted  SmtpCredentialLifecycleStateEnum = "DELETED"
)

var mappingSmtpCredentialLifecycleState = map[string]SmtpCredentialLifecycleStateEnum{
	"CREATING": SmtpCredentialLifecycleStateCreating,
	"ACTIVE":   SmtpCredentialLifecycleStateActive,
	"INACTIVE": SmtpCredentialLifecycleStateInactive,
	"DELETING": SmtpCredentialLifecycleStateDeleting,
	"DELETED":  SmtpCredentialLifecycleStateDeleted,
}

// GetSmtpCredentialLifecycleStateEnumValues Enumerates the set of values for SmtpCredentialLifecycleStateEnum
func GetSmtpCredentialLifecycleStateEnumValues() []SmtpCredentialLifecycleStateEnum {
	values := make([]SmtpCredentialLifecycleStateEnum, 0)
	for _, v := range mappingSmtpCredentialLifecycleState {
		values = append(values, v)
	}
	return values
}
