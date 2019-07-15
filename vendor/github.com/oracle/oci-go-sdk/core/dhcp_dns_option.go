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
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// DhcpDnsOption DHCP option for specifying how DNS (hostname resolution) is handled in the subnets in the VCN.
// For more information, see
// DNS in Your Virtual Cloud Network (https://docs.cloud.oracle.com/Content/Network/Concepts/dns.htm).
type DhcpDnsOption struct {

	// If you set `serverType` to `CustomDnsServer`, specify the
	// IP address of at least one DNS server of your choice (three maximum).
	CustomDnsServers []string `mandatory:"false" json:"customDnsServers"`

	// * **VcnLocal:** Reserved for future use.
	// * **VcnLocalPlusInternet:** Also referred to as "Internet and VCN Resolver".
	// Instances can resolve internet hostnames (no internet gateway is required),
	// and can resolve hostnames of instances in the VCN. This is the default
	// value in the default set of DHCP options in the VCN. For the Internet and
	// VCN Resolver to work across the VCN, there must also be a DNS label set for
	// the VCN, a DNS label set for each subnet, and a hostname for each instance.
	// The Internet and VCN Resolver also enables reverse DNS lookup, which lets
	// you determine the hostname corresponding to the private IP address. For more
	// information, see
	// DNS in Your Virtual Cloud Network (https://docs.cloud.oracle.com/Content/Network/Concepts/dns.htm).
	// * **CustomDnsServer:** Instances use a DNS server of your choice (three
	// maximum).
	ServerType DhcpDnsOptionServerTypeEnum `mandatory:"true" json:"serverType"`
}

func (m DhcpDnsOption) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m DhcpDnsOption) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeDhcpDnsOption DhcpDnsOption
	s := struct {
		DiscriminatorParam string `json:"type"`
		MarshalTypeDhcpDnsOption
	}{
		"DomainNameServer",
		(MarshalTypeDhcpDnsOption)(m),
	}

	return json.Marshal(&s)
}

// DhcpDnsOptionServerTypeEnum Enum with underlying type: string
type DhcpDnsOptionServerTypeEnum string

// Set of constants representing the allowable values for DhcpDnsOptionServerTypeEnum
const (
	DhcpDnsOptionServerTypeVcnlocal             DhcpDnsOptionServerTypeEnum = "VcnLocal"
	DhcpDnsOptionServerTypeVcnlocalplusinternet DhcpDnsOptionServerTypeEnum = "VcnLocalPlusInternet"
	DhcpDnsOptionServerTypeCustomdnsserver      DhcpDnsOptionServerTypeEnum = "CustomDnsServer"
)

var mappingDhcpDnsOptionServerType = map[string]DhcpDnsOptionServerTypeEnum{
	"VcnLocal":             DhcpDnsOptionServerTypeVcnlocal,
	"VcnLocalPlusInternet": DhcpDnsOptionServerTypeVcnlocalplusinternet,
	"CustomDnsServer":      DhcpDnsOptionServerTypeCustomdnsserver,
}

// GetDhcpDnsOptionServerTypeEnumValues Enumerates the set of values for DhcpDnsOptionServerTypeEnum
func GetDhcpDnsOptionServerTypeEnumValues() []DhcpDnsOptionServerTypeEnum {
	values := make([]DhcpDnsOptionServerTypeEnum, 0)
	for _, v := range mappingDhcpDnsOptionServerType {
		values = append(values, v)
	}
	return values
}
