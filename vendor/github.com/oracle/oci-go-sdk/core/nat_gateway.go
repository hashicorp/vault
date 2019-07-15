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

// NatGateway A NAT (Network Address Translation) gateway, which represents a router that lets instances
// without public IPs contact the public internet without exposing the instance to inbound
// internet traffic. For more information, see
// NAT Gateway (https://docs.cloud.oracle.com/Content/Network/Tasks/NATgateway.htm).
// To use any of the API operations, you must be authorized in an
// IAM policy. If you are not authorized, talk to an
// administrator. If you are an administrator who needs to write
// policies to give users access, see Getting Started with
// Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type NatGateway struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment that contains
	// the NAT gateway.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the NAT gateway.
	Id *string `mandatory:"true" json:"id"`

	// Whether the NAT gateway blocks traffic through it. The default is `false`.
	// Example: `true`
	BlockTraffic *bool `mandatory:"true" json:"blockTraffic"`

	// The NAT gateway's current state.
	LifecycleState NatGatewayLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The IP address associated with the NAT gateway.
	NatIp *string `mandatory:"true" json:"natIp"`

	// The date and time the NAT gateway was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the VCN the NAT gateway
	// belongs to.
	VcnId *string `mandatory:"true" json:"vcnId"`

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
}

func (m NatGateway) String() string {
	return common.PointerString(m)
}

// NatGatewayLifecycleStateEnum Enum with underlying type: string
type NatGatewayLifecycleStateEnum string

// Set of constants representing the allowable values for NatGatewayLifecycleStateEnum
const (
	NatGatewayLifecycleStateProvisioning NatGatewayLifecycleStateEnum = "PROVISIONING"
	NatGatewayLifecycleStateAvailable    NatGatewayLifecycleStateEnum = "AVAILABLE"
	NatGatewayLifecycleStateTerminating  NatGatewayLifecycleStateEnum = "TERMINATING"
	NatGatewayLifecycleStateTerminated   NatGatewayLifecycleStateEnum = "TERMINATED"
)

var mappingNatGatewayLifecycleState = map[string]NatGatewayLifecycleStateEnum{
	"PROVISIONING": NatGatewayLifecycleStateProvisioning,
	"AVAILABLE":    NatGatewayLifecycleStateAvailable,
	"TERMINATING":  NatGatewayLifecycleStateTerminating,
	"TERMINATED":   NatGatewayLifecycleStateTerminated,
}

// GetNatGatewayLifecycleStateEnumValues Enumerates the set of values for NatGatewayLifecycleStateEnum
func GetNatGatewayLifecycleStateEnumValues() []NatGatewayLifecycleStateEnum {
	values := make([]NatGatewayLifecycleStateEnum, 0)
	for _, v := range mappingNatGatewayLifecycleState {
		values = append(values, v)
	}
	return values
}
