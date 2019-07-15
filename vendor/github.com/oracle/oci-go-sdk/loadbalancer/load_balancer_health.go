// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Load Balancing API
//
// API for the Load Balancing service. Use this API to manage load balancers, backend sets, and related items. For more
// information, see Overview of Load Balancing (https://docs.cloud.oracle.com/iaas/Content/Balance/Concepts/balanceoverview.htm).
//

package loadbalancer

import (
	"github.com/oracle/oci-go-sdk/common"
)

// LoadBalancerHealth The health status details for the specified load balancer.
// This object does not explicitly enumerate backend sets with a status of `OK`. However, they are included in the
// `totalBackendSetCount` sum.
type LoadBalancerHealth struct {

	// The overall health status of the load balancer.
	// *  **OK:** All backend sets associated with the load balancer return a status of `OK`.
	// *  **WARNING:** At least one of the backend sets associated with the load balancer returns a status of `WARNING`,
	// no backend sets return a status of `CRITICAL`, and the load balancer life cycle state is `ACTIVE`.
	// *  **CRITICAL:** One or more of the backend sets associated with the load balancer return a status of `CRITICAL`.
	// *  **UNKNOWN:** If any one of the following conditions is true:
	//     *  The load balancer life cycle state is not `ACTIVE`.
	//     *  No backend sets are defined for the load balancer.
	//     *  More than half of the backend sets associated with the load balancer return a status of `UNKNOWN`, none of the backend
	//        sets return a status of `WARNING` or `CRITICAL`, and the load balancer life cycle state is `ACTIVE`.
	//     *  The system could not retrieve metrics for any reason.
	Status LoadBalancerHealthStatusEnum `mandatory:"true" json:"status"`

	// A list of backend sets that are currently in the `WARNING` health state. The list identifies each backend set by the
	// friendly name you assigned when you created it.
	// Example: `example_backend_set3`
	WarningStateBackendSetNames []string `mandatory:"true" json:"warningStateBackendSetNames"`

	// A list of backend sets that are currently in the `CRITICAL` health state. The list identifies each backend set by the
	// friendly name you assigned when you created it.
	// Example: `example_backend_set`
	CriticalStateBackendSetNames []string `mandatory:"true" json:"criticalStateBackendSetNames"`

	// A list of backend sets that are currently in the `UNKNOWN` health state. The list identifies each backend set by the
	// friendly name you assigned when you created it.
	// Example: `example_backend_set2`
	UnknownStateBackendSetNames []string `mandatory:"true" json:"unknownStateBackendSetNames"`

	// The total number of backend sets associated with this load balancer.
	// Example: `4`
	TotalBackendSetCount *int `mandatory:"true" json:"totalBackendSetCount"`
}

func (m LoadBalancerHealth) String() string {
	return common.PointerString(m)
}

// LoadBalancerHealthStatusEnum Enum with underlying type: string
type LoadBalancerHealthStatusEnum string

// Set of constants representing the allowable values for LoadBalancerHealthStatusEnum
const (
	LoadBalancerHealthStatusOk       LoadBalancerHealthStatusEnum = "OK"
	LoadBalancerHealthStatusWarning  LoadBalancerHealthStatusEnum = "WARNING"
	LoadBalancerHealthStatusCritical LoadBalancerHealthStatusEnum = "CRITICAL"
	LoadBalancerHealthStatusUnknown  LoadBalancerHealthStatusEnum = "UNKNOWN"
)

var mappingLoadBalancerHealthStatus = map[string]LoadBalancerHealthStatusEnum{
	"OK":       LoadBalancerHealthStatusOk,
	"WARNING":  LoadBalancerHealthStatusWarning,
	"CRITICAL": LoadBalancerHealthStatusCritical,
	"UNKNOWN":  LoadBalancerHealthStatusUnknown,
}

// GetLoadBalancerHealthStatusEnumValues Enumerates the set of values for LoadBalancerHealthStatusEnum
func GetLoadBalancerHealthStatusEnumValues() []LoadBalancerHealthStatusEnum {
	values := make([]LoadBalancerHealthStatusEnum, 0)
	for _, v := range mappingLoadBalancerHealthStatus {
		values = append(values, v)
	}
	return values
}
