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

// IpAddress A load balancer IP address.
type IpAddress struct {

	// An IP address.
	// Example: `192.168.0.3`
	IpAddress *string `mandatory:"true" json:"ipAddress"`

	// Whether the IP address is public or private.
	// If "true", the IP address is public and accessible from the internet.
	// If "false", the IP address is private and accessible only from within the associated VCN.
	IsPublic *bool `mandatory:"false" json:"isPublic"`
}

func (m IpAddress) String() string {
	return common.PointerString(m)
}
