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
