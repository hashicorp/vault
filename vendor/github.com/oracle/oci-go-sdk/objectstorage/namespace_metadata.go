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

// NamespaceMetadata NamespaceMetadata maps a namespace string to defaultS3CompartmentId and defaultSwiftCompartmentId values.
type NamespaceMetadata struct {

	// The Object Storage namespace to which the metadata belongs.
	Namespace *string `mandatory:"true" json:"namespace"`

	// If the field is set, specifies the default compartment assignment for the Amazon S3 Compatibility API.
	DefaultS3CompartmentId *string `mandatory:"true" json:"defaultS3CompartmentId"`

	// If the field is set, specifies the default compartment assignment for the Swift API.
	DefaultSwiftCompartmentId *string `mandatory:"true" json:"defaultSwiftCompartmentId"`
}

func (m NamespaceMetadata) String() string {
	return common.PointerString(m)
}
