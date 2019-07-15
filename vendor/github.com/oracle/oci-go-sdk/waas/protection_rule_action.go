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

// ProtectionRuleAction A protection rule key and the associated action to apply to that rule.
type ProtectionRuleAction struct {

	// The unique key of the protection rule.
	Key *string `mandatory:"true" json:"key"`

	// The action to apply to the protection rule. If unspecified, defaults to `OFF`.
	Action ProtectionRuleActionActionEnum `mandatory:"true" json:"action"`

	// The types of requests excluded from the protection rule action. If the requests matches the criteria in the `exclusions`, the protection rule action will not be executed.
	Exclusions []ProtectionRuleExclusion `mandatory:"false" json:"exclusions"`
}

func (m ProtectionRuleAction) String() string {
	return common.PointerString(m)
}

// ProtectionRuleActionActionEnum Enum with underlying type: string
type ProtectionRuleActionActionEnum string

// Set of constants representing the allowable values for ProtectionRuleActionActionEnum
const (
	ProtectionRuleActionActionOff    ProtectionRuleActionActionEnum = "OFF"
	ProtectionRuleActionActionDetect ProtectionRuleActionActionEnum = "DETECT"
	ProtectionRuleActionActionBlock  ProtectionRuleActionActionEnum = "BLOCK"
)

var mappingProtectionRuleActionAction = map[string]ProtectionRuleActionActionEnum{
	"OFF":    ProtectionRuleActionActionOff,
	"DETECT": ProtectionRuleActionActionDetect,
	"BLOCK":  ProtectionRuleActionActionBlock,
}

// GetProtectionRuleActionActionEnumValues Enumerates the set of values for ProtectionRuleActionActionEnum
func GetProtectionRuleActionActionEnumValues() []ProtectionRuleActionActionEnum {
	values := make([]ProtectionRuleActionActionEnum, 0)
	for _, v := range mappingProtectionRuleActionAction {
		values = append(values, v)
	}
	return values
}
