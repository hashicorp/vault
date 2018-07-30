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

// NamespaceMetadata A NamespaceMetadta is a map for storing namespace and defaultS3CompartmentId, defaultSwiftCompartmentId.
type NamespaceMetadata struct {

	// The namespace to which the metadata belongs.
	Namespace *string `mandatory:"true" json:"namespace"`

	// The default compartment ID for an S3 client.
	DefaultS3CompartmentId *string `mandatory:"true" json:"defaultS3CompartmentId"`

	// The default compartment ID for a Swift client.
	DefaultSwiftCompartmentId *string `mandatory:"true" json:"defaultSwiftCompartmentId"`
}

func (m NamespaceMetadata) String() string {
	return common.PointerString(m)
}
