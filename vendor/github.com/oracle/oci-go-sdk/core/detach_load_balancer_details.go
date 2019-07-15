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

// DetachLoadBalancerDetails Represents a load balancer that is to be detached from an instance pool.
type DetachLoadBalancerDetails struct {

	// The OCID of the load balancer to detach from the instance pool.
	LoadBalancerId *string `mandatory:"true" json:"loadBalancerId"`

	// The name of the backend set on the load balancer to detach from the instance pool.
	BackendSetName *string `mandatory:"true" json:"backendSetName"`
}

func (m DetachLoadBalancerDetails) String() string {
	return common.PointerString(m)
}
