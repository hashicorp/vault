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

// PreauthenticatedRequestSummary Get summary information about pre-authenticated requests.
type PreauthenticatedRequestSummary struct {

	// The unique identifier to use when directly addressing the pre-authenticated request.
	Id *string `mandatory:"true" json:"id"`

	// The user-provided name of the pre-authenticated request.
	Name *string `mandatory:"true" json:"name"`

	// The operation that can be performed on this resource.
	AccessType PreauthenticatedRequestSummaryAccessTypeEnum `mandatory:"true" json:"accessType"`

	// The expiration date for the pre-authenticated request as per RFC 3339 (https://tools.ietf.org/html/rfc3339). After this date the pre-authenticated request will no longer be valid.
	TimeExpires *common.SDKTime `mandatory:"true" json:"timeExpires"`

	// The date when the pre-authenticated request was created as per RFC 3339 (https://tools.ietf.org/html/rfc3339).
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The name of object that is being granted access to by the pre-authenticated request. This can be null and if it is,
	// the pre-authenticated request grants access to the entire bucket.
	ObjectName *string `mandatory:"false" json:"objectName"`
}

func (m PreauthenticatedRequestSummary) String() string {
	return common.PointerString(m)
}

// PreauthenticatedRequestSummaryAccessTypeEnum Enum with underlying type: string
type PreauthenticatedRequestSummaryAccessTypeEnum string

// Set of constants representing the allowable values for PreauthenticatedRequestSummaryAccessTypeEnum
const (
	PreauthenticatedRequestSummaryAccessTypeObjectread      PreauthenticatedRequestSummaryAccessTypeEnum = "ObjectRead"
	PreauthenticatedRequestSummaryAccessTypeObjectwrite     PreauthenticatedRequestSummaryAccessTypeEnum = "ObjectWrite"
	PreauthenticatedRequestSummaryAccessTypeObjectreadwrite PreauthenticatedRequestSummaryAccessTypeEnum = "ObjectReadWrite"
	PreauthenticatedRequestSummaryAccessTypeAnyobjectwrite  PreauthenticatedRequestSummaryAccessTypeEnum = "AnyObjectWrite"
)

var mappingPreauthenticatedRequestSummaryAccessType = map[string]PreauthenticatedRequestSummaryAccessTypeEnum{
	"ObjectRead":      PreauthenticatedRequestSummaryAccessTypeObjectread,
	"ObjectWrite":     PreauthenticatedRequestSummaryAccessTypeObjectwrite,
	"ObjectReadWrite": PreauthenticatedRequestSummaryAccessTypeObjectreadwrite,
	"AnyObjectWrite":  PreauthenticatedRequestSummaryAccessTypeAnyobjectwrite,
}

// GetPreauthenticatedRequestSummaryAccessTypeEnumValues Enumerates the set of values for PreauthenticatedRequestSummaryAccessTypeEnum
func GetPreauthenticatedRequestSummaryAccessTypeEnumValues() []PreauthenticatedRequestSummaryAccessTypeEnum {
	values := make([]PreauthenticatedRequestSummaryAccessTypeEnum, 0)
	for _, v := range mappingPreauthenticatedRequestSummaryAccessType {
		values = append(values, v)
	}
	return values
}
