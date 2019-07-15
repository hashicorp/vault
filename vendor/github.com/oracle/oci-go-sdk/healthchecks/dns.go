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

// Dns The DNS resolution results.
type Dns struct {

	// Total DNS resolution duration, in milliseconds. Calculated using `domainLookupEnd`
	// minus `domainLookupStart`.
	DomainLookupDuration *float64 `mandatory:"false" json:"domainLookupDuration"`

	// The addresses returned by DNS resolution.
	Addresses []string `mandatory:"false" json:"addresses"`
}

func (m Dns) String() string {
	return common.PointerString(m)
}
