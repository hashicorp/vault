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

// AttachLoadBalancerDetails Represents a load balancer that is to be attached to an instance pool.
type AttachLoadBalancerDetails struct {

	// The OCID of the load balancer to attach to the instance pool.
	LoadBalancerId *string `mandatory:"true" json:"loadBalancerId"`

	// The name of the backend set on the load balancer to add instances to.
	BackendSetName *string `mandatory:"true" json:"backendSetName"`

	// The port value to use when creating the backend set.
	Port *int `mandatory:"true" json:"port"`

	// Indicates which VNIC on each instance in the pool should be used to associate with the load balancer. Possible values are "PrimaryVnic" or the displayName of one of the secondary VNICs on the instance configuration that is associated with the instance pool.
	VnicSelection *string `mandatory:"true" json:"vnicSelection"`
}

func (m AttachLoadBalancerDetails) String() string {
	return common.PointerString(m)
}
