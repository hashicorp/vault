// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object and Archive Storage APIs for managing buckets and objects.
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ListObjects To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.us-phoenix-1.oraclecloud.com/Content/Identity/Concepts/policygetstarted.htm).
type ListObjects struct {

	// An array of object summaries.
	Objects []ObjectSummary `mandatory:"true" json:"objects"`

	// Prefixes that are common to the results returned by the request if the request specified a delimiter.
	Prefixes []string `mandatory:"false" json:"prefixes"`

	// The name of the object to use in the 'startWith' parameter to obtain the next page of
	// a truncated ListObjects response. Avoid entering confidential information.
	// Example: test/object1.log
	NextStartWith *string `mandatory:"false" json:"nextStartWith"`
}

func (m ListObjects) String() string {
	return common.PointerString(m)
}
