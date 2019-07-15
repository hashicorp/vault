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

// CreateIpSecConnectionDetails The representation of CreateIpSecConnectionDetails
type CreateIpSecConnectionDetails struct {

	// The OCID of the compartment to contain the IPSec connection.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the Cpe object.
	CpeId *string `mandatory:"true" json:"cpeId"`

	// The OCID of the DRG.
	DrgId *string `mandatory:"true" json:"drgId"`

	// Static routes to the CPE. A static route's CIDR must not be a
	// multicast address or class E address.
	// Used for routing a given IPSec tunnel's traffic only if the tunnel
	// is using static routing. If you configure at least one tunnel to use static routing, then
	// you must provide at least one valid static route. If you configure both
	// tunnels to use BGP dynamic routing, you can provide an empty list for the static routes.
	// For more information, see the important note in IPSecConnection.
	//
	// Example: `10.0.1.0/24`
	StaticRoutes []string `mandatory:"true" json:"staticRoutes"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Your identifier for your CPE device. Can be either an IP address or a hostname (specifically, the
	// fully qualified domain name (FQDN)). The type of identifier you provide here must correspond
	// to the value for `cpeLocalIdentifierType`.
	// If you don't provide a value, the `ipAddress` attribute for the Cpe
	// object specified by `cpeId` is used as the `cpeLocalIdentifier`.
	// For information about why you'd provide this value, see
	// If Your CPE Is Behind a NAT Device (https://docs.cloud.oracle.com/Content/Network/Tasks/overviewIPsec.htm#nat).
	// Example IP address: `10.0.3.3`
	// Example hostname: `cpe.example.com`
	CpeLocalIdentifier *string `mandatory:"false" json:"cpeLocalIdentifier"`

	// The type of identifier for your CPE device. The value you provide here must correspond to the value
	// for `cpeLocalIdentifier`.
	CpeLocalIdentifierType CreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnum `mandatory:"false" json:"cpeLocalIdentifierType,omitempty"`

	// Information for creating the individual tunnels in the IPSec connection. You can provide a
	// maximum of 2 `tunnelConfiguration` objects in the array (one for each of the
	// two tunnels).
	TunnelConfiguration []CreateIpSecConnectionTunnelDetails `mandatory:"false" json:"tunnelConfiguration"`
}

func (m CreateIpSecConnectionDetails) String() string {
	return common.PointerString(m)
}

// CreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnum Enum with underlying type: string
type CreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnum string

// Set of constants representing the allowable values for CreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnum
const (
	CreateIpSecConnectionDetailsCpeLocalIdentifierTypeIpAddress CreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnum = "IP_ADDRESS"
	CreateIpSecConnectionDetailsCpeLocalIdentifierTypeHostname  CreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnum = "HOSTNAME"
)

var mappingCreateIpSecConnectionDetailsCpeLocalIdentifierType = map[string]CreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnum{
	"IP_ADDRESS": CreateIpSecConnectionDetailsCpeLocalIdentifierTypeIpAddress,
	"HOSTNAME":   CreateIpSecConnectionDetailsCpeLocalIdentifierTypeHostname,
}

// GetCreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnumValues Enumerates the set of values for CreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnum
func GetCreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnumValues() []CreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnum {
	values := make([]CreateIpSecConnectionDetailsCpeLocalIdentifierTypeEnum, 0)
	for _, v := range mappingCreateIpSecConnectionDetailsCpeLocalIdentifierType {
		values = append(values, v)
	}
	return values
}
