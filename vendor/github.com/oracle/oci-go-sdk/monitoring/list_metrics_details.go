// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Monitoring API
//
// Use the Monitoring API to manage metric queries and alarms for assessing the health, capacity, and performance of your cloud resources.
// For information about monitoring, see Monitoring Overview (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm).
//

package monitoring

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ListMetricsDetails The request details for retrieving metric definitions. Specify optional properties to filter the returned results.
// Use an asterisk ("\*") as a wildcard character, placed anywhere in the string.
// For example, to search for all metrics with names that begin with "disk", specify "name" as "disk\*".
// If no properties are specified, then all metric definitions within the request scope are returned.
type ListMetricsDetails struct {

	// The metric name to use when searching for metric definitions.
	// Example: `CpuUtilization`
	Name *string `mandatory:"false" json:"name"`

	// The source service or application to use when searching for metric definitions.
	// Example: `oci_computeagent`
	Namespace *string `mandatory:"false" json:"namespace"`

	// Qualifiers that you want to use when searching for metric definitions.
	// Available dimensions vary by metric namespace. Each dimension takes the form of a key-value pair.
	// Example: { "resourceId": "<var>&lt;instance_OCID&gt;</var>" }
	DimensionFilters map[string]string `mandatory:"false" json:"dimensionFilters"`

	// Group metrics by these fields in the response. For example, to list all metric namespaces available
	// in a compartment, groupBy the "namespace" field.
	// Example - group by namespace and resource:
	// `[ "namespace", "resourceId" ]`
	GroupBy []string `mandatory:"false" json:"groupBy"`

	// The field to use when sorting returned metric definitions. Only one sorting level is provided.
	// Example: `NAMESPACE`
	SortBy ListMetricsDetailsSortByEnum `mandatory:"false" json:"sortBy,omitempty"`

	// The sort order to use when sorting returned metric definitions. Ascending (ASC) or
	// descending (DESC).
	// Example: `ASC`
	SortOrder ListMetricsDetailsSortOrderEnum `mandatory:"false" json:"sortOrder,omitempty"`
}

func (m ListMetricsDetails) String() string {
	return common.PointerString(m)
}

// ListMetricsDetailsSortByEnum Enum with underlying type: string
type ListMetricsDetailsSortByEnum string

// Set of constants representing the allowable values for ListMetricsDetailsSortByEnum
const (
	ListMetricsDetailsSortByNamespace ListMetricsDetailsSortByEnum = "NAMESPACE"
	ListMetricsDetailsSortByName      ListMetricsDetailsSortByEnum = "NAME"
)

var mappingListMetricsDetailsSortBy = map[string]ListMetricsDetailsSortByEnum{
	"NAMESPACE": ListMetricsDetailsSortByNamespace,
	"NAME":      ListMetricsDetailsSortByName,
}

// GetListMetricsDetailsSortByEnumValues Enumerates the set of values for ListMetricsDetailsSortByEnum
func GetListMetricsDetailsSortByEnumValues() []ListMetricsDetailsSortByEnum {
	values := make([]ListMetricsDetailsSortByEnum, 0)
	for _, v := range mappingListMetricsDetailsSortBy {
		values = append(values, v)
	}
	return values
}

// ListMetricsDetailsSortOrderEnum Enum with underlying type: string
type ListMetricsDetailsSortOrderEnum string

// Set of constants representing the allowable values for ListMetricsDetailsSortOrderEnum
const (
	ListMetricsDetailsSortOrderAsc  ListMetricsDetailsSortOrderEnum = "ASC"
	ListMetricsDetailsSortOrderDesc ListMetricsDetailsSortOrderEnum = "DESC"
)

var mappingListMetricsDetailsSortOrder = map[string]ListMetricsDetailsSortOrderEnum{
	"ASC":  ListMetricsDetailsSortOrderAsc,
	"DESC": ListMetricsDetailsSortOrderDesc,
}

// GetListMetricsDetailsSortOrderEnumValues Enumerates the set of values for ListMetricsDetailsSortOrderEnum
func GetListMetricsDetailsSortOrderEnumValues() []ListMetricsDetailsSortOrderEnum {
	values := make([]ListMetricsDetailsSortOrderEnum, 0)
	for _, v := range mappingListMetricsDetailsSortOrder {
		values = append(values, v)
	}
	return values
}
