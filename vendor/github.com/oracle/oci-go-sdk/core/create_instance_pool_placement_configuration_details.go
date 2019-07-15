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

// CreateInstancePoolPlacementConfigurationDetails The location for where an instance pool will place instances.
type CreateInstancePoolPlacementConfigurationDetails struct {

	// The availability domain to place instances.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID of the primary subnet to place instances.
	PrimarySubnetId *string `mandatory:"true" json:"primarySubnetId"`

	// The set of secondary VNIC data for instances in the pool.
	SecondaryVnicSubnets []InstancePoolPlacementSecondaryVnicSubnet `mandatory:"false" json:"secondaryVnicSubnets"`
}

func (m CreateInstancePoolPlacementConfigurationDetails) String() string {
	return common.PointerString(m)
}
