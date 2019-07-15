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

// UpdatedNetworkSecurityGroupSecurityRules The representation of UpdatedNetworkSecurityGroupSecurityRules
type UpdatedNetworkSecurityGroupSecurityRules struct {

	// The NSG security rules that were updated.
	SecurityRules []SecurityRule `mandatory:"false" json:"securityRules"`
}

func (m UpdatedNetworkSecurityGroupSecurityRules) String() string {
	return common.PointerString(m)
}
