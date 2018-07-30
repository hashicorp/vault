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

// CreatePreauthenticatedRequestDetails The representation of CreatePreauthenticatedRequestDetails
type CreatePreauthenticatedRequestDetails struct {

	// A user-specified name for the pre-authenticated request. Helpful for management purposes.
	Name *string `mandatory:"true" json:"name"`

	// The operation that can be performed on this resource.
	AccessType CreatePreauthenticatedRequestDetailsAccessTypeEnum `mandatory:"true" json:"accessType"`

	// The expiration date for the pre-authenticated request as per RFC 3339 (https://tools.ietf.org/rfc/rfc3339). After this date the pre-authenticated request will no longer be valid.
	TimeExpires *common.SDKTime `mandatory:"true" json:"timeExpires"`

	// The name of object that is being granted access to by the pre-authenticated request. This can be null and if it is, the pre-authenticated request grants access to the entire bucket.
	ObjectName *string `mandatory:"false" json:"objectName"`
}

func (m CreatePreauthenticatedRequestDetails) String() string {
	return common.PointerString(m)
}

// CreatePreauthenticatedRequestDetailsAccessTypeEnum Enum with underlying type: string
type CreatePreauthenticatedRequestDetailsAccessTypeEnum string

// Set of constants representing the allowable values for CreatePreauthenticatedRequestDetailsAccessType
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

// GetCreatePreauthenticatedRequestDetailsAccessTypeEnumValues Enumerates the set of values for CreatePreauthenticatedRequestDetailsAccessType
func GetCreatePreauthenticatedRequestDetailsAccessTypeEnumValues() []CreatePreauthenticatedRequestDetailsAccessTypeEnum {
	values := make([]CreatePreauthenticatedRequestDetailsAccessTypeEnum, 0)
	for _, v := range mappingCreatePreauthenticatedRequestDetailsAccessType {
		values = append(values, v)
	}
	return values
}
