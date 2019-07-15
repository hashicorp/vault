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

// WafMeterDatum The representation of WafMeterDatum
type WafMeterDatum struct {

	// The date and time the traffic was observed, rounded down to the start of a range, and expressed in RFC 3339 timestamp format.
	TimeObserved *common.SDKTime `mandatory:"false" json:"timeObserved"`

	// The number of seconds this data covers.
	TimeRangeInSeconds *int `mandatory:"false" json:"timeRangeInSeconds"`

	// The tenancy OCID of the data.
	TenancyId *string `mandatory:"false" json:"tenancyId"`

	// The compartment OCID of the data.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The policy OCID of the data.
	WaasPolicyId *string `mandatory:"false" json:"waasPolicyId"`

	// True if origin (endpoint) is an OCI resource. False if external.
	IsOciOrigin *bool `mandatory:"false" json:"isOciOrigin"`

	// True if bot manager is enabled.
	IsBotEnabled *bool `mandatory:"false" json:"isBotEnabled"`

	// The number of incoming requests.
	RequestCount *int `mandatory:"false" json:"requestCount"`

	// Traffic in bytes.
	TrafficInBytes *int `mandatory:"false" json:"trafficInBytes"`

	// The tag slug for the specified `waasPolicyId`.
	TagSlug []byte `mandatory:"false" json:"tagSlug"`
}

func (m WafMeterDatum) String() string {
	return common.PointerString(m)
}
