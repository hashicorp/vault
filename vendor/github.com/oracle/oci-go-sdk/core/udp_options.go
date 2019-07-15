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

// UdpOptions Optional object to specify ports for a UDP rule. If you specify UDP as the
// protocol but omit this object, then all ports are allowed.
type UdpOptions struct {

	// An inclusive range of allowed destination ports. Use the same number for the min and max
	// to indicate a single port. Defaults to all ports if not specified.
	DestinationPortRange *PortRange `mandatory:"false" json:"destinationPortRange"`

	// An inclusive range of allowed source ports. Use the same number for the min and max to
	// indicate a single port. Defaults to all ports if not specified.
	SourcePortRange *PortRange `mandatory:"false" json:"sourcePortRange"`
}

func (m UdpOptions) String() string {
	return common.PointerString(m)
}
