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

// Service An object that represents one or multiple Oracle services that you can enable for a
// ServiceGateway. In the User Guide topic
// Access to Oracle Services: Service Gateway (https://docs.cloud.oracle.com/Content/Network/Tasks/servicegateway.htm), the
// term *service CIDR label* is used to refer to the string that represents the regional public
// IP address ranges of the Oracle service or services covered by a given `Service` object. That
// unique string is the value of the `Service` object's `cidrBlock` attribute.
type Service struct {

	// A string that represents the regional public IP address ranges for the Oracle service or
	// services covered by this `Service` object. Also known as the `Service` object's *service
	// CIDR label*.
	// When you set up a route rule to route traffic to the service gateway, use this value as the
	// rule's destination. See RouteTable. Also, when you set up
	// a security list rule to cover traffic with the service gateway, use the `cidrBlock` value
	// as the rule's destination (for an egress rule) or the source (for an ingress rule).
	// See SecurityList.
	// Example: `oci-phx-objectstorage`
	CidrBlock *string `mandatory:"true" json:"cidrBlock"`

	// Description of the Oracle service or services covered by this `Service` object.
	// Example: `OCI PHX Object Storage`
	Description *string `mandatory:"true" json:"description"`

	// The `Service` object's OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	Id *string `mandatory:"true" json:"id"`

	// Name of the `Service` object. This name can change and is not guaranteed to be unique.
	// Example: `OCI PHX Object Storage`
	Name *string `mandatory:"true" json:"name"`
}

func (m Service) String() string {
	return common.PointerString(m)
}
