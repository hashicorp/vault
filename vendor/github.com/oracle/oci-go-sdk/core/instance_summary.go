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

// InstanceSummary Condensed instance data when listing instances in an instance pool.
type InstanceSummary struct {

	// The OCID of the instance.
	Id *string `mandatory:"true" json:"id"`

	// The availability domain the instance is running in.
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID of the compartment that contains the instance.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the instance confgiuration used to create the instance.
	InstanceConfigurationId *string `mandatory:"true" json:"instanceConfigurationId"`

	// The region that contains the availability domain the instance is running in.
	Region *string `mandatory:"true" json:"region"`

	// The current state of the instance pool instance.
	State *string `mandatory:"true" json:"state"`

	// The date and time the instance pool instance was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The user-friendly name.  Does not have to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The name of the Fault Domain the instance is running in.
	FaultDomain *string `mandatory:"false" json:"faultDomain"`

	// The shape of an instance. The shape determines the number of CPUs, amount of memory,
	// and other resources allocated to the instance.
	// You can enumerate all available shapes by calling ListShapes.
	Shape *string `mandatory:"false" json:"shape"`

	// The load balancer backends that are configured for the instance pool instance.
	LoadBalancerBackends []InstancePoolInstanceLoadBalancerBackend `mandatory:"false" json:"loadBalancerBackends"`
}

func (m InstanceSummary) String() string {
	return common.PointerString(m)
}
