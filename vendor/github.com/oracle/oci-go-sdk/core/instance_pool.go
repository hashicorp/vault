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

// InstancePool An instance pool is a group of instances within the same region that are created based off of the same
// instance configuration. For more information about instance pools and instance configurations, see
// Managing Compute Instances (https://docs.cloud.oracle.com/Content/Compute/Concepts/instancemanagement.htm).
type InstancePool struct {

	// The OCID of the instance pool.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment containing the instance pool.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the instance configuration associated with the instance pool.
	InstanceConfigurationId *string `mandatory:"true" json:"instanceConfigurationId"`

	// The current state of the instance pool.
	LifecycleState InstancePoolLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The placement configurations for the instance pool.
	PlacementConfigurations []InstancePoolPlacementConfiguration `mandatory:"true" json:"placementConfigurations"`

	// The number of instances that should be in the instance pool.
	Size *int `mandatory:"true" json:"size"`

	// The date and time the instance pool was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The user-friendly name.  Does not have to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The load balancers attached to the instance pool.
	LoadBalancers []InstancePoolLoadBalancerAttachment `mandatory:"false" json:"loadBalancers"`
}

func (m InstancePool) String() string {
	return common.PointerString(m)
}

// InstancePoolLifecycleStateEnum Enum with underlying type: string
type InstancePoolLifecycleStateEnum string

// Set of constants representing the allowable values for InstancePoolLifecycleStateEnum
const (
	InstancePoolLifecycleStateProvisioning InstancePoolLifecycleStateEnum = "PROVISIONING"
	InstancePoolLifecycleStateScaling      InstancePoolLifecycleStateEnum = "SCALING"
	InstancePoolLifecycleStateStarting     InstancePoolLifecycleStateEnum = "STARTING"
	InstancePoolLifecycleStateStopping     InstancePoolLifecycleStateEnum = "STOPPING"
	InstancePoolLifecycleStateTerminating  InstancePoolLifecycleStateEnum = "TERMINATING"
	InstancePoolLifecycleStateStopped      InstancePoolLifecycleStateEnum = "STOPPED"
	InstancePoolLifecycleStateTerminated   InstancePoolLifecycleStateEnum = "TERMINATED"
	InstancePoolLifecycleStateRunning      InstancePoolLifecycleStateEnum = "RUNNING"
)

var mappingInstancePoolLifecycleState = map[string]InstancePoolLifecycleStateEnum{
	"PROVISIONING": InstancePoolLifecycleStateProvisioning,
	"SCALING":      InstancePoolLifecycleStateScaling,
	"STARTING":     InstancePoolLifecycleStateStarting,
	"STOPPING":     InstancePoolLifecycleStateStopping,
	"TERMINATING":  InstancePoolLifecycleStateTerminating,
	"STOPPED":      InstancePoolLifecycleStateStopped,
	"TERMINATED":   InstancePoolLifecycleStateTerminated,
	"RUNNING":      InstancePoolLifecycleStateRunning,
}

// GetInstancePoolLifecycleStateEnumValues Enumerates the set of values for InstancePoolLifecycleStateEnum
func GetInstancePoolLifecycleStateEnumValues() []InstancePoolLifecycleStateEnum {
	values := make([]InstancePoolLifecycleStateEnum, 0)
	for _, v := range mappingInstancePoolLifecycleState {
		values = append(values, v)
	}
	return values
}
