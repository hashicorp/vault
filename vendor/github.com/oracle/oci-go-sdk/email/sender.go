// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Email Delivery API
//
// API for the Email Delivery service. Use this API to send high-volume, application-generated
// emails. For more information, see Overview of the Email Delivery Service (https://docs.cloud.oracle.com/iaas/Content/Email/Concepts/overview.htm).
//
// **Note:** Write actions (POST, UPDATE, DELETE) may take several minutes to propagate and be reflected by the API. If a subsequent read request fails to reflect your changes, wait a few minutes and try again.
//

package email

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Sender The full information representing an approved sender.
type Sender struct {

	// The OCID for the compartment.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// Email address of the sender.
	EmailAddress *string `mandatory:"false" json:"emailAddress"`

	// The unique OCID of the sender.
	Id *string `mandatory:"false" json:"id"`

	// Value of the SPF field. For more information about SPF, please see
	// SPF Authentication (https://docs.cloud.oracle.com/Content/Email/Concepts/overview.htm#components).
	IsSpf *bool `mandatory:"false" json:"isSpf"`

	// The sender's current lifecycle state.
	LifecycleState SenderLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The date and time the approved sender was added in "YYYY-MM-ddThh:mmZ"
	// format with a Z offset, as defined by RFC 3339.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m Sender) String() string {
	return common.PointerString(m)
}

// SenderLifecycleStateEnum Enum with underlying type: string
type SenderLifecycleStateEnum string

// Set of constants representing the allowable values for SenderLifecycleStateEnum
const (
	SenderLifecycleStateCreating SenderLifecycleStateEnum = "CREATING"
	SenderLifecycleStateActive   SenderLifecycleStateEnum = "ACTIVE"
	SenderLifecycleStateDeleting SenderLifecycleStateEnum = "DELETING"
	SenderLifecycleStateDeleted  SenderLifecycleStateEnum = "DELETED"
)

var mappingSenderLifecycleState = map[string]SenderLifecycleStateEnum{
	"CREATING": SenderLifecycleStateCreating,
	"ACTIVE":   SenderLifecycleStateActive,
	"DELETING": SenderLifecycleStateDeleting,
	"DELETED":  SenderLifecycleStateDeleted,
}

// GetSenderLifecycleStateEnumValues Enumerates the set of values for SenderLifecycleStateEnum
func GetSenderLifecycleStateEnumValues() []SenderLifecycleStateEnum {
	values := make([]SenderLifecycleStateEnum, 0)
	for _, v := range mappingSenderLifecycleState {
		values = append(values, v)
	}
	return values
}
