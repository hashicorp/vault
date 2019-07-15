// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"github.com/oracle/oci-go-sdk/common"
)

// DbVersionSummary The Oracle Database software version.
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized, talk to an administrator. If you're an administrator who needs to write policies to give users access, see Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type DbVersionSummary struct {

	// A valid Oracle Database version.
	Version *string `mandatory:"true" json:"version"`

	// True if this version of the Oracle Database software is the latest version for a release.
	IsLatestForMajorVersion *bool `mandatory:"false" json:"isLatestForMajorVersion"`

	// True if this version of the Oracle Database software supports pluggable databases.
	SupportsPdb *bool `mandatory:"false" json:"supportsPdb"`
}

func (m DbVersionSummary) String() string {
	return common.PointerString(m)
}
