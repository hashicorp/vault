// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Search Service API
//
// Search for resources in your cloud network.
//

package resourcesearch

import (
	"github.com/oracle/oci-go-sdk/common"
)

// QueryableFieldDescription An individual field that can be used as part of a query filter.
type QueryableFieldDescription struct {

	// The type of the field, which dictates what semantics and query constraints you can use when searching or querying.
	FieldType QueryableFieldDescriptionFieldTypeEnum `mandatory:"true" json:"fieldType"`

	// The name of the field to use when constructing the query. Field names are present for all types except `OBJECT`.
	FieldName *string `mandatory:"true" json:"fieldName"`

	// Indicates this field is actually an array of the specified field type.
	IsArray *bool `mandatory:"false" json:"isArray"`

	// If the field type is `OBJECT`, then this property will provide all the individual properties on the object that can
	// be queried.
	ObjectProperties []QueryableFieldDescription `mandatory:"false" json:"objectProperties"`
}

func (m QueryableFieldDescription) String() string {
	return common.PointerString(m)
}

// QueryableFieldDescriptionFieldTypeEnum Enum with underlying type: string
type QueryableFieldDescriptionFieldTypeEnum string

// Set of constants representing the allowable values for QueryableFieldDescriptionFieldTypeEnum
const (
	QueryableFieldDescriptionFieldTypeIdentifier QueryableFieldDescriptionFieldTypeEnum = "IDENTIFIER"
	QueryableFieldDescriptionFieldTypeString     QueryableFieldDescriptionFieldTypeEnum = "STRING"
	QueryableFieldDescriptionFieldTypeInteger    QueryableFieldDescriptionFieldTypeEnum = "INTEGER"
	QueryableFieldDescriptionFieldTypeRational   QueryableFieldDescriptionFieldTypeEnum = "RATIONAL"
	QueryableFieldDescriptionFieldTypeBoolean    QueryableFieldDescriptionFieldTypeEnum = "BOOLEAN"
	QueryableFieldDescriptionFieldTypeDatetime   QueryableFieldDescriptionFieldTypeEnum = "DATETIME"
	QueryableFieldDescriptionFieldTypeIp         QueryableFieldDescriptionFieldTypeEnum = "IP"
	QueryableFieldDescriptionFieldTypeObject     QueryableFieldDescriptionFieldTypeEnum = "OBJECT"
)

var mappingQueryableFieldDescriptionFieldType = map[string]QueryableFieldDescriptionFieldTypeEnum{
	"IDENTIFIER": QueryableFieldDescriptionFieldTypeIdentifier,
	"STRING":     QueryableFieldDescriptionFieldTypeString,
	"INTEGER":    QueryableFieldDescriptionFieldTypeInteger,
	"RATIONAL":   QueryableFieldDescriptionFieldTypeRational,
	"BOOLEAN":    QueryableFieldDescriptionFieldTypeBoolean,
	"DATETIME":   QueryableFieldDescriptionFieldTypeDatetime,
	"IP":         QueryableFieldDescriptionFieldTypeIp,
	"OBJECT":     QueryableFieldDescriptionFieldTypeObject,
}

// GetQueryableFieldDescriptionFieldTypeEnumValues Enumerates the set of values for QueryableFieldDescriptionFieldTypeEnum
func GetQueryableFieldDescriptionFieldTypeEnumValues() []QueryableFieldDescriptionFieldTypeEnum {
	values := make([]QueryableFieldDescriptionFieldTypeEnum, 0)
	for _, v := range mappingQueryableFieldDescriptionFieldType {
		values = append(values, v)
	}
	return values
}
