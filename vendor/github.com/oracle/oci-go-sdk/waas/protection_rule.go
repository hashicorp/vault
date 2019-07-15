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

// ProtectionRule The protection rule settings. Protection rules can allow, block, or trigger an alert if a request meets the parameters of an applied rule.
type ProtectionRule struct {

	// The unique key of the protection rule.
	Key *string `mandatory:"false" json:"key"`

	// The list of the ModSecurity rule IDs that apply to this protection rule. For more information about ModSecurity's open source WAF rules, see Mod Security's documentation (https://www.modsecurity.org/CRS/Documentation/index.html).
	ModSecurityRuleIds []string `mandatory:"false" json:"modSecurityRuleIds"`

	// The name of the protection rule.
	Name *string `mandatory:"false" json:"name"`

	// The description of the protection rule.
	Description *string `mandatory:"false" json:"description"`

	// The action to take when the traffic is detected as malicious. If unspecified, defaults to `OFF`.
	Action ProtectionRuleActionEnum `mandatory:"false" json:"action,omitempty"`

	// The list of labels for the protection rule.
	// **Note:** Protection rules with a `ResponseBody` label will have no effect unless `isResponseInspected` is true.
	Labels []string `mandatory:"false" json:"labels"`

	Exclusions []ProtectionRuleExclusion `mandatory:"false" json:"exclusions"`
}

func (m ProtectionRule) String() string {
	return common.PointerString(m)
}

// ProtectionRuleActionEnum Enum with underlying type: string
type ProtectionRuleActionEnum string

// Set of constants representing the allowable values for ProtectionRuleActionEnum
const (
	ProtectionRuleActionOff    ProtectionRuleActionEnum = "OFF"
	ProtectionRuleActionDetect ProtectionRuleActionEnum = "DETECT"
	ProtectionRuleActionBlock  ProtectionRuleActionEnum = "BLOCK"
)

var mappingProtectionRuleAction = map[string]ProtectionRuleActionEnum{
	"OFF":    ProtectionRuleActionOff,
	"DETECT": ProtectionRuleActionDetect,
	"BLOCK":  ProtectionRuleActionBlock,
}

// GetProtectionRuleActionEnumValues Enumerates the set of values for ProtectionRuleActionEnum
func GetProtectionRuleActionEnumValues() []ProtectionRuleActionEnum {
	values := make([]ProtectionRuleActionEnum, 0)
	for _, v := range mappingProtectionRuleAction {
		values = append(values, v)
	}
	return values
}
