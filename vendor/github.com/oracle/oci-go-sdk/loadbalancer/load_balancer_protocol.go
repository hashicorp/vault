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

// LoadBalancerProtocol A protocol that defines the type of traffic accepted by a listener.
type LoadBalancerProtocol struct {

	// The name of a protocol.
	// Example: 'HTTP'
	Name *string `mandatory:"true" json:"name"`
}

func (m LoadBalancerProtocol) String() string {
	return common.PointerString(m)
}
