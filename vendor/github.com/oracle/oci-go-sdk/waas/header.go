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

// Header An HTTP header name and value. You can configure your origin server to only allow requests that contain the custom header values that you specify.
type Header struct {

	// The name of the header.
	Name *string `mandatory:"true" json:"name"`

	// The value of the header.
	Value *string `mandatory:"true" json:"value"`
}

func (m Header) String() string {
	return common.PointerString(m)
}
