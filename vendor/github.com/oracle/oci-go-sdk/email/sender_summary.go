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

// SenderSummary The email addresses and `senderId` representing an approved sender.
type SenderSummary struct {

	// The OCID for the compartment.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The email address of the sender.
	EmailAddress *string `mandatory:"false" json:"emailAddress"`

	// The unique ID of the sender.
	Id *string `mandatory:"false" json:"id"`

	// The current status of the approved sender.
	LifecycleState SenderSummaryLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// Date time the approved sender was added, in "YYYY-MM-ddThh:mmZ"
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

func (m SenderSummary) String() string {
	return common.PointerString(m)
}

// SenderSummaryLifecycleStateEnum Enum with underlying type: string
type SenderSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for SenderSummaryLifecycleStateEnum
const (
	SenderSummaryLifecycleStateCreating SenderSummaryLifecycleStateEnum = "CREATING"
	SenderSummaryLifecycleStateActive   SenderSummaryLifecycleStateEnum = "ACTIVE"
	SenderSummaryLifecycleStateDeleting SenderSummaryLifecycleStateEnum = "DELETING"
	SenderSummaryLifecycleStateDeleted  SenderSummaryLifecycleStateEnum = "DELETED"
)

var mappingSenderSummaryLifecycleState = map[string]SenderSummaryLifecycleStateEnum{
	"CREATING": SenderSummaryLifecycleStateCreating,
	"ACTIVE":   SenderSummaryLifecycleStateActive,
	"DELETING": SenderSummaryLifecycleStateDeleting,
	"DELETED":  SenderSummaryLifecycleStateDeleted,
}

// GetSenderSummaryLifecycleStateEnumValues Enumerates the set of values for SenderSummaryLifecycleStateEnum
func GetSenderSummaryLifecycleStateEnumValues() []SenderSummaryLifecycleStateEnum {
	values := make([]SenderSummaryLifecycleStateEnum, 0)
	for _, v := range mappingSenderSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
