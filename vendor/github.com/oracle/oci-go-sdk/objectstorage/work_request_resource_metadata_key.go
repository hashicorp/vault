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

// WorkRequestResourceMetadataKeyEnum Enum with underlying type: string
type WorkRequestResourceMetadataKeyEnum string

// Set of constants representing the allowable values for WorkRequestResourceMetadataKeyEnum
const (
	WorkRequestResourceMetadataKeyRegion    WorkRequestResourceMetadataKeyEnum = "REGION"
	WorkRequestResourceMetadataKeyNamespace WorkRequestResourceMetadataKeyEnum = "NAMESPACE"
	WorkRequestResourceMetadataKeyBucket    WorkRequestResourceMetadataKeyEnum = "BUCKET"
	WorkRequestResourceMetadataKeyObject    WorkRequestResourceMetadataKeyEnum = "OBJECT"
)

var mappingWorkRequestResourceMetadataKey = map[string]WorkRequestResourceMetadataKeyEnum{
	"REGION":    WorkRequestResourceMetadataKeyRegion,
	"NAMESPACE": WorkRequestResourceMetadataKeyNamespace,
	"BUCKET":    WorkRequestResourceMetadataKeyBucket,
	"OBJECT":    WorkRequestResourceMetadataKeyObject,
}

// GetWorkRequestResourceMetadataKeyEnumValues Enumerates the set of values for WorkRequestResourceMetadataKeyEnum
func GetWorkRequestResourceMetadataKeyEnumValues() []WorkRequestResourceMetadataKeyEnum {
	values := make([]WorkRequestResourceMetadataKeyEnum, 0)
	for _, v := range mappingWorkRequestResourceMetadataKey {
		values = append(values, v)
	}
	return values
}
