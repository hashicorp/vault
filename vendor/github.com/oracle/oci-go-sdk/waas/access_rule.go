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

// AccessRule A content access rule. An access rule specifies an action to take if a set of criteria is matched by a request.
type AccessRule struct {

	// The unique name of the access rule.
	Name *string `mandatory:"true" json:"name"`

	// The list of access rule criteria.
	Criteria []AccessRuleCriteria `mandatory:"true" json:"criteria"`

	// The action to take when the access criteria are met for a rule. If unspecified, defaults to `ALLOW`.
	Action AccessRuleActionEnum `mandatory:"true" json:"action"`

	// The method used to block requests if `action` is set to `BLOCK` and the access criteria are met. If unspecified, defaults to `SET_RESPONSE_CODE`.
	BlockAction AccessRuleBlockActionEnum `mandatory:"false" json:"blockAction,omitempty"`

	// The response status code to return when `action` is set to `BLOCK`, `blockAction` is set to `SET_RESPONSE_CODE`, and the access criteria are met. If unspecified, defaults to `403`.
	BlockResponseCode *int `mandatory:"false" json:"blockResponseCode"`

	// The message to show on the error page when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_ERROR_PAGE`, and the access criteria are met. If unspecified, defaults to 'Access to the website is blocked.'
	BlockErrorPageMessage *string `mandatory:"false" json:"blockErrorPageMessage"`

	// The error code to show on the error page when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_ERROR_PAGE`, and the access criteria are met. If unspecified, defaults to 'Access rules'.
	BlockErrorPageCode *string `mandatory:"false" json:"blockErrorPageCode"`

	// The description text to show on the error page when `action` is set to `BLOCK`, `blockAction` is set to `SHOW_ERROR_PAGE`, and the access criteria are met. If unspecified, defaults to 'Access blocked by website owner. Please contact support.'
	BlockErrorPageDescription *string `mandatory:"false" json:"blockErrorPageDescription"`
}

func (m AccessRule) String() string {
	return common.PointerString(m)
}

// AccessRuleActionEnum Enum with underlying type: string
type AccessRuleActionEnum string

// Set of constants representing the allowable values for AccessRuleActionEnum
const (
	AccessRuleActionAllow  AccessRuleActionEnum = "ALLOW"
	AccessRuleActionDetect AccessRuleActionEnum = "DETECT"
	AccessRuleActionBlock  AccessRuleActionEnum = "BLOCK"
)

var mappingAccessRuleAction = map[string]AccessRuleActionEnum{
	"ALLOW":  AccessRuleActionAllow,
	"DETECT": AccessRuleActionDetect,
	"BLOCK":  AccessRuleActionBlock,
}

// GetAccessRuleActionEnumValues Enumerates the set of values for AccessRuleActionEnum
func GetAccessRuleActionEnumValues() []AccessRuleActionEnum {
	values := make([]AccessRuleActionEnum, 0)
	for _, v := range mappingAccessRuleAction {
		values = append(values, v)
	}
	return values
}

// AccessRuleBlockActionEnum Enum with underlying type: string
type AccessRuleBlockActionEnum string

// Set of constants representing the allowable values for AccessRuleBlockActionEnum
const (
	AccessRuleBlockActionSetResponseCode AccessRuleBlockActionEnum = "SET_RESPONSE_CODE"
	AccessRuleBlockActionShowErrorPage   AccessRuleBlockActionEnum = "SHOW_ERROR_PAGE"
)

var mappingAccessRuleBlockAction = map[string]AccessRuleBlockActionEnum{
	"SET_RESPONSE_CODE": AccessRuleBlockActionSetResponseCode,
	"SHOW_ERROR_PAGE":   AccessRuleBlockActionShowErrorPage,
}

// GetAccessRuleBlockActionEnumValues Enumerates the set of values for AccessRuleBlockActionEnum
func GetAccessRuleBlockActionEnumValues() []AccessRuleBlockActionEnum {
	values := make([]AccessRuleBlockActionEnum, 0)
	for _, v := range mappingAccessRuleBlockAction {
		values = append(values, v)
	}
	return values
}
