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

// SmtpCredentialSummary As the name suggests, an `SmtpCredentialSummary` object contains information about an `SmtpCredential`.
// The SMTP credential is used for SMTP authentication with
// the Email Delivery Service (https://docs.cloud.oracle.com/Content/Email/Concepts/overview.htm).
type SmtpCredentialSummary struct {

	// The SMTP user name.
	Username *string `mandatory:"false" json:"username"`

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
	LifecycleState SmtpCredentialSummaryLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`
}

func (m SmtpCredentialSummary) String() string {
	return common.PointerString(m)
}

// SmtpCredentialSummaryLifecycleStateEnum Enum with underlying type: string
type SmtpCredentialSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for SmtpCredentialSummaryLifecycleStateEnum
const (
	SmtpCredentialSummaryLifecycleStateCreating SmtpCredentialSummaryLifecycleStateEnum = "CREATING"
	SmtpCredentialSummaryLifecycleStateActive   SmtpCredentialSummaryLifecycleStateEnum = "ACTIVE"
	SmtpCredentialSummaryLifecycleStateInactive SmtpCredentialSummaryLifecycleStateEnum = "INACTIVE"
	SmtpCredentialSummaryLifecycleStateDeleting SmtpCredentialSummaryLifecycleStateEnum = "DELETING"
	SmtpCredentialSummaryLifecycleStateDeleted  SmtpCredentialSummaryLifecycleStateEnum = "DELETED"
)

var mappingSmtpCredentialSummaryLifecycleState = map[string]SmtpCredentialSummaryLifecycleStateEnum{
	"CREATING": SmtpCredentialSummaryLifecycleStateCreating,
	"ACTIVE":   SmtpCredentialSummaryLifecycleStateActive,
	"INACTIVE": SmtpCredentialSummaryLifecycleStateInactive,
	"DELETING": SmtpCredentialSummaryLifecycleStateDeleting,
	"DELETED":  SmtpCredentialSummaryLifecycleStateDeleted,
}

// GetSmtpCredentialSummaryLifecycleStateEnumValues Enumerates the set of values for SmtpCredentialSummaryLifecycleStateEnum
func GetSmtpCredentialSummaryLifecycleStateEnumValues() []SmtpCredentialSummaryLifecycleStateEnum {
	values := make([]SmtpCredentialSummaryLifecycleStateEnum, 0)
	for _, v := range mappingSmtpCredentialSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
