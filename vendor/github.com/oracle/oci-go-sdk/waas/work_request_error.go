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

// WorkRequestError An object returned in the event of a work request error.
type WorkRequestError struct {

	// A machine-usable code for the error that occurred.
	Code *string `mandatory:"false" json:"code"`

	// The error message.
	Message *string `mandatory:"false" json:"message"`

	// The date and time the work request error happened, expressed in RFC 3339 timestamp format.
	Timestamp *common.SDKTime `mandatory:"false" json:"timestamp"`
}

func (m WorkRequestError) String() string {
	return common.PointerString(m)
}
