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

// IpSecConnection A connection between a DRG and CPE. This connection consists of multiple IPSec
// tunnels. Creating this connection is one of the steps required when setting up
// an IPSec VPN.
// **Important:**  Each tunnel in an IPSec connection can use either static routing or BGP dynamic
// routing (see the IPSecConnectionTunnel object's
// `routing` attribute). Originally only static routing was supported and
// every IPSec connection was required to have at least one static route configured.
// To maintain backward compatibility in the API when support for BPG dynamic routing was introduced,
// the API accepts an empty list of static routes if you configure both of the IPSec tunnels to use
// BGP dynamic routing. If you switch a tunnel's routing from `BGP` to `STATIC`, you must first
// ensure that the IPSec connection is configured with at least one valid CIDR block static route.
// Oracle uses the IPSec connection's static routes when routing a tunnel's traffic *only*
// if that tunnel's `routing` attribute = `STATIC`. Otherwise the static routes are ignored.
// For more information about the workflow for setting up an IPSec connection, see
// IPSec VPN (https://docs.cloud.oracle.com/Content/Network/Tasks/managingIPsec.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type IpSecConnection struct {

	// The OCID of the compartment containing the IPSec connection.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the Cpe object.
	CpeId *string `mandatory:"true" json:"cpeId"`

	// The OCID of the DRG.
	DrgId *string `mandatory:"true" json:"drgId"`

	// The IPSec connection's Oracle ID (OCID).
	Id *string `mandatory:"true" json:"id"`

	// The IPSec connection's current state.
	LifecycleState IpSecConnectionLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Static routes to the CPE. The CIDR must not be a
	// multicast address or class E address.
	// Used for routing a given IPSec tunnel's traffic only if the tunnel
	// is using static routing. If you configure at least one tunnel to use static routing, then
	// you must provide at least one valid static route. If you configure both
	// tunnels to use BGP dynamic routing, you can provide an empty list for the static routes.
	//
	// Example: `10.0.1.0/24`
	StaticRoutes []string `mandatory:"true" json:"staticRoutes"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Your identifier for your CPE device. Can be either an IP address or a hostname (specifically,
	// the fully qualified domain name (FQDN)). The type of identifier here must correspond
	// to the value for `cpeLocalIdentifierType`.
	// If you don't provide a value when creating the IPSec connection, the `ipAddress` attribute
	// for the Cpe object specified by `cpeId` is used as the `cpeLocalIdentifier`.
	// For information about why you'd provide this value, see
	// If Your CPE Is Behind a NAT Device (https://docs.cloud.oracle.com/Content/Network/Tasks/overviewIPsec.htm#nat).
	// Example IP address: `10.0.3.3`
	// Example hostname: `cpe.example.com`
	CpeLocalIdentifier *string `mandatory:"false" json:"cpeLocalIdentifier"`

	// The type of identifier for your CPE device. The value here must correspond to the value
	// for `cpeLocalIdentifier`.
	CpeLocalIdentifierType IpSecConnectionCpeLocalIdentifierTypeEnum `mandatory:"false" json:"cpeLocalIdentifierType,omitempty"`

	// The date and time the IPSec connection was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m IpSecConnection) String() string {
	return common.PointerString(m)
}

// IpSecConnectionLifecycleStateEnum Enum with underlying type: string
type IpSecConnectionLifecycleStateEnum string

// Set of constants representing the allowable values for IpSecConnectionLifecycleStateEnum
const (
	IpSecConnectionLifecycleStateProvisioning IpSecConnectionLifecycleStateEnum = "PROVISIONING"
	IpSecConnectionLifecycleStateAvailable    IpSecConnectionLifecycleStateEnum = "AVAILABLE"
	IpSecConnectionLifecycleStateTerminating  IpSecConnectionLifecycleStateEnum = "TERMINATING"
	IpSecConnectionLifecycleStateTerminated   IpSecConnectionLifecycleStateEnum = "TERMINATED"
)

var mappingIpSecConnectionLifecycleState = map[string]IpSecConnectionLifecycleStateEnum{
	"PROVISIONING": IpSecConnectionLifecycleStateProvisioning,
	"AVAILABLE":    IpSecConnectionLifecycleStateAvailable,
	"TERMINATING":  IpSecConnectionLifecycleStateTerminating,
	"TERMINATED":   IpSecConnectionLifecycleStateTerminated,
}

// GetIpSecConnectionLifecycleStateEnumValues Enumerates the set of values for IpSecConnectionLifecycleStateEnum
func GetIpSecConnectionLifecycleStateEnumValues() []IpSecConnectionLifecycleStateEnum {
	values := make([]IpSecConnectionLifecycleStateEnum, 0)
	for _, v := range mappingIpSecConnectionLifecycleState {
		values = append(values, v)
	}
	return values
}

// IpSecConnectionCpeLocalIdentifierTypeEnum Enum with underlying type: string
type IpSecConnectionCpeLocalIdentifierTypeEnum string

// Set of constants representing the allowable values for IpSecConnectionCpeLocalIdentifierTypeEnum
const (
	IpSecConnectionCpeLocalIdentifierTypeIpAddress IpSecConnectionCpeLocalIdentifierTypeEnum = "IP_ADDRESS"
	IpSecConnectionCpeLocalIdentifierTypeHostname  IpSecConnectionCpeLocalIdentifierTypeEnum = "HOSTNAME"
)

var mappingIpSecConnectionCpeLocalIdentifierType = map[string]IpSecConnectionCpeLocalIdentifierTypeEnum{
	"IP_ADDRESS": IpSecConnectionCpeLocalIdentifierTypeIpAddress,
	"HOSTNAME":   IpSecConnectionCpeLocalIdentifierTypeHostname,
}

// GetIpSecConnectionCpeLocalIdentifierTypeEnumValues Enumerates the set of values for IpSecConnectionCpeLocalIdentifierTypeEnum
func GetIpSecConnectionCpeLocalIdentifierTypeEnumValues() []IpSecConnectionCpeLocalIdentifierTypeEnum {
	values := make([]IpSecConnectionCpeLocalIdentifierTypeEnum, 0)
	for _, v := range mappingIpSecConnectionCpeLocalIdentifierType {
		values = append(values, v)
	}
	return values
}
