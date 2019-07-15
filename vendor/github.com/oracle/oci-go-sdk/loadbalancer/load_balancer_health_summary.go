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

// LoadBalancerHealthSummary A health status summary for the specified load balancer.
type LoadBalancerHealthSummary struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the load balancer the health status is associated with.
	LoadBalancerId *string `mandatory:"true" json:"loadBalancerId"`

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
	Status LoadBalancerHealthSummaryStatusEnum `mandatory:"true" json:"status"`
}

func (m LoadBalancerHealthSummary) String() string {
	return common.PointerString(m)
}

// LoadBalancerHealthSummaryStatusEnum Enum with underlying type: string
type LoadBalancerHealthSummaryStatusEnum string

// Set of constants representing the allowable values for LoadBalancerHealthSummaryStatusEnum
const (
	LoadBalancerHealthSummaryStatusOk       LoadBalancerHealthSummaryStatusEnum = "OK"
	LoadBalancerHealthSummaryStatusWarning  LoadBalancerHealthSummaryStatusEnum = "WARNING"
	LoadBalancerHealthSummaryStatusCritical LoadBalancerHealthSummaryStatusEnum = "CRITICAL"
	LoadBalancerHealthSummaryStatusUnknown  LoadBalancerHealthSummaryStatusEnum = "UNKNOWN"
)

var mappingLoadBalancerHealthSummaryStatus = map[string]LoadBalancerHealthSummaryStatusEnum{
	"OK":       LoadBalancerHealthSummaryStatusOk,
	"WARNING":  LoadBalancerHealthSummaryStatusWarning,
	"CRITICAL": LoadBalancerHealthSummaryStatusCritical,
	"UNKNOWN":  LoadBalancerHealthSummaryStatusUnknown,
}

// GetLoadBalancerHealthSummaryStatusEnumValues Enumerates the set of values for LoadBalancerHealthSummaryStatusEnum
func GetLoadBalancerHealthSummaryStatusEnumValues() []LoadBalancerHealthSummaryStatusEnum {
	values := make([]LoadBalancerHealthSummaryStatusEnum, 0)
	for _, v := range mappingLoadBalancerHealthSummaryStatus {
		values = append(values, v)
	}
	return values
}
