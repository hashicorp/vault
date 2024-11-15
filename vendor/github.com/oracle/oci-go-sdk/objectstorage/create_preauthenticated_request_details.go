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

// CreatePreauthenticatedRequestDetails The representation of CreatePreauthenticatedRequestDetails
type CreatePreauthenticatedRequestDetails struct {

	// A user-specified name for the pre-authenticated request. Names can be helpful in managing pre-authenticated requests.
	Name *string `mandatory:"true" json:"name"`

	// The operation that can be performed on this resource.
	AccessType CreatePreauthenticatedRequestDetailsAccessTypeEnum `mandatory:"true" json:"accessType"`

	// The expiration date for the pre-authenticated request as per RFC 3339 (https://tools.ietf.org/html/rfc3339).
	// After this date the pre-authenticated request will no longer be valid.
	TimeExpires *common.SDKTime `mandatory:"true" json:"timeExpires"`

	// The name of the object that is being granted access to by the pre-authenticated request. Avoid entering confidential
	// information. The object name can be null and if so, the pre-authenticated request grants access to the entire bucket.
	ObjectName *string `mandatory:"false" json:"objectName"`
}

func (m CreatePreauthenticatedRequestDetails) String() string {
	return common.PointerString(m)
}

// CreatePreauthenticatedRequestDetailsAccessTypeEnum Enum with underlying type: string
type CreatePreauthenticatedRequestDetailsAccessTypeEnum string

// Set of constants representing the allowable values for CreatePreauthenticatedRequestDetailsAccessTypeEnum
const (
	CreatePreauthenticatedRequestDetailsAccessTypeObjectread      CreatePreauthenticatedRequestDetailsAccessTypeEnum = "ObjectRead"
	CreatePreauthenticatedRequestDetailsAccessTypeObjectwrite     CreatePreauthenticatedRequestDetailsAccessTypeEnum = "ObjectWrite"
	CreatePreauthenticatedRequestDetailsAccessTypeObjectreadwrite CreatePreauthenticatedRequestDetailsAccessTypeEnum = "ObjectReadWrite"
	CreatePreauthenticatedRequestDetailsAccessTypeAnyobjectwrite  CreatePreauthenticatedRequestDetailsAccessTypeEnum = "AnyObjectWrite"
)

var mappingCreatePreauthenticatedRequestDetailsAccessType = map[string]CreatePreauthenticatedRequestDetailsAccessTypeEnum{
	"ObjectRead":      CreatePreauthenticatedRequestDetailsAccessTypeObjectread,
	"ObjectWrite":     CreatePreauthenticatedRequestDetailsAccessTypeObjectwrite,
	"ObjectReadWrite": CreatePreauthenticatedRequestDetailsAccessTypeObjectreadwrite,
	"AnyObjectWrite":  CreatePreauthenticatedRequestDetailsAccessTypeAnyobjectwrite,
}

// GetCreatePreauthenticatedRequestDetailsAccessTypeEnumValues Enumerates the set of values for CreatePreauthenticatedRequestDetailsAccessTypeEnum
func GetCreatePreauthenticatedRequestDetailsAccessTypeEnumValues() []CreatePreauthenticatedRequestDetailsAccessTypeEnum {
	values := make([]CreatePreauthenticatedRequestDetailsAccessTypeEnum, 0)
	for _, v := range mappingCreatePreauthenticatedRequestDetailsAccessType {
		values = append(values, v)
	}
	return values
}
