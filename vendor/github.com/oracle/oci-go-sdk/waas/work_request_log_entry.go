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

// WorkRequestLogEntry A log message for a work request.
type WorkRequestLogEntry struct {

	// The log message.
	Message *string `mandatory:"false" json:"message"`

	// The date and time the work request log event happend, expressed in RFC 3339 timestamp format.
	Timestamp *common.SDKTime `mandatory:"false" json:"timestamp"`
}

func (m WorkRequestLogEntry) String() string {
	return common.PointerString(m)
}
