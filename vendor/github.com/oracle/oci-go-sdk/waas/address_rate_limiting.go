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

// AddressRateLimiting The IP rate limiting configuration. Defines the amount of allowed requests from a unique IP address and the resulting block response code when that threshold is exceeded.
type AddressRateLimiting struct {

	// Enables or disables the address rate limiting Web Application Firewall feature.
	IsEnabled *bool `mandatory:"true" json:"isEnabled"`

	// The number of allowed requests per second from one IP address. If unspecified, defaults to `1`.
	AllowedRatePerAddress *int `mandatory:"false" json:"allowedRatePerAddress"`

	// The maximum number of requests allowed to be queued before subsequent requests are dropped. If unspecified, defaults to `10`.
	MaxDelayedCountPerAddress *int `mandatory:"false" json:"maxDelayedCountPerAddress"`

	// The response status code returned when a request is blocked. If unspecified, defaults to `503`.
	BlockResponseCode *int `mandatory:"false" json:"blockResponseCode"`
}

func (m AddressRateLimiting) String() string {
	return common.PointerString(m)
}
