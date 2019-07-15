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

// PeerRegionForRemotePeering Details about a region that supports remote VCN peering. For more information, see VCN Peering (https://docs.cloud.oracle.com/Content/Network/Tasks/VCNpeering.htm).
type PeerRegionForRemotePeering struct {

	// The region's name.
	// Example: `us-phoenix-1`
	Name *string `mandatory:"true" json:"name"`
}

func (m PeerRegionForRemotePeering) String() string {
	return common.PointerString(m)
}
