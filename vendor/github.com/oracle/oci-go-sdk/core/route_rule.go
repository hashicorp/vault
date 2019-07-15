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

// RouteRule A mapping between a destination IP address range and a virtual device to route matching
// packets to (a target).
type RouteRule struct {

	// The OCID for the route rule's target. For information about the type of
	// targets you can specify, see
	// Route Tables (https://docs.cloud.oracle.com/Content/Network/Tasks/managingroutetables.htm).
	NetworkEntityId *string `mandatory:"true" json:"networkEntityId"`

	// Deprecated. Instead use `destination` and `destinationType`. Requests that include both
	// `cidrBlock` and `destination` will be rejected.
	// A destination IP address range in CIDR notation. Matching packets will
	// be routed to the indicated network entity (the target).
	//
	// Example: `0.0.0.0/0`
	CidrBlock *string `mandatory:"false" json:"cidrBlock"`

	// Conceptually, this is the range of IP addresses used for matching when routing
	// traffic. Required if you provide a `destinationType`.
	// Allowed values:
	//   * IP address range in CIDR notation. For example: `192.168.1.0/24`
	//   * The `cidrBlock` value for a Service, if you're
	//     setting up a route rule for traffic destined for a particular `Service` through
	//     a service gateway. For example: `oci-phx-objectstorage`.
	Destination *string `mandatory:"false" json:"destination"`

	// Type of destination for the rule. Required if you provide a `destination`.
	//   * `CIDR_BLOCK`: If the rule's `destination` is an IP address range in CIDR notation.
	//   * `SERVICE_CIDR_BLOCK`: If the rule's `destination` is the `cidrBlock` value for a
	//     Service (the rule is for traffic destined for a
	//     particular `Service` through a service gateway).
	DestinationType RouteRuleDestinationTypeEnum `mandatory:"false" json:"destinationType,omitempty"`
}

func (m RouteRule) String() string {
	return common.PointerString(m)
}

// RouteRuleDestinationTypeEnum Enum with underlying type: string
type RouteRuleDestinationTypeEnum string

// Set of constants representing the allowable values for RouteRuleDestinationTypeEnum
const (
	RouteRuleDestinationTypeCidrBlock        RouteRuleDestinationTypeEnum = "CIDR_BLOCK"
	RouteRuleDestinationTypeServiceCidrBlock RouteRuleDestinationTypeEnum = "SERVICE_CIDR_BLOCK"
)

var mappingRouteRuleDestinationType = map[string]RouteRuleDestinationTypeEnum{
	"CIDR_BLOCK":         RouteRuleDestinationTypeCidrBlock,
	"SERVICE_CIDR_BLOCK": RouteRuleDestinationTypeServiceCidrBlock,
}

// GetRouteRuleDestinationTypeEnumValues Enumerates the set of values for RouteRuleDestinationTypeEnum
func GetRouteRuleDestinationTypeEnumValues() []RouteRuleDestinationTypeEnum {
	values := make([]RouteRuleDestinationTypeEnum, 0)
	for _, v := range mappingRouteRuleDestinationType {
		values = append(values, v)
	}
	return values
}
