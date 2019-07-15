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

// NetworkSecurityGroupVnic Information about a VNIC that belongs to a network security group.
type NetworkSecurityGroupVnic struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the VNIC.
	VnicId *string `mandatory:"true" json:"vnicId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the parent resource that the VNIC
	// is attached to (for example, a Compute instance).
	ResourceId *string `mandatory:"false" json:"resourceId"`

	// The date and time the VNIC was added to the network security group, in the format
	// defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeAssociated *common.SDKTime `mandatory:"false" json:"timeAssociated"`
}

func (m NetworkSecurityGroupVnic) String() string {
	return common.PointerString(m)
}
