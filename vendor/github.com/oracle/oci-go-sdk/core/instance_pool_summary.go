// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"github.com/oracle/oci-go-sdk/common"
)

// InstancePoolSummary Summary information for an instance pool.
type InstancePoolSummary struct {

	// The OCID of the instance pool.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment containing the instance pool.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the instance configuration associated with the instance pool.
	InstanceConfigurationId *string `mandatory:"true" json:"instanceConfigurationId"`

	// The current state of the instance pool.
	LifecycleState InstancePoolSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The availability domains for the instance pool.
	AvailabilityDomains []string `mandatory:"true" json:"availabilityDomains"`

	// The number of instances that should be in the instance pool.
	Size *int `mandatory:"true" json:"size"`

	// The date and time the instance pool was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The user-friendly name.  Does not have to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m InstancePoolSummary) String() string {
	return common.PointerString(m)
}

// InstancePoolSummaryLifecycleStateEnum Enum with underlying type: string
type InstancePoolSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for InstancePoolSummaryLifecycleStateEnum
const (
	InstancePoolSummaryLifecycleStateProvisioning InstancePoolSummaryLifecycleStateEnum = "PROVISIONING"
	InstancePoolSummaryLifecycleStateScaling      InstancePoolSummaryLifecycleStateEnum = "SCALING"
	InstancePoolSummaryLifecycleStateStarting     InstancePoolSummaryLifecycleStateEnum = "STARTING"
	InstancePoolSummaryLifecycleStateStopping     InstancePoolSummaryLifecycleStateEnum = "STOPPING"
	InstancePoolSummaryLifecycleStateTerminating  InstancePoolSummaryLifecycleStateEnum = "TERMINATING"
	InstancePoolSummaryLifecycleStateStopped      InstancePoolSummaryLifecycleStateEnum = "STOPPED"
	InstancePoolSummaryLifecycleStateTerminated   InstancePoolSummaryLifecycleStateEnum = "TERMINATED"
	InstancePoolSummaryLifecycleStateRunning      InstancePoolSummaryLifecycleStateEnum = "RUNNING"
)

var mappingInstancePoolSummaryLifecycleState = map[string]InstancePoolSummaryLifecycleStateEnum{
	"PROVISIONING": InstancePoolSummaryLifecycleStateProvisioning,
	"SCALING":      InstancePoolSummaryLifecycleStateScaling,
	"STARTING":     InstancePoolSummaryLifecycleStateStarting,
	"STOPPING":     InstancePoolSummaryLifecycleStateStopping,
	"TERMINATING":  InstancePoolSummaryLifecycleStateTerminating,
	"STOPPED":      InstancePoolSummaryLifecycleStateStopped,
	"TERMINATED":   InstancePoolSummaryLifecycleStateTerminated,
	"RUNNING":      InstancePoolSummaryLifecycleStateRunning,
}

// GetInstancePoolSummaryLifecycleStateEnumValues Enumerates the set of values for InstancePoolSummaryLifecycleStateEnum
func GetInstancePoolSummaryLifecycleStateEnumValues() []InstancePoolSummaryLifecycleStateEnum {
	values := make([]InstancePoolSummaryLifecycleStateEnum, 0)
	for _, v := range mappingInstancePoolSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
