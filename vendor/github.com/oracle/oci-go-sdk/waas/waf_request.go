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

// WafRequest A time series of request counts handled by the Web Application Firewall, including blocked requests.
type WafRequest struct {

	// The date and time the traffic was observed, rounded down to the start of a range, and expressed in RFC 3339 timestamp format.
	TimeObserved *common.SDKTime `mandatory:"false" json:"timeObserved"`

	// The number of seconds this data covers.
	TimeRangeInSeconds *int `mandatory:"false" json:"timeRangeInSeconds"`

	// The total number of requests received in this time period.
	Count *int `mandatory:"false" json:"count"`
}

func (m WafRequest) String() string {
	return common.PointerString(m)
}
