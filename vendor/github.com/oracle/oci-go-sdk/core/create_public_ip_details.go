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

// CreatePublicIpDetails The representation of CreatePublicIpDetails
type CreatePublicIpDetails struct {

	// The OCID of the compartment to contain the public IP. For ephemeral public IPs,
	// you must set this to the private IP's compartment OCID.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// Defines when the public IP is deleted and released back to the Oracle Cloud
	// Infrastructure public IP pool. For more information, see
	// Public IP Addresses (https://docs.cloud.oracle.com/Content/Network/Tasks/managingpublicIPs.htm).
	Lifetime CreatePublicIpDetailsLifetimeEnum `mandatory:"true" json:"lifetime"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid
	// entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The OCID of the private IP to assign the public IP to.
	// Required for an ephemeral public IP because it must always be assigned to a private IP
	// (specifically a *primary* private IP).
	// Optional for a reserved public IP. If you don't provide it, the public IP is created but not
	// assigned to a private IP. You can later assign the public IP with
	// UpdatePublicIp.
	PrivateIpId *string `mandatory:"false" json:"privateIpId"`
}

func (m CreatePublicIpDetails) String() string {
	return common.PointerString(m)
}

// CreatePublicIpDetailsLifetimeEnum Enum with underlying type: string
type CreatePublicIpDetailsLifetimeEnum string

// Set of constants representing the allowable values for CreatePublicIpDetailsLifetimeEnum
const (
	CreatePublicIpDetailsLifetimeEphemeral CreatePublicIpDetailsLifetimeEnum = "EPHEMERAL"
	CreatePublicIpDetailsLifetimeReserved  CreatePublicIpDetailsLifetimeEnum = "RESERVED"
)

var mappingCreatePublicIpDetailsLifetime = map[string]CreatePublicIpDetailsLifetimeEnum{
	"EPHEMERAL": CreatePublicIpDetailsLifetimeEphemeral,
	"RESERVED":  CreatePublicIpDetailsLifetimeReserved,
}

// GetCreatePublicIpDetailsLifetimeEnumValues Enumerates the set of values for CreatePublicIpDetailsLifetimeEnum
func GetCreatePublicIpDetailsLifetimeEnumValues() []CreatePublicIpDetailsLifetimeEnum {
	values := make([]CreatePublicIpDetailsLifetimeEnum, 0)
	for _, v := range mappingCreatePublicIpDetailsLifetime {
		values = append(values, v)
	}
	return values
}
