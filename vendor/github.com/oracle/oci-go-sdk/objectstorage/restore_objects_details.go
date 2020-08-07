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

// RestoreObjectsDetails The representation of RestoreObjectsDetails
type RestoreObjectsDetails struct {

	// An object that is in an archive storage tier and needs to be restored.
	ObjectName *string `mandatory:"true" json:"objectName"`

	// The number of hours for which this object will be restored.
	// By default objects will be restored for 24 hours. You can instead configure the duration using the hours parameter.
	Hours *int `mandatory:"false" json:"hours"`
}

func (m RestoreObjectsDetails) String() string {
	return common.PointerString(m)
}
