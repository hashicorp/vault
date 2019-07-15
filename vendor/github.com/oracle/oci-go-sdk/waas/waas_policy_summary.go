// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Web Application Acceleration and Security Services API
//
// OCI Web Application Acceleration and Security Services
//

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
)

// WaasPolicySummary Summary information about a WAAS policy.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type WaasPolicySummary struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the WAAS policy.
	Id *string `mandatory:"false" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the WAAS policy's compartment.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The user-friendly name of the WAAS policy. The name can be changed and does not need to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The web application domain that the WAAS policy protects.
	Domain *string `mandatory:"false" json:"domain"`

	// The current lifecycle state of the WAAS policy.
	LifecycleState WaasPolicySummaryLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The date and time the policy was created, expressed in RFC 3339 timestamp format.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// A simple key-value pair without any defined schema.
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// A key-value pair with a defined schema that restricts the values of tags. These predefined keys are scoped to namespaces.
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m WaasPolicySummary) String() string {
	return common.PointerString(m)
}

// WaasPolicySummaryLifecycleStateEnum Enum with underlying type: string
type WaasPolicySummaryLifecycleStateEnum string

// Set of constants representing the allowable values for WaasPolicySummaryLifecycleStateEnum
const (
	WaasPolicySummaryLifecycleStateCreating WaasPolicySummaryLifecycleStateEnum = "CREATING"
	WaasPolicySummaryLifecycleStateActive   WaasPolicySummaryLifecycleStateEnum = "ACTIVE"
	WaasPolicySummaryLifecycleStateFailed   WaasPolicySummaryLifecycleStateEnum = "FAILED"
	WaasPolicySummaryLifecycleStateUpdating WaasPolicySummaryLifecycleStateEnum = "UPDATING"
	WaasPolicySummaryLifecycleStateDeleting WaasPolicySummaryLifecycleStateEnum = "DELETING"
	WaasPolicySummaryLifecycleStateDeleted  WaasPolicySummaryLifecycleStateEnum = "DELETED"
)

var mappingWaasPolicySummaryLifecycleState = map[string]WaasPolicySummaryLifecycleStateEnum{
	"CREATING": WaasPolicySummaryLifecycleStateCreating,
	"ACTIVE":   WaasPolicySummaryLifecycleStateActive,
	"FAILED":   WaasPolicySummaryLifecycleStateFailed,
	"UPDATING": WaasPolicySummaryLifecycleStateUpdating,
	"DELETING": WaasPolicySummaryLifecycleStateDeleting,
	"DELETED":  WaasPolicySummaryLifecycleStateDeleted,
}

// GetWaasPolicySummaryLifecycleStateEnumValues Enumerates the set of values for WaasPolicySummaryLifecycleStateEnum
func GetWaasPolicySummaryLifecycleStateEnumValues() []WaasPolicySummaryLifecycleStateEnum {
	values := make([]WaasPolicySummaryLifecycleStateEnum, 0)
	for _, v := range mappingWaasPolicySummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
