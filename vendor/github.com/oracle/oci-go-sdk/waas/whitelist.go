// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Web Application Acceleration and Security Services API
//
// OCI Web Application Acceleration and Security Services
//

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Whitelist An array of IP addresses that bypass the Web Application Firewall. Supports both single IP addresses or subnet masks (CIDR notation).
type Whitelist struct {

	// The unique name of the whitelist.
	Name *string `mandatory:"true" json:"name"`

	// A set of IP addresses or CIDR notations to include in the whitelist.
	Addresses []string `mandatory:"true" json:"addresses"`
}

func (m Whitelist) String() string {
	return common.PointerString(m)
}
