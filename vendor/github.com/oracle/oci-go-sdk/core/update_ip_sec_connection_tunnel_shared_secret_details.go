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

// UpdateIpSecConnectionTunnelSharedSecretDetails The representation of UpdateIpSecConnectionTunnelSharedSecretDetails
type UpdateIpSecConnectionTunnelSharedSecretDetails struct {

	// The shared secret (pre-shared key) to use for the tunnel. Only numbers, letters, and spaces
	// are allowed.
	// Example: `EXAMPLEToUis6j1cp8GdVQxcmdfMO0yXMLilZTbYCMDGu4V8o`
	SharedSecret *string `mandatory:"false" json:"sharedSecret"`
}

func (m UpdateIpSecConnectionTunnelSharedSecretDetails) String() string {
	return common.PointerString(m)
}
