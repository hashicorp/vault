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

// HealthChecksVantagePointSummary Information about a vantage point.
type HealthChecksVantagePointSummary struct {

	// The display name for the vantage point. Display names are determined by
	// the best information available and may change over time.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The organization on whose infrastructure this vantage point resides.
	// Provider names are not unique, as Oracle Cloud Infrastructure maintains
	// many vantage points in each major provider.
	ProviderName *string `mandatory:"false" json:"providerName"`

	// The unique, permanent name for the vantage point.
	Name *string `mandatory:"false" json:"name"`

	Geo *Geolocation `mandatory:"false" json:"geo"`

	// An array of objects that describe how traffic to this vantage point is
	// routed, including which prefixes and ASNs connect it to the internet.
	// The addresses are sorted from the most-specific to least-specific
	// prefix (the smallest network to largest network). When a prefix has
	// multiple origin ASNs (MOAS routing), they are sorted by weight
	// (highest to lowest). Weight is determined by the total percentage of
	// peers observing the prefix originating from an ASN. Only present if
	// `fields` includes `routing`. The field will be null if the address's
	// routing information is unknown.
	Routing []Routing `mandatory:"false" json:"routing"`
}

func (m HealthChecksVantagePointSummary) String() string {
	return common.PointerString(m)
}
