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

// UpdateIpSecTunnelBgpSessionDetails The representation of UpdateIpSecTunnelBgpSessionDetails
type UpdateIpSecTunnelBgpSessionDetails struct {

	// The IP address for the Oracle end of the inside tunnel interface.
	// If the tunnel's `routing` attribute is set to `BGP`
	// (see UpdateIPSecConnectionTunnelDetails), this IP address
	// is used for the tunnel's BGP session.
	// If `routing` is instead set to `STATIC`, you can set this IP address to troubleshoot or
	// monitor the tunnel.
	// The value must be a /30 or /31.
	// If you are switching the tunnel from using BGP dynamic routing to static routing and want
	// to remove the value for `oracleInterfaceIp`, you can set the value to an empty string.
	// Example: `10.0.0.4/31`
	OracleInterfaceIp *string `mandatory:"false" json:"oracleInterfaceIp"`

	// The IP address for the CPE end of the inside tunnel interface.
	// If the tunnel's `routing` attribute is set to `BGP`
	// (see UpdateIPSecConnectionTunnelDetails), this IP address
	// is used for the tunnel's BGP session.
	// If `routing` is instead set to `STATIC`, you can set this IP address to troubleshoot or
	// monitor the tunnel.
	// The value must be a /30 or /31.
	// If you are switching the tunnel from using BGP dynamic routing to static routing and want
	// to remove the value for `customerInterfaceIp`, you can set the value to an empty string.
	// Example: `10.0.0.5/31`
	CustomerInterfaceIp *string `mandatory:"false" json:"customerInterfaceIp"`

	// The BGP ASN of the network on the CPE end of the BGP session. Can be a 2-byte or 4-byte ASN.
	// Uses "asplain" format.
	// If you are switching the tunnel from using BGP dynamic routing to static routing, the
	// `customerBgpAsn` must be null.
	// Example: `12345` (2-byte) or `1587232876` (4-byte)
	CustomerBgpAsn *string `mandatory:"false" json:"customerBgpAsn"`
}

func (m UpdateIpSecTunnelBgpSessionDetails) String() string {
	return common.PointerString(m)
}
