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

// InstancePoolInstanceLoadBalancerBackend Represents the load balancer Backend that is configured for an instance pool instance.
type InstancePoolInstanceLoadBalancerBackend struct {

	// The OCID of the load balancer attached to the instance pool.
	LoadBalancerId *string `mandatory:"true" json:"loadBalancerId"`

	// The name of the backend set on the load balancer.
	BackendSetName *string `mandatory:"true" json:"backendSetName"`

	// The name of the backend in the backend set.
	BackendName *string `mandatory:"true" json:"backendName"`

	// The health of the backend as observed by the load balancer.
	BackendHealthStatus InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum `mandatory:"true" json:"backendHealthStatus"`
}

func (m InstancePoolInstanceLoadBalancerBackend) String() string {
	return common.PointerString(m)
}

// InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum Enum with underlying type: string
type InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum string

// Set of constants representing the allowable values for InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum
const (
	InstancePoolInstanceLoadBalancerBackendBackendHealthStatusOk       InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum = "OK"
	InstancePoolInstanceLoadBalancerBackendBackendHealthStatusWarning  InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum = "WARNING"
	InstancePoolInstanceLoadBalancerBackendBackendHealthStatusCritical InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum = "CRITICAL"
	InstancePoolInstanceLoadBalancerBackendBackendHealthStatusUnknown  InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum = "UNKNOWN"
)

var mappingInstancePoolInstanceLoadBalancerBackendBackendHealthStatus = map[string]InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum{
	"OK":       InstancePoolInstanceLoadBalancerBackendBackendHealthStatusOk,
	"WARNING":  InstancePoolInstanceLoadBalancerBackendBackendHealthStatusWarning,
	"CRITICAL": InstancePoolInstanceLoadBalancerBackendBackendHealthStatusCritical,
	"UNKNOWN":  InstancePoolInstanceLoadBalancerBackendBackendHealthStatusUnknown,
}

// GetInstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnumValues Enumerates the set of values for InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum
func GetInstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnumValues() []InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum {
	values := make([]InstancePoolInstanceLoadBalancerBackendBackendHealthStatusEnum, 0)
	for _, v := range mappingInstancePoolInstanceLoadBalancerBackendBackendHealthStatus {
		values = append(values, v)
	}
	return values
}
