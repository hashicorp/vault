// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Health Checks API
//
// API for the Health Checks service. Use this API to manage endpoint probes and monitors.
// For more information, see
// Overview of the Health Checks Service (https://docs.cloud.oracle.com/iaas/Content/HealthChecks/Concepts/healthchecks.htm).
//

package healthchecks

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Routing The routing information for a vantage point.
type Routing struct {

	// The registry label for `asn`, usually the name of the organization that
	// owns the ASN. May be omitted or null.
	AsLabel *string `mandatory:"false" json:"asLabel"`

	// The Autonomous System Number (ASN) identifying the organization
	// responsible for routing packets to `prefix`.
	Asn *int `mandatory:"false" json:"asn"`

	// An IP prefix (CIDR syntax) that is less specific than
	// `address`, through which `address` is routed.
	Prefix *string `mandatory:"false" json:"prefix"`

	// An integer between 0 and 100 used to select between multiple
	// origin ASNs when routing to `prefix`. Most prefixes have
	// exactly one origin ASN, in which case `weight` will be 100.
	Weight *int `mandatory:"false" json:"weight"`
}

func (m Routing) String() string {
	return common.PointerString(m)
}
