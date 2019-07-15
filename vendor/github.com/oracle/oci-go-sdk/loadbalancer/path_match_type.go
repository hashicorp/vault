// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Load Balancing API
//
// API for the Load Balancing service. Use this API to manage load balancers, backend sets, and related items. For more
// information, see Overview of Load Balancing (https://docs.cloud.oracle.com/iaas/Content/Balance/Concepts/balanceoverview.htm).
//

package loadbalancer

import (
	"github.com/oracle/oci-go-sdk/common"
)

// PathMatchType The type of matching to apply to incoming URIs.
type PathMatchType struct {

	// Specifies how the load balancing service compares a PathRoute
	// object's `path` string against the incoming URI.
	// *  **EXACT_MATCH** - Looks for a `path` string that exactly matches the incoming URI path.
	// *  **FORCE_LONGEST_PREFIX_MATCH** - Looks for the `path` string with the best, longest match of the beginning
	//    portion of the incoming URI path.
	// *  **PREFIX_MATCH** - Looks for a `path` string that matches the beginning portion of the incoming URI path.
	// *  **SUFFIX_MATCH** - Looks for a `path` string that matches the ending portion of the incoming URI path.
	// For a full description of how the system handles `matchType` in a path route set containing multiple rules, see
	// Managing Request Routing (https://docs.cloud.oracle.com/Content/Balance/Tasks/managingrequest.htm).
	MatchType PathMatchTypeMatchTypeEnum `mandatory:"true" json:"matchType"`
}

func (m PathMatchType) String() string {
	return common.PointerString(m)
}

// PathMatchTypeMatchTypeEnum Enum with underlying type: string
type PathMatchTypeMatchTypeEnum string

// Set of constants representing the allowable values for PathMatchTypeMatchTypeEnum
const (
	PathMatchTypeMatchTypeExactMatch              PathMatchTypeMatchTypeEnum = "EXACT_MATCH"
	PathMatchTypeMatchTypeForceLongestPrefixMatch PathMatchTypeMatchTypeEnum = "FORCE_LONGEST_PREFIX_MATCH"
	PathMatchTypeMatchTypePrefixMatch             PathMatchTypeMatchTypeEnum = "PREFIX_MATCH"
	PathMatchTypeMatchTypeSuffixMatch             PathMatchTypeMatchTypeEnum = "SUFFIX_MATCH"
)

var mappingPathMatchTypeMatchType = map[string]PathMatchTypeMatchTypeEnum{
	"EXACT_MATCH":                PathMatchTypeMatchTypeExactMatch,
	"FORCE_LONGEST_PREFIX_MATCH": PathMatchTypeMatchTypeForceLongestPrefixMatch,
	"PREFIX_MATCH":               PathMatchTypeMatchTypePrefixMatch,
	"SUFFIX_MATCH":               PathMatchTypeMatchTypeSuffixMatch,
}

// GetPathMatchTypeMatchTypeEnumValues Enumerates the set of values for PathMatchTypeMatchTypeEnum
func GetPathMatchTypeMatchTypeEnumValues() []PathMatchTypeMatchTypeEnum {
	values := make([]PathMatchTypeMatchTypeEnum, 0)
	for _, v := range mappingPathMatchTypeMatchType {
		values = append(values, v)
	}
	return values
}
