// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Identity and Access Management Service API
//
// APIs for managing users, groups, compartments, and policies.
//

package identity

import (
	"github.com/oracle/oci-go-sdk/common"
)

// CreatePolicyDetails The representation of CreatePolicyDetails
type CreatePolicyDetails struct {

	// The OCID of the compartment containing the policy (either the tenancy or another compartment).
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The name you assign to the policy during creation. The name must be unique across all policies
	// in the tenancy and cannot be changed.
	Name *string `mandatory:"true" json:"name"`

	// An array of policy statements written in the policy language. See
	// How Policies Work (https://docs.cloud.oracle.com/Content/Identity/Concepts/policies.htm) and
	// Common Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/commonpolicies.htm).
	Statements []string `mandatory:"true" json:"statements"`

	// The description you assign to the policy during creation. Does not have to be unique, and it's changeable.
	Description *string `mandatory:"true" json:"description"`

	// The version of the policy. If null or set to an empty string, when a request comes in for authorization, the
	// policy will be evaluated according to the current behavior of the services at that moment. If set to a particular
	// date (YYYY-MM-DD), the policy will be evaluated according to the behavior of the services on that date.
	VersionDate *common.SDKDate `mandatory:"false" json:"versionDate"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m CreatePolicyDetails) String() string {
	return common.PointerString(m)
}
