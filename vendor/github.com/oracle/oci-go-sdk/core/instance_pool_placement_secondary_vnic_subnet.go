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

// InstancePoolPlacementSecondaryVnicSubnet The secondary VNIC object for the placement configuration for an instance pool.
type InstancePoolPlacementSecondaryVnicSubnet struct {

	// The subnet OCID for the secondary vnic
	SubnetId *string `mandatory:"true" json:"subnetId"`

	// The displayName of the vnic. This is also use to match against the Instance Configuration defined
	// secondary vnic.
	DisplayName *string `mandatory:"false" json:"displayName"`
}

func (m InstancePoolPlacementSecondaryVnicSubnet) String() string {
	return common.PointerString(m)
}
