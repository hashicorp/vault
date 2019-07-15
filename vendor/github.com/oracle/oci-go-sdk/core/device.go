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

// Device Device Path corresponding to the block devices attached to instances having a name and isAvailable flag.
type Device struct {

	// The device name.
	Name *string `mandatory:"true" json:"name"`

	// The flag denoting whether device is available.
	IsAvailable *bool `mandatory:"true" json:"isAvailable"`
}

func (m Device) String() string {
	return common.PointerString(m)
}
