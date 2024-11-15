// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object Storage and Archive Storage APIs for managing buckets, objects, and related resources.
// For more information, see Overview of Object Storage (https://docs.cloud.oracle.com/Content/Object/Concepts/objectstorageoverview.htm) and
// Overview of Archive Storage (https://docs.cloud.oracle.com/Content/Archive/Concepts/archivestorageoverview.htm).
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ListObjects To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type ListObjects struct {

	// An array of object summaries.
	Objects []ObjectSummary `mandatory:"true" json:"objects"`

	// Prefixes that are common to the results returned by the request if the request specified a delimiter.
	Prefixes []string `mandatory:"false" json:"prefixes"`

	// The name of the object to use in the 'start' parameter to obtain the next page of
	// a truncated ListObjects response. Avoid entering confidential information.
	// Example: test/object1.log
	NextStartWith *string `mandatory:"false" json:"nextStartWith"`
}

func (m ListObjects) String() string {
	return common.PointerString(m)
}
