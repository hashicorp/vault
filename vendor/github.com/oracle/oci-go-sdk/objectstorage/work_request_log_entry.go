// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object Storage and Archive Storage APIs for managing buckets, objects, and related resources.
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// WorkRequestLogEntry The representation of WorkRequestLogEntry
type WorkRequestLogEntry struct {

	// Human-readable log message.
	Message *string `mandatory:"false" json:"message"`

	// The date and time the log message was written, as described in
	// RFC 3339 (https://tools.ietf.org/rfc/rfc3339), section 14.29.
	Timestamp *common.SDKTime `mandatory:"false" json:"timestamp"`
}

func (m WorkRequestLogEntry) String() string {
	return common.PointerString(m)
}
