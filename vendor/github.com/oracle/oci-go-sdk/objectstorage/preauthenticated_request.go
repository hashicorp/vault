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

// PreauthenticatedRequest Pre-authenticated requests provide a way to let users access a bucket or an object without having their own credentials.
// When you create a pre-authenticated request, a unique URL is generated. Users in your organization, partners, or third
// parties can use this URL to access the targets identified in the pre-authenticated request. See Managing Access to Buckets and Objects (https://docs.us-phoenix-1.oraclecloud.com/Content/Object/Tasks/managingaccess.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized, talk to an administrator.
// If you're an administrator who needs to write policies to give users access, see Getting Started with Policies (https://docs.us-phoenix-1.oraclecloud.com/Content/Identity/Concepts/policygetstarted.htm).
type PreauthenticatedRequest struct {

	// The unique identifier to use when directly addressing the pre-authenticated request.
	Id *string `mandatory:"true" json:"id"`

	// The user-provided name of the pre-authenticated request.
	Name *string `mandatory:"true" json:"name"`

	// The URI to embed in the URL when using the pre-authenticated request.
	AccessUri *string `mandatory:"true" json:"accessUri"`

	// The operation that can be performed on this resource.
	AccessType PreauthenticatedRequestAccessTypeEnum `mandatory:"true" json:"accessType"`

	// The expiration date for the pre-authenticated request as per RFC 3339 (https://tools.ietf.org/rfc/rfc3339). After this date the pre-authenticated request will no longer be valid.
	TimeExpires *common.SDKTime `mandatory:"true" json:"timeExpires"`

	// The date when the pre-authenticated request was created as per specification
	// RFC 3339 (https://tools.ietf.org/rfc/rfc3339).
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The name of the object that is being granted access to by the pre-authenticated request. This can be null and
	// if so, the pre-authenticated request grants access to the entire bucket. Avoid entering confidential information.
	// Example: test/object1.log
	ObjectName *string `mandatory:"false" json:"objectName"`
}

func (m PreauthenticatedRequest) String() string {
	return common.PointerString(m)
}

// PreauthenticatedRequestAccessTypeEnum Enum with underlying type: string
type PreauthenticatedRequestAccessTypeEnum string

// Set of constants representing the allowable values for PreauthenticatedRequestAccessType
const (
	PreauthenticatedRequestAccessTypeObjectread      PreauthenticatedRequestAccessTypeEnum = "ObjectRead"
	PreauthenticatedRequestAccessTypeObjectwrite     PreauthenticatedRequestAccessTypeEnum = "ObjectWrite"
	PreauthenticatedRequestAccessTypeObjectreadwrite PreauthenticatedRequestAccessTypeEnum = "ObjectReadWrite"
	PreauthenticatedRequestAccessTypeAnyobjectwrite  PreauthenticatedRequestAccessTypeEnum = "AnyObjectWrite"
)

var mappingPreauthenticatedRequestAccessType = map[string]PreauthenticatedRequestAccessTypeEnum{
	"ObjectRead":      PreauthenticatedRequestAccessTypeObjectread,
	"ObjectWrite":     PreauthenticatedRequestAccessTypeObjectwrite,
	"ObjectReadWrite": PreauthenticatedRequestAccessTypeObjectreadwrite,
	"AnyObjectWrite":  PreauthenticatedRequestAccessTypeAnyobjectwrite,
}

// GetPreauthenticatedRequestAccessTypeEnumValues Enumerates the set of values for PreauthenticatedRequestAccessType
func GetPreauthenticatedRequestAccessTypeEnumValues() []PreauthenticatedRequestAccessTypeEnum {
	values := make([]PreauthenticatedRequestAccessTypeEnum, 0)
	for _, v := range mappingPreauthenticatedRequestAccessType {
		values = append(values, v)
	}
	return values
}
