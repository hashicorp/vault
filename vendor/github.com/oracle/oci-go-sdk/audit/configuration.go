// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Audit API
//
// API for the Audit Service. You can use this API for queries, but not bulk-export operations.
//

package audit

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Configuration The representation of Configuration
type Configuration struct {

	// The retention period days
	RetentionPeriodDays *int `mandatory:"false" json:"retentionPeriodDays"`
}

func (m Configuration) String() string {
	return common.PointerString(m)
}
