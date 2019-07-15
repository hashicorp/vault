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

// WaasPolicy The details of a Web Application Acceleration and Security (WAAS) policy. A policy describes how the WAAS service should operate for the configured web application.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type WaasPolicy struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the WAAS policy.
	Id *string `mandatory:"false" json:"id"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the WAAS policy's compartment.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The user-friendly name of the WAAS policy. The name can be changed and does not need to be unique.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The web application domain that the WAAS policy protects.
	Domain *string `mandatory:"false" json:"domain"`

	// An array of additional domains for this web application.
	AdditionalDomains []string `mandatory:"false" json:"additionalDomains"`

	// The CNAME record to add to your DNS configuration to route traffic for the domain, and all additional domains, through the WAF.
	Cname *string `mandatory:"false" json:"cname"`

	// The current lifecycle state of the WAAS policy.
	LifecycleState WaasPolicyLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// The date and time the policy was created, expressed in RFC 3339 timestamp format.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// A map of host to origin for the web application. The key should be a customer friendly name for the host, ex. primary, secondary, etc.
	Origins map[string]Origin `mandatory:"false" json:"origins"`

	PolicyConfig *PolicyConfig `mandatory:"false" json:"policyConfig"`

	WafConfig *WafConfig `mandatory:"false" json:"wafConfig"`

	// A simple key-value pair without any defined schema.
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// A key-value pair with a defined schema that restricts the values of tags. These predefined keys are scoped to namespaces.
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m WaasPolicy) String() string {
	return common.PointerString(m)
}

// WaasPolicyLifecycleStateEnum Enum with underlying type: string
type WaasPolicyLifecycleStateEnum string

// Set of constants representing the allowable values for WaasPolicyLifecycleStateEnum
const (
	WaasPolicyLifecycleStateCreating WaasPolicyLifecycleStateEnum = "CREATING"
	WaasPolicyLifecycleStateActive   WaasPolicyLifecycleStateEnum = "ACTIVE"
	WaasPolicyLifecycleStateFailed   WaasPolicyLifecycleStateEnum = "FAILED"
	WaasPolicyLifecycleStateUpdating WaasPolicyLifecycleStateEnum = "UPDATING"
	WaasPolicyLifecycleStateDeleting WaasPolicyLifecycleStateEnum = "DELETING"
	WaasPolicyLifecycleStateDeleted  WaasPolicyLifecycleStateEnum = "DELETED"
)

var mappingWaasPolicyLifecycleState = map[string]WaasPolicyLifecycleStateEnum{
	"CREATING": WaasPolicyLifecycleStateCreating,
	"ACTIVE":   WaasPolicyLifecycleStateActive,
	"FAILED":   WaasPolicyLifecycleStateFailed,
	"UPDATING": WaasPolicyLifecycleStateUpdating,
	"DELETING": WaasPolicyLifecycleStateDeleting,
	"DELETED":  WaasPolicyLifecycleStateDeleted,
}

// GetWaasPolicyLifecycleStateEnumValues Enumerates the set of values for WaasPolicyLifecycleStateEnum
func GetWaasPolicyLifecycleStateEnumValues() []WaasPolicyLifecycleStateEnum {
	values := make([]WaasPolicyLifecycleStateEnum, 0)
	for _, v := range mappingWaasPolicyLifecycleState {
		values = append(values, v)
	}
	return values
}
