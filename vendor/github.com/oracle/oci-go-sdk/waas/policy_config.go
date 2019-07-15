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

// PolicyConfig The configuration details for the WAAS policy.
type PolicyConfig struct {

	// The OCID of the SSL certificate to use if HTTPS is supported.
	CertificateId *string `mandatory:"false" json:"certificateId"`

	// Enable or disable HTTPS support. If true, a `certificateId` is required. If unspecified, defaults to `false`.
	IsHttpsEnabled *bool `mandatory:"false" json:"isHttpsEnabled"`

	// Force HTTP to HTTPS redirection. If unspecified, defaults to `false`.
	IsHttpsForced *bool `mandatory:"false" json:"isHttpsForced"`
}

func (m PolicyConfig) String() string {
	return common.PointerString(m)
}
