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

// VirtualCircuitBandwidthShape An individual bandwidth level for virtual circuits.
type VirtualCircuitBandwidthShape struct {

	// The name of the bandwidth shape.
	// Example: `10 Gbps`
	Name *string `mandatory:"true" json:"name"`

	// The bandwidth in Mbps.
	// Example: `10000`
	BandwidthInMbps *int `mandatory:"false" json:"bandwidthInMbps"`
}

func (m VirtualCircuitBandwidthShape) String() string {
	return common.PointerString(m)
}
