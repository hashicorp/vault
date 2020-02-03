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

// UpdateNamespaceMetadataDetails UpdateNamespaceMetadataDetails is used to update the NamespaceMetadata. To update NamespaceMetadata, a user
// must have OBJECTSTORAGE_NAMESPACE_UPDATE permission.
type UpdateNamespaceMetadataDetails struct {

	// The updated compartment id for use by an S3 client, if this field is set.
	DefaultS3CompartmentId *string `mandatory:"false" json:"defaultS3CompartmentId"`

	// The updated compartment id for use by a Swift client, if this field is set.
	DefaultSwiftCompartmentId *string `mandatory:"false" json:"defaultSwiftCompartmentId"`
}

func (m UpdateNamespaceMetadataDetails) String() string {
	return common.PointerString(m)
}
