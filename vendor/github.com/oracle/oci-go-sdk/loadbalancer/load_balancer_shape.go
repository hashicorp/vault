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

// LoadBalancerShape A shape is a template that determines the total pre-provisioned bandwidth (ingress plus egress) for the
// load balancer.
// Note that the pre-provisioned maximum capacity applies to aggregated connections, not to a single client
// attempting to use the full bandwidth.
type LoadBalancerShape struct {

	// The name of the shape.
	// Example: `100Mbps`
	Name *string `mandatory:"true" json:"name"`
}

func (m LoadBalancerShape) String() string {
	return common.PointerString(m)
}
