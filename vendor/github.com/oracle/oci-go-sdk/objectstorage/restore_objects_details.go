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

// RestoreObjectsDetails The representation of RestoreObjectsDetails
type RestoreObjectsDetails struct {

	// A object which was in an archived state and need to be restored.
	ObjectName *string `mandatory:"true" json:"objectName"`

	// The number of hours for which this object will be restored.
	// By default object will be restored for 24 hours.It can be configured using hours parameter.
	Hours *int `mandatory:"false" json:"hours"`
}

func (m RestoreObjectsDetails) String() string {
	return common.PointerString(m)
}
