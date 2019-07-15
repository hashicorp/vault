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

// InternetGateway Represents a router that connects the edge of a VCN with the Internet. For an example scenario
// that uses an internet gateway, see
// Typical Networking Service Scenarios (https://docs.cloud.oracle.com/Content/Network/Concepts/overview.htm#scenarios).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type InternetGateway struct {

	// The OCID of the compartment containing the internet gateway.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The internet gateway's Oracle ID (OCID).
	Id *string `mandatory:"true" json:"id"`

	// The internet gateway's current state.
	LifecycleState InternetGatewayLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The OCID of the VCN the internet gateway belongs to.
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

	// Whether the gateway is enabled. When the gateway is disabled, traffic is not
	// routed to/from the Internet, regardless of route rules.
	IsEnabled *bool `mandatory:"false" json:"isEnabled"`

	// The date and time the internet gateway was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m InternetGateway) String() string {
	return common.PointerString(m)
}

// InternetGatewayLifecycleStateEnum Enum with underlying type: string
type InternetGatewayLifecycleStateEnum string

// Set of constants representing the allowable values for InternetGatewayLifecycleStateEnum
const (
	InternetGatewayLifecycleStateProvisioning InternetGatewayLifecycleStateEnum = "PROVISIONING"
	InternetGatewayLifecycleStateAvailable    InternetGatewayLifecycleStateEnum = "AVAILABLE"
	InternetGatewayLifecycleStateTerminating  InternetGatewayLifecycleStateEnum = "TERMINATING"
	InternetGatewayLifecycleStateTerminated   InternetGatewayLifecycleStateEnum = "TERMINATED"
)

var mappingInternetGatewayLifecycleState = map[string]InternetGatewayLifecycleStateEnum{
	"PROVISIONING": InternetGatewayLifecycleStateProvisioning,
	"AVAILABLE":    InternetGatewayLifecycleStateAvailable,
	"TERMINATING":  InternetGatewayLifecycleStateTerminating,
	"TERMINATED":   InternetGatewayLifecycleStateTerminated,
}

// GetInternetGatewayLifecycleStateEnumValues Enumerates the set of values for InternetGatewayLifecycleStateEnum
func GetInternetGatewayLifecycleStateEnumValues() []InternetGatewayLifecycleStateEnum {
	values := make([]InternetGatewayLifecycleStateEnum, 0)
	for _, v := range mappingInternetGatewayLifecycleState {
		values = append(values, v)
	}
	return values
}
