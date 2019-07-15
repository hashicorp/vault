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

// IpSecConnectionDeviceConfig Deprecated. For tunnel information, instead see:
//   * IPSecConnectionTunnel
//   * IPSecConnectionTunnelSharedSecret
type IpSecConnectionDeviceConfig struct {

	// The OCID of the compartment containing the IPSec connection.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The IPSec connection's Oracle ID (OCID).
	Id *string `mandatory:"true" json:"id"`

	// The date and time the IPSec connection was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// Two TunnelConfig objects.
	Tunnels []TunnelConfig `mandatory:"false" json:"tunnels"`
}

func (m IpSecConnectionDeviceConfig) String() string {
	return common.PointerString(m)
}
