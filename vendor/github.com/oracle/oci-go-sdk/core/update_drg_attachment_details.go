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

// UpdateDrgAttachmentDetails The representation of UpdateDrgAttachmentDetails
type UpdateDrgAttachmentDetails struct {

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The OCID of the route table the DRG attachment will use. For information about why you
	// would associate a route table with a DRG attachment, see
	// Advanced Scenario: Transit Routing (https://docs.cloud.oracle.com/Content/Network/Tasks/transitrouting.htm).
	RouteTableId *string `mandatory:"false" json:"routeTableId"`
}

func (m UpdateDrgAttachmentDetails) String() string {
	return common.PointerString(m)
}
