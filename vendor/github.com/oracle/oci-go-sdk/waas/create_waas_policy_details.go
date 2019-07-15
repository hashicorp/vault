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

// CreateWaasPolicyDetails The required data to create a WAAS policy.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type CreateWaasPolicyDetails struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment in which to create the WAAS policy.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The web application domain that the WAAS policy protects.
	Domain *string `mandatory:"true" json:"domain"`

	// A user-friendly name for the WAAS policy. The name is can be changed and does not need to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// An array of additional domains for the specified web application.
	AdditionalDomains []string `mandatory:"false" json:"additionalDomains"`

	// A map of host to origin for the web application. The key should be a customer friendly name for the host, ex. primary, secondary, etc.
	Origins map[string]Origin `mandatory:"false" json:"origins"`

	PolicyConfig *PolicyConfig `mandatory:"false" json:"policyConfig"`

	WafConfig *WafConfigDetails `mandatory:"false" json:"wafConfig"`

	// A simple key-value pair without any defined schema.
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// A key-value pair with a defined schema that restricts the values of tags. These predefined keys are scoped to namespaces.
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m CreateWaasPolicyDetails) String() string {
	return common.PointerString(m)
}
