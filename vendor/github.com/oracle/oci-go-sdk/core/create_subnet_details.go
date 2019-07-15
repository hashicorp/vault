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

// CreateSubnetDetails The representation of CreateSubnetDetails
type CreateSubnetDetails struct {

	// The CIDR IP address range of the subnet.
	// Example: `172.16.1.0/24`
	CidrBlock *string `mandatory:"true" json:"cidrBlock"`

	// The OCID of the compartment to contain the subnet.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the VCN to contain the subnet.
	VcnId *string `mandatory:"true" json:"vcnId"`

	// Controls whether the subnet is regional or specific to an availability domain. Oracle
	// recommends creating regional subnets because they're more flexible and make it easier to
	// implement failover across availability domains. Originally, AD-specific subnets were the
	// only kind available to use.
	// To create a regional subnet, omit this attribute. Then any resources later created in this
	// subnet (such as a Compute instance) can be created in any availability domain in the region.
	// To instead create an AD-specific subnet, set this attribute to the availability domain you
	// want this subnet to be in. Then any resources later created in this subnet can only be
	// created in that availability domain.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"false" json:"availabilityDomain"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The OCID of the set of DHCP options the subnet will use. If you don't
	// provide a value, the subnet uses the VCN's default set of DHCP options.
	DhcpOptionsId *string `mandatory:"false" json:"dhcpOptionsId"`

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// A DNS label for the subnet, used in conjunction with the VNIC's hostname and
	// VCN's DNS label to form a fully qualified domain name (FQDN) for each VNIC
	// within this subnet (for example, `bminstance-1.subnet123.vcn1.oraclevcn.com`).
	// Must be an alphanumeric string that begins with a letter and is unique within the VCN.
	// The value cannot be changed.
	// This value must be set if you want to use the Internet and VCN Resolver to resolve the
	// hostnames of instances in the subnet. It can only be set if the VCN itself
	// was created with a DNS label.
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
	// `assignPublicIp` flag in CreateVnicDetails).
	// If `prohibitPublicIpOnVnic` is set to true, VNICs created in this
	// subnet cannot have public IP addresses (that is, it's a private
	// subnet).
	//
	// Example: `true`
	ProhibitPublicIpOnVnic *bool `mandatory:"false" json:"prohibitPublicIpOnVnic"`

	// The OCID of the route table the subnet will use. If you don't provide a value,
	// the subnet uses the VCN's default route table.
	RouteTableId *string `mandatory:"false" json:"routeTableId"`

	// The OCIDs of the security list or lists the subnet will use. If you don't
	// provide a value, the subnet uses the VCN's default security list.
	// Remember that security lists are associated *with the subnet*, but the
	// rules are applied to the individual VNICs in the subnet.
	SecurityListIds []string `mandatory:"false" json:"securityListIds"`
}

func (m CreateSubnetDetails) String() string {
	return common.PointerString(m)
}
