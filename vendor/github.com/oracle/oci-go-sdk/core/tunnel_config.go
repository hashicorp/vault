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

// TunnelConfig Deprecated. For tunnel information, instead see:
//   * IPSecConnectionTunnel
//   * IPSecConnectionTunnelSharedSecret
type TunnelConfig struct {

	// The IP address of Oracle's VPN headend.
	// Example: `129.146.17.50`
	IpAddress *string `mandatory:"true" json:"ipAddress"`

	// The shared secret of the IPSec tunnel.
	// Example: `EXAMPLEToUis6j1cp8GdVQxcmdfMO0yXMLilZTbYCMDGu4V8o`
	SharedSecret *string `mandatory:"true" json:"sharedSecret"`

	// The date and time the IPSec connection was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m TunnelConfig) String() string {
	return common.PointerString(m)
}
