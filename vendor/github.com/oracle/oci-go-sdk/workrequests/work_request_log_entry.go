// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Work Requests API
//
// A description of the work requests API
//

package workrequests

import (
	"github.com/oracle/oci-go-sdk/common"
)

// WorkRequestLogEntry A log message from the execution of a work request.
type WorkRequestLogEntry struct {

	// A human-readable log message.
	Message *string `mandatory:"true" json:"message"`

	// The time the log message was written.
	Timestamp *common.SDKTime `mandatory:"true" json:"timestamp"`
}

func (m WorkRequestLogEntry) String() string {
	return common.PointerString(m)
}
