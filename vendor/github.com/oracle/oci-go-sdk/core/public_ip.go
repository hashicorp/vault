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

// PublicIp A *public IP* is a conceptual term that refers to a public IP address and related properties.
// The `publicIp` object is the API representation of a public IP.
// There are two types of public IPs:
// 1. Ephemeral
// 2. Reserved
// For more information and comparison of the two types,
// see Public IP Addresses (https://docs.cloud.oracle.com/Content/Network/Tasks/managingpublicIPs.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type PublicIp struct {

	// The OCID of the entity the public IP is assigned to, or in the process of
	// being assigned to.
	AssignedEntityId *string `mandatory:"false" json:"assignedEntityId"`

	// The type of entity the public IP is assigned to, or in the process of being
	// assigned to.
	AssignedEntityType PublicIpAssignedEntityTypeEnum `mandatory:"false" json:"assignedEntityType,omitempty"`

	// The public IP's availability domain. This property is set only for ephemeral public IPs
	// that are assigned to a private IP (that is, when the `scope` of the public IP is set to
	// AVAILABILITY_DOMAIN). The value is the availability domain of the assigned private IP.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"false" json:"availabilityDomain"`

	// The OCID of the compartment containing the public IP. For an ephemeral public IP, this is
	// the compartment of its assigned entity (which can be a private IP or a regional entity such
	// as a NAT gateway). For a reserved public IP that is currently assigned,
	// its compartment can be different from the assigned private IP's.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

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

	// The public IP's Oracle ID (OCID).
	Id *string `mandatory:"false" json:"id"`

	// The public IP address of the `publicIp` object.
	// Example: `129.146.2.1`
	IpAddress *string `mandatory:"false" json:"ipAddress"`

	// The public IP's current state.
	LifecycleState PublicIpLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// Defines when the public IP is deleted and released back to Oracle's public IP pool.
	// * `EPHEMERAL`: The lifetime is tied to the lifetime of its assigned entity. An ephemeral
	// public IP must always be assigned to an entity. If the assigned entity is a private IP,
	// the ephemeral public IP is automatically deleted when the private IP is deleted, when
	// the VNIC is terminated, or when the instance is terminated. If the assigned entity is a
	// NatGateway, the ephemeral public IP is automatically
	// deleted when the NAT gateway is terminated.
	// * `RESERVED`: You control the public IP's lifetime. You can delete a reserved public IP
	// whenever you like. It does not need to be assigned to a private IP at all times.
	// For more information and comparison of the two types,
	// see Public IP Addresses (https://docs.cloud.oracle.com/Content/Network/Tasks/managingpublicIPs.htm).
	Lifetime PublicIpLifetimeEnum `mandatory:"false" json:"lifetime,omitempty"`

	// Deprecated. Use `assignedEntityId` instead.
	// The OCID of the private IP that the public IP is currently assigned to, or in the
	// process of being assigned to.
	// **Note:** This is `null` if the public IP is not assigned to a private IP, or is
	// in the process of being assigned to one.
	PrivateIpId *string `mandatory:"false" json:"privateIpId"`

	// Whether the public IP is regional or specific to a particular availability domain.
	// * `REGION`: The public IP exists within a region and is assigned to a regional entity
	// (such as a NatGateway), or can be assigned to a private
	// IP in any availability domain in the region. Reserved public IPs and ephemeral public IPs
	// assigned to a regional entity have `scope` = `REGION`.
	// * `AVAILABILITY_DOMAIN`: The public IP exists within the availability domain of the entity
	// it's assigned to, which is specified by the `availabilityDomain` property of the public IP object.
	// Ephemeral public IPs that are assigned to private IPs have `scope` = `AVAILABILITY_DOMAIN`.
	Scope PublicIpScopeEnum `mandatory:"false" json:"scope,omitempty"`

	// The date and time the public IP was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m PublicIp) String() string {
	return common.PointerString(m)
}

// PublicIpAssignedEntityTypeEnum Enum with underlying type: string
type PublicIpAssignedEntityTypeEnum string

// Set of constants representing the allowable values for PublicIpAssignedEntityTypeEnum
const (
	PublicIpAssignedEntityTypePrivateIp  PublicIpAssignedEntityTypeEnum = "PRIVATE_IP"
	PublicIpAssignedEntityTypeNatGateway PublicIpAssignedEntityTypeEnum = "NAT_GATEWAY"
)

var mappingPublicIpAssignedEntityType = map[string]PublicIpAssignedEntityTypeEnum{
	"PRIVATE_IP":  PublicIpAssignedEntityTypePrivateIp,
	"NAT_GATEWAY": PublicIpAssignedEntityTypeNatGateway,
}

// GetPublicIpAssignedEntityTypeEnumValues Enumerates the set of values for PublicIpAssignedEntityTypeEnum
func GetPublicIpAssignedEntityTypeEnumValues() []PublicIpAssignedEntityTypeEnum {
	values := make([]PublicIpAssignedEntityTypeEnum, 0)
	for _, v := range mappingPublicIpAssignedEntityType {
		values = append(values, v)
	}
	return values
}

// PublicIpLifecycleStateEnum Enum with underlying type: string
type PublicIpLifecycleStateEnum string

// Set of constants representing the allowable values for PublicIpLifecycleStateEnum
const (
	PublicIpLifecycleStateProvisioning PublicIpLifecycleStateEnum = "PROVISIONING"
	PublicIpLifecycleStateAvailable    PublicIpLifecycleStateEnum = "AVAILABLE"
	PublicIpLifecycleStateAssigning    PublicIpLifecycleStateEnum = "ASSIGNING"
	PublicIpLifecycleStateAssigned     PublicIpLifecycleStateEnum = "ASSIGNED"
	PublicIpLifecycleStateUnassigning  PublicIpLifecycleStateEnum = "UNASSIGNING"
	PublicIpLifecycleStateUnassigned   PublicIpLifecycleStateEnum = "UNASSIGNED"
	PublicIpLifecycleStateTerminating  PublicIpLifecycleStateEnum = "TERMINATING"
	PublicIpLifecycleStateTerminated   PublicIpLifecycleStateEnum = "TERMINATED"
)

var mappingPublicIpLifecycleState = map[string]PublicIpLifecycleStateEnum{
	"PROVISIONING": PublicIpLifecycleStateProvisioning,
	"AVAILABLE":    PublicIpLifecycleStateAvailable,
	"ASSIGNING":    PublicIpLifecycleStateAssigning,
	"ASSIGNED":     PublicIpLifecycleStateAssigned,
	"UNASSIGNING":  PublicIpLifecycleStateUnassigning,
	"UNASSIGNED":   PublicIpLifecycleStateUnassigned,
	"TERMINATING":  PublicIpLifecycleStateTerminating,
	"TERMINATED":   PublicIpLifecycleStateTerminated,
}

// GetPublicIpLifecycleStateEnumValues Enumerates the set of values for PublicIpLifecycleStateEnum
func GetPublicIpLifecycleStateEnumValues() []PublicIpLifecycleStateEnum {
	values := make([]PublicIpLifecycleStateEnum, 0)
	for _, v := range mappingPublicIpLifecycleState {
		values = append(values, v)
	}
	return values
}

// PublicIpLifetimeEnum Enum with underlying type: string
type PublicIpLifetimeEnum string

// Set of constants representing the allowable values for PublicIpLifetimeEnum
const (
	PublicIpLifetimeEphemeral PublicIpLifetimeEnum = "EPHEMERAL"
	PublicIpLifetimeReserved  PublicIpLifetimeEnum = "RESERVED"
)

var mappingPublicIpLifetime = map[string]PublicIpLifetimeEnum{
	"EPHEMERAL": PublicIpLifetimeEphemeral,
	"RESERVED":  PublicIpLifetimeReserved,
}

// GetPublicIpLifetimeEnumValues Enumerates the set of values for PublicIpLifetimeEnum
func GetPublicIpLifetimeEnumValues() []PublicIpLifetimeEnum {
	values := make([]PublicIpLifetimeEnum, 0)
	for _, v := range mappingPublicIpLifetime {
		values = append(values, v)
	}
	return values
}

// PublicIpScopeEnum Enum with underlying type: string
type PublicIpScopeEnum string

// Set of constants representing the allowable values for PublicIpScopeEnum
const (
	PublicIpScopeRegion             PublicIpScopeEnum = "REGION"
	PublicIpScopeAvailabilityDomain PublicIpScopeEnum = "AVAILABILITY_DOMAIN"
)

var mappingPublicIpScope = map[string]PublicIpScopeEnum{
	"REGION":              PublicIpScopeRegion,
	"AVAILABILITY_DOMAIN": PublicIpScopeAvailabilityDomain,
}

// GetPublicIpScopeEnumValues Enumerates the set of values for PublicIpScopeEnum
func GetPublicIpScopeEnumValues() []PublicIpScopeEnum {
	values := make([]PublicIpScopeEnum, 0)
	for _, v := range mappingPublicIpScope {
		values = append(values, v)
	}
	return values
}
