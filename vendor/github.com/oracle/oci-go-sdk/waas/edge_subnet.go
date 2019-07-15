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

// EdgeSubnet The details about an edge node subnet.
type EdgeSubnet struct {

	// An edge node subnet. This can include /24 or /8 addresses.
	Cidr *string `mandatory:"false" json:"cidr"`

	// The date and time the last change was made to the indicated edge node subnet, expressed in RFC 3339 timestamp format.
	TimeModified *common.SDKTime `mandatory:"false" json:"timeModified"`

	// The name of the region containing the indicated subnet.
	Region *string `mandatory:"false" json:"region"`
}

func (m EdgeSubnet) String() string {
	return common.PointerString(m)
}
