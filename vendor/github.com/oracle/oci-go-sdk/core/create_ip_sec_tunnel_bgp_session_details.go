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

// CreateIpSecTunnelBgpSessionDetails The representation of CreateIpSecTunnelBgpSessionDetails
type CreateIpSecTunnelBgpSessionDetails struct {

	// The IP address for the Oracle end of the inside tunnel interface.
	// If the tunnel's `routing` attribute is set to `BGP`
	// (see IPSecConnectionTunnel), this IP address
	// is required and used for the tunnel's BGP session.
	// If `routing` is instead set to `STATIC`, this IP address is optional. You can set this IP
	// address to troubleshoot or monitor the tunnel.
	// The value must be a /30 or /31.
	// Example: `10.0.0.4/31`
	OracleInterfaceIp *string `mandatory:"false" json:"oracleInterfaceIp"`

	// The IP address for the CPE end of the inside tunnel interface.
	// If the tunnel's `routing` attribute is set to `BGP`
	// (see IPSecConnectionTunnel), this IP address
	// is required and used for the tunnel's BGP session.
	// If `routing` is instead set to `STATIC`, this IP address is optional. You can set this IP
	// address to troubleshoot or monitor the tunnel.
	// The value must be a /30 or /31.
	// Example: `10.0.0.5/31`
	CustomerInterfaceIp *string `mandatory:"false" json:"customerInterfaceIp"`

	// If the tunnel's `routing` attribute is set to `BGP`
	// (see IPSecConnectionTunnel), this ASN
	// is required and used for the tunnel's BGP session. This is the ASN of the network on the
	// CPE end of the BGP session. Can be a 2-byte or 4-byte ASN. Uses "asplain" format.
	// If the tunnel's `routing` attribute is set to `STATIC`, the `customerBgpAsn` must be null.
	// Example: `12345` (2-byte) or `1587232876` (4-byte)
	CustomerBgpAsn *string `mandatory:"false" json:"customerBgpAsn"`
}

func (m CreateIpSecTunnelBgpSessionDetails) String() string {
	return common.PointerString(m)
}
