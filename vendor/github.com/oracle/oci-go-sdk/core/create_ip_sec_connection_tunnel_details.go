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

// CreateIpSecConnectionTunnelDetails The representation of CreateIpSecConnectionTunnelDetails
type CreateIpSecConnectionTunnelDetails struct {

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid
	// entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The type of routing to use for this tunnel (either BGP dynamic routing or static routing).
	Routing CreateIpSecConnectionTunnelDetailsRoutingEnum `mandatory:"false" json:"routing,omitempty"`

	// The shared secret (pre-shared key) to use for the IPSec tunnel. Only numbers, letters, and
	// spaces are allowed. If you don't provide a value,
	// Oracle generates a value for you. You can specify your own shared secret later if
	// you like with UpdateIPSecConnectionTunnelSharedSecret.
	// Example: `EXAMPLEToUis6j1cp8GdVQxcmdfMO0yXMLilZTbYCMDGu4V8o`
	SharedSecret *string `mandatory:"false" json:"sharedSecret"`

	// Information for establishing a BGP session for the IPSec tunnel. Required if the tunnel uses
	// BGP dynamic routing.
	// If the tunnel instead uses static routing, you may optionally provide
	// this object and set an IP address for one or both ends of the IPSec tunnel for the purposes
	// of troubleshooting or monitoring the tunnel.
	BgpSessionConfig *CreateIpSecTunnelBgpSessionDetails `mandatory:"false" json:"bgpSessionConfig"`
}

func (m CreateIpSecConnectionTunnelDetails) String() string {
	return common.PointerString(m)
}

// CreateIpSecConnectionTunnelDetailsRoutingEnum Enum with underlying type: string
type CreateIpSecConnectionTunnelDetailsRoutingEnum string

// Set of constants representing the allowable values for CreateIpSecConnectionTunnelDetailsRoutingEnum
const (
	CreateIpSecConnectionTunnelDetailsRoutingBgp    CreateIpSecConnectionTunnelDetailsRoutingEnum = "BGP"
	CreateIpSecConnectionTunnelDetailsRoutingStatic CreateIpSecConnectionTunnelDetailsRoutingEnum = "STATIC"
)

var mappingCreateIpSecConnectionTunnelDetailsRouting = map[string]CreateIpSecConnectionTunnelDetailsRoutingEnum{
	"BGP":    CreateIpSecConnectionTunnelDetailsRoutingBgp,
	"STATIC": CreateIpSecConnectionTunnelDetailsRoutingStatic,
}

// GetCreateIpSecConnectionTunnelDetailsRoutingEnumValues Enumerates the set of values for CreateIpSecConnectionTunnelDetailsRoutingEnum
func GetCreateIpSecConnectionTunnelDetailsRoutingEnumValues() []CreateIpSecConnectionTunnelDetailsRoutingEnum {
	values := make([]CreateIpSecConnectionTunnelDetailsRoutingEnum, 0)
	for _, v := range mappingCreateIpSecConnectionTunnelDetailsRouting {
		values = append(values, v)
	}
	return values
}
