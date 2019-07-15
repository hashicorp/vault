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

// SecurityRule A security rule is one of the items in a NetworkSecurityGroup.
// It is a virtual firewall rule for the VNICs in the network security group. A rule can be for
// either inbound (`direction`= INGRESS) or outbound (`direction`= EGRESS) IP packets.
type SecurityRule struct {

	// Direction of the security rule. Set to `EGRESS` for rules to allow outbound IP packets, or `INGRESS` for rules to allow inbound IP packets.
	Direction SecurityRuleDirectionEnum `mandatory:"true" json:"direction"`

	// The transport protocol. Specify either `all` or an IPv4 protocol number as
	// defined in
	// Protocol Numbers (http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml).
	// Options are supported only for ICMP ("1"), TCP ("6"), UDP ("17"), and ICMPv6 ("58").
	Protocol *string `mandatory:"true" json:"protocol"`

	// An optional description of your choice for the rule.
	Description *string `mandatory:"false" json:"description"`

	// Conceptually, this is the range of IP addresses that a packet originating from the instance
	// can go to.
	// Allowed values:
	//   * An IP address range in CIDR notation. For example: `192.168.1.0/24`
	//   * The `cidrBlock` value for a Service, if you're
	//     setting up a security rule for traffic destined for a particular `Service` through
	//     a service gateway. For example: `oci-phx-objectstorage`.
	//   * The OCID of a NetworkSecurityGroup in the same
	//     VCN. The value can be the NSG that the rule belongs to if the rule's intent is to control
	//     traffic between VNICs in the same NSG.
	Destination *string `mandatory:"false" json:"destination"`

	// Type of destination for the rule. Required if `direction` = `EGRESS`.
	// Allowed values:
	//   * `CIDR_BLOCK`: If the rule's `destination` is an IP address range in CIDR notation.
	//   * `SERVICE_CIDR_BLOCK`: If the rule's `destination` is the `cidrBlock` value for a
	//     Service (the rule is for traffic destined for a
	//     particular `Service` through a service gateway).
	//   * `NETWORK_SECURITY_GROUP`: If the rule's `destination` is the OCID of a
	//     NetworkSecurityGroup.
	DestinationType SecurityRuleDestinationTypeEnum `mandatory:"false" json:"destinationType,omitempty"`

	// Optional and valid only for ICMP and ICMPv6. Use to specify a particular ICMP type and code
	// as defined in:
	// - ICMP Parameters (http://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml)
	// - ICMPv6 Parameters (https://www.iana.org/assignments/icmpv6-parameters/icmpv6-parameters.xhtml)
	// If you specify ICMP or ICMPv6 as the protocol but omit this object, then all ICMP types and
	// codes are allowed. If you do provide this object, the type is required and the code is optional.
	// To enable MTU negotiation for ingress internet traffic via IPv4, make sure to allow type 3 ("Destination
	// Unreachable") code 4 ("Fragmentation Needed and Don't Fragment was Set"). If you need to specify
	// multiple codes for a single type, create a separate security rule for each.
	IcmpOptions *IcmpOptions `mandatory:"false" json:"icmpOptions"`

	// An Oracle-assigned identifier for the security rule. You specify this ID when you want to
	// update or delete the rule.
	// Example: `04ABEC`
	Id *string `mandatory:"false" json:"id"`

	// A stateless rule allows traffic in one direction. Remember to add a corresponding
	// stateless rule in the other direction if you need to support bidirectional traffic. For
	// example, if egress traffic allows TCP destination port 80, there should be an ingress
	// rule to allow TCP source port 80. Defaults to false, which means the rule is stateful
	// and a corresponding rule is not necessary for bidirectional traffic.
	IsStateless *bool `mandatory:"false" json:"isStateless"`

	// Whether the rule is valid. The value is `True` when the rule is first created. If
	// the rule's `source` or `destination` is a network security group, the value changes to
	// `False` if that network security group is deleted.
	IsValid *bool `mandatory:"false" json:"isValid"`

	// Conceptually, this is the range of IP addresses that a packet coming into the instance
	// can come from.
	// Allowed values:
	//   * An IP address range in CIDR notation. For example: `192.168.1.0/24`
	//   * The `cidrBlock` value for a Service, if you're
	//     setting up a security rule for traffic coming from a particular `Service` through
	//     a service gateway. For example: `oci-phx-objectstorage`.
	//   * The OCID of a NetworkSecurityGroup in the same
	//     VCN. The value can be the NSG that the rule belongs to if the rule's intent is to control
	//     traffic between VNICs in the same NSG.
	Source *string `mandatory:"false" json:"source"`

	// Type of source for the rule. Required if `direction` = `INGRESS`.
	//   * `CIDR_BLOCK`: If the rule's `source` is an IP address range in CIDR notation.
	//   * `SERVICE_CIDR_BLOCK`: If the rule's `source` is the `cidrBlock` value for a
	//     Service (the rule is for traffic coming from a
	//     particular `Service` through a service gateway).
	//   * `NETWORK_SECURITY_GROUP`: If the rule's `destination` is the OCID of a
	//     NetworkSecurityGroup.
	SourceType SecurityRuleSourceTypeEnum `mandatory:"false" json:"sourceType,omitempty"`

	// Optional and valid only for TCP. Use to specify particular destination ports for TCP rules.
	// If you specify TCP as the protocol but omit this object, then all destination ports are allowed.
	TcpOptions *TcpOptions `mandatory:"false" json:"tcpOptions"`

	// The date and time the security rule was created. Format defined by RFC3339.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// Optional and valid only for UDP. Use to specify particular destination ports for UDP rules.
	// If you specify UDP as the protocol but omit this object, then all destination ports are allowed.
	UdpOptions *UdpOptions `mandatory:"false" json:"udpOptions"`
}

func (m SecurityRule) String() string {
	return common.PointerString(m)
}

// SecurityRuleDestinationTypeEnum Enum with underlying type: string
type SecurityRuleDestinationTypeEnum string

// Set of constants representing the allowable values for SecurityRuleDestinationTypeEnum
const (
	SecurityRuleDestinationTypeCidrBlock            SecurityRuleDestinationTypeEnum = "CIDR_BLOCK"
	SecurityRuleDestinationTypeServiceCidrBlock     SecurityRuleDestinationTypeEnum = "SERVICE_CIDR_BLOCK"
	SecurityRuleDestinationTypeNetworkSecurityGroup SecurityRuleDestinationTypeEnum = "NETWORK_SECURITY_GROUP"
)

var mappingSecurityRuleDestinationType = map[string]SecurityRuleDestinationTypeEnum{
	"CIDR_BLOCK":             SecurityRuleDestinationTypeCidrBlock,
	"SERVICE_CIDR_BLOCK":     SecurityRuleDestinationTypeServiceCidrBlock,
	"NETWORK_SECURITY_GROUP": SecurityRuleDestinationTypeNetworkSecurityGroup,
}

// GetSecurityRuleDestinationTypeEnumValues Enumerates the set of values for SecurityRuleDestinationTypeEnum
func GetSecurityRuleDestinationTypeEnumValues() []SecurityRuleDestinationTypeEnum {
	values := make([]SecurityRuleDestinationTypeEnum, 0)
	for _, v := range mappingSecurityRuleDestinationType {
		values = append(values, v)
	}
	return values
}

// SecurityRuleDirectionEnum Enum with underlying type: string
type SecurityRuleDirectionEnum string

// Set of constants representing the allowable values for SecurityRuleDirectionEnum
const (
	SecurityRuleDirectionEgress  SecurityRuleDirectionEnum = "EGRESS"
	SecurityRuleDirectionIngress SecurityRuleDirectionEnum = "INGRESS"
)

var mappingSecurityRuleDirection = map[string]SecurityRuleDirectionEnum{
	"EGRESS":  SecurityRuleDirectionEgress,
	"INGRESS": SecurityRuleDirectionIngress,
}

// GetSecurityRuleDirectionEnumValues Enumerates the set of values for SecurityRuleDirectionEnum
func GetSecurityRuleDirectionEnumValues() []SecurityRuleDirectionEnum {
	values := make([]SecurityRuleDirectionEnum, 0)
	for _, v := range mappingSecurityRuleDirection {
		values = append(values, v)
	}
	return values
}

// SecurityRuleSourceTypeEnum Enum with underlying type: string
type SecurityRuleSourceTypeEnum string

// Set of constants representing the allowable values for SecurityRuleSourceTypeEnum
const (
	SecurityRuleSourceTypeCidrBlock            SecurityRuleSourceTypeEnum = "CIDR_BLOCK"
	SecurityRuleSourceTypeServiceCidrBlock     SecurityRuleSourceTypeEnum = "SERVICE_CIDR_BLOCK"
	SecurityRuleSourceTypeNetworkSecurityGroup SecurityRuleSourceTypeEnum = "NETWORK_SECURITY_GROUP"
)

var mappingSecurityRuleSourceType = map[string]SecurityRuleSourceTypeEnum{
	"CIDR_BLOCK":             SecurityRuleSourceTypeCidrBlock,
	"SERVICE_CIDR_BLOCK":     SecurityRuleSourceTypeServiceCidrBlock,
	"NETWORK_SECURITY_GROUP": SecurityRuleSourceTypeNetworkSecurityGroup,
}

// GetSecurityRuleSourceTypeEnumValues Enumerates the set of values for SecurityRuleSourceTypeEnum
func GetSecurityRuleSourceTypeEnumValues() []SecurityRuleSourceTypeEnum {
	values := make([]SecurityRuleSourceTypeEnum, 0)
	for _, v := range mappingSecurityRuleSourceType {
		values = append(values, v)
	}
	return values
}
