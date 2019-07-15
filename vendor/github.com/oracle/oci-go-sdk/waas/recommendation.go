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

// Recommendation A recommended protection rule for a web application. This recommendation can be accepted to apply it to the Web Application Firewall configuration for this policy.
// Use the `POST /waasPolicies/{waasPolicyId}/actions/acceptWafConfigRecommendations` method to accept recommended protection rules.
type Recommendation struct {

	// The unique key for the recommended protection rule.
	Key *string `mandatory:"false" json:"key"`

	// The list of the ModSecurity rule IDs associated with the protection rule.
	// For more information about ModSecurity's open source WAF rules, see Mod Security's documentation (https://www.modsecurity.org/CRS/Documentation/index.html).
	ModSecurityRuleIds []string `mandatory:"false" json:"modSecurityRuleIds"`

	// The name of the recommended protection rule.
	Name *string `mandatory:"false" json:"name"`

	// The description of the recommended protection rule.
	Description *string `mandatory:"false" json:"description"`

	// The list of labels for the recommended protection rule.
	Labels []string `mandatory:"false" json:"labels"`

	// The recommended action to apply to the protection rule.
	RecommendedAction *string `mandatory:"false" json:"recommendedAction"`
}

func (m Recommendation) String() string {
	return common.PointerString(m)
}
