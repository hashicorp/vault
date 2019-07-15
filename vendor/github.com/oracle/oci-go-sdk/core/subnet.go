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

// Subnet A logical subdivision of a VCN. Each subnet exists in a single availability domain and
// consists of a contiguous range of IP addresses that do not overlap with
// other subnets in the VCN. Example: 172.16.1.0/24. For more information, see
// Overview of the Networking Service (https://docs.cloud.oracle.com/Content/Network/Concepts/overview.htm) and
// VCNs and Subnets (https://docs.cloud.oracle.com/Content/Network/Tasks/managingVCNs.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type Subnet struct {

	// The subnet's CIDR block.
	// Example: `172.16.1.0/24`
	CidrBlock *string `mandatory:"true" json:"cidrBlock"`

	// The OCID of the compartment containing the subnet.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The subnet's Oracle ID (OCID).
	Id *string `mandatory:"true" json:"id"`

	// The subnet's current state.
	LifecycleState SubnetLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The OCID of the route table that the subnet uses.
	RouteTableId *string `mandatory:"true" json:"routeTableId"`

	// The OCID of the VCN the subnet is in.
	VcnId *string `mandatory:"true" json:"vcnId"`

	// The IP address of the virtual router.
	// Example: `10.0.14.1`
	VirtualRouterIp *string `mandatory:"true" json:"virtualRouterIp"`

	// The MAC address of the virtual router.
	// Example: `00:00:17:B6:4D:DD`
	VirtualRouterMac *string `mandatory:"true" json:"virtualRouterMac"`

	// The subnet's availability domain. This attribute will be null if this is a regional subnet
	// instead of an AD-specific subnet. Oracle recommends creating regional subnets.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"false" json:"availabilityDomain"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The OCID of the set of DHCP options that the subnet uses.
	DhcpOptionsId *string `mandatory:"false" json:"dhcpOptionsId"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// A DNS label for the subnet, used in conjunction with the VNIC's hostname and
	// VCN's DNS label to form a fully qualified domain name (FQDN) for each VNIC
	// within this subnet (for example, `bminstance-1.subnet123.vcn1.oraclevcn.com`).
	// Must be an alphanumeric string that begins with a letter and is unique within the VCN.
	// The value cannot be changed.
	// The absence of this parameter means the Internet and VCN Resolver
	// will not resolve hostnames of instances in this subnet.
	// For more information, see
	// DNS in Your Virtual Cloud Network (https://docs.cloud.oracle.com/Content/Network/Concepts/dns.htm).
	// Example: `subnet123`
	DnsLabel *string `mandatory:"false" json:"dnsLabel"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Whether VNICs within this subnet can have public IP addresses.
	// Defaults to false, which means VNICs created in this subnet will
	// automatically be assigned public IP addresses unless specified
	// otherwise during instance launch or VNIC creation (with the
	// `assignPublicIp` flag in
	// CreateVnicDetails).
	// If `prohibitPublicIpOnVnic` is set to true, VNICs created in this
	// subnet cannot have public IP addresses (that is, it's a private
	// subnet).
	// Example: `true`
	ProhibitPublicIpOnVnic *bool `mandatory:"false" json:"prohibitPublicIpOnVnic"`

	// The OCIDs of the security list or lists that the subnet uses. Remember
	// that security lists are associated *with the subnet*, but the
	// rules are applied to the individual VNICs in the subnet.
	SecurityListIds []string `mandatory:"false" json:"securityListIds"`

	// The subnet's domain name, which consists of the subnet's DNS label,
	// the VCN's DNS label, and the `oraclevcn.com` domain.
	// For more information, see
	// DNS in Your Virtual Cloud Network (https://docs.cloud.oracle.com/Content/Network/Concepts/dns.htm).
	// Example: `subnet123.vcn1.oraclevcn.com`
	SubnetDomainName *string `mandatory:"false" json:"subnetDomainName"`

	// The date and time the subnet was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m Subnet) String() string {
	return common.PointerString(m)
}

// SubnetLifecycleStateEnum Enum with underlying type: string
type SubnetLifecycleStateEnum string

// Set of constants representing the allowable values for SubnetLifecycleStateEnum
const (
	SubnetLifecycleStateProvisioning SubnetLifecycleStateEnum = "PROVISIONING"
	SubnetLifecycleStateAvailable    SubnetLifecycleStateEnum = "AVAILABLE"
	SubnetLifecycleStateTerminating  SubnetLifecycleStateEnum = "TERMINATING"
	SubnetLifecycleStateTerminated   SubnetLifecycleStateEnum = "TERMINATED"
)

var mappingSubnetLifecycleState = map[string]SubnetLifecycleStateEnum{
	"PROVISIONING": SubnetLifecycleStateProvisioning,
	"AVAILABLE":    SubnetLifecycleStateAvailable,
	"TERMINATING":  SubnetLifecycleStateTerminating,
	"TERMINATED":   SubnetLifecycleStateTerminated,
}

// GetSubnetLifecycleStateEnumValues Enumerates the set of values for SubnetLifecycleStateEnum
func GetSubnetLifecycleStateEnumValues() []SubnetLifecycleStateEnum {
	values := make([]SubnetLifecycleStateEnum, 0)
	for _, v := range mappingSubnetLifecycleState {
		values = append(values, v)
	}
	return values
}
