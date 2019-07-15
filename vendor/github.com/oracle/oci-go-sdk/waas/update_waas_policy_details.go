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

// UpdateWaasPolicyDetails Updates the configuration details of a WAAS policy.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type UpdateWaasPolicyDetails struct {

	// A user-friendly name for the WAAS policy. The name is can be changed and does not need to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// An array of additional domains protected by this WAAS policy.
	AdditionalDomains []string `mandatory:"false" json:"additionalDomains"`

	// A map of host to origin for the web application. The key should be a customer friendly name for the host, ex. primary, secondary, etc.
	Origins map[string]Origin `mandatory:"false" json:"origins"`

	PolicyConfig *PolicyConfig `mandatory:"false" json:"policyConfig"`

	WafConfig *WafConfig `mandatory:"false" json:"wafConfig"`

	// A simple key-value pair without any defined schema.
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m UpdateWaasPolicyDetails) String() string {
	return common.PointerString(m)
}
