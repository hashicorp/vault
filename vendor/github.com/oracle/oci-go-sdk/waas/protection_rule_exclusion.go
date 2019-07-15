// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Web Application Acceleration and Security Services API
//
// OCI Web Application Acceleration and Security Services
//

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ProtectionRuleExclusion Allows specified types of requests to bypass the protection rule. If the requests matches any of the criteria in the `exclusions` field, the protection rule will not be executed. Rules can have more than one exclusion and exclusions are applied to requests disjunctively.
type ProtectionRuleExclusion struct {

	// The target of the exclusion.
	Target ProtectionRuleExclusionTargetEnum `mandatory:"false" json:"target,omitempty"`

	Exclusions []string `mandatory:"false" json:"exclusions"`
}

func (m ProtectionRuleExclusion) String() string {
	return common.PointerString(m)
}

// ProtectionRuleExclusionTargetEnum Enum with underlying type: string
type ProtectionRuleExclusionTargetEnum string

// Set of constants representing the allowable values for ProtectionRuleExclusionTargetEnum
const (
	ProtectionRuleExclusionTargetRequestCookies     ProtectionRuleExclusionTargetEnum = "REQUEST_COOKIES"
	ProtectionRuleExclusionTargetRequestCookieNames ProtectionRuleExclusionTargetEnum = "REQUEST_COOKIE_NAMES"
	ProtectionRuleExclusionTargetArgs               ProtectionRuleExclusionTargetEnum = "ARGS"
	ProtectionRuleExclusionTargetArgsNames          ProtectionRuleExclusionTargetEnum = "ARGS_NAMES"
)

var mappingProtectionRuleExclusionTarget = map[string]ProtectionRuleExclusionTargetEnum{
	"REQUEST_COOKIES":      ProtectionRuleExclusionTargetRequestCookies,
	"REQUEST_COOKIE_NAMES": ProtectionRuleExclusionTargetRequestCookieNames,
	"ARGS":                 ProtectionRuleExclusionTargetArgs,
	"ARGS_NAMES":           ProtectionRuleExclusionTargetArgsNames,
}

// GetProtectionRuleExclusionTargetEnumValues Enumerates the set of values for ProtectionRuleExclusionTargetEnum
func GetProtectionRuleExclusionTargetEnumValues() []ProtectionRuleExclusionTargetEnum {
	values := make([]ProtectionRuleExclusionTargetEnum, 0)
	for _, v := range mappingProtectionRuleExclusionTarget {
		values = append(values, v)
	}
	return values
}
