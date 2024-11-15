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

// RenameObjectDetails To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type RenameObjectDetails struct {

	// The name of the source object to be renamed.
	SourceName *string `mandatory:"true" json:"sourceName"`

	// The new name of the source object.
	NewName *string `mandatory:"true" json:"newName"`

	// The if-match entity tag (ETag) of the source object.
	SrcObjIfMatchETag *string `mandatory:"false" json:"srcObjIfMatchETag"`

	// The if-match entity tag (ETag) of the new object.
	NewObjIfMatchETag *string `mandatory:"false" json:"newObjIfMatchETag"`

	// The if-none-match entity tag (ETag) of the new object.
	NewObjIfNoneMatchETag *string `mandatory:"false" json:"newObjIfNoneMatchETag"`
}

func (m RenameObjectDetails) String() string {
	return common.PointerString(m)
}
