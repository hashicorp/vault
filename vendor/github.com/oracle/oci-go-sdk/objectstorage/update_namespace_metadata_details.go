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

// UpdateNamespaceMetadataDetails An UpdateNamespaceMetadataDetails is used for update NamespaceMetadata. To be able to upate the NamespaceMetadata, a user
//  must have NAMESPACE_UPDATE permission.
type UpdateNamespaceMetadataDetails struct {

	// The update compartment id for an S3 client if this field is set.
	DefaultS3CompartmentId *string `mandatory:"false" json:"defaultS3CompartmentId"`

	// The update compartment id for a Swift client if this field is set.
	DefaultSwiftCompartmentId *string `mandatory:"false" json:"defaultSwiftCompartmentId"`
}

func (m UpdateNamespaceMetadataDetails) String() string {
	return common.PointerString(m)
}
