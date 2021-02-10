// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Key Management Service API
//
// API for managing and performing operations with keys and vaults.
//

package keymanagement

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ScheduleKeyDeletionDetails Details for scheduling key deletion.
type ScheduleKeyDeletionDetails struct {

	// An optional property to indicate when to delete the vault, expressed in
	// RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format. The specified
	// time must be between 7 and 30 days from when the request is received.
	// If this property is missing, it will be set to 30 days from the time of the request
	// by default.
	TimeOfDeletion *common.SDKTime `mandatory:"false" json:"timeOfDeletion"`
}

func (m ScheduleKeyDeletionDetails) String() string {
	return common.PointerString(m)
}
