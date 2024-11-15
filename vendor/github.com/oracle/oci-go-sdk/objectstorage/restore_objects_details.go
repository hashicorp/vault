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

// RestoreObjectsDetails The representation of RestoreObjectsDetails
type RestoreObjectsDetails struct {

	// An object that is in an archive storage tier and needs to be restored.
	ObjectName *string `mandatory:"true" json:"objectName"`

	// The number of hours for which this object will be restored.
	// By default objects will be restored for 24 hours. You can instead configure the duration using the hours parameter.
	Hours *int `mandatory:"false" json:"hours"`

	// The versionId of the object to restore. Current object version is used by default.
	VersionId *string `mandatory:"false" json:"versionId"`
}

func (m RestoreObjectsDetails) String() string {
	return common.PointerString(m)
}
