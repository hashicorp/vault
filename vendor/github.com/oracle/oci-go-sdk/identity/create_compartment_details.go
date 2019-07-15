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

// CreateCompartmentDetails The representation of CreateCompartmentDetails
type CreateCompartmentDetails struct {

	// The OCID of the parent compartment containing the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The name you assign to the compartment during creation. The name must be unique across all compartments
	// in the parent compartment. Avoid entering confidential information.
	Name *string `mandatory:"true" json:"name"`

	// The description you assign to the compartment during creation. Does not have to be unique, and it's changeable.
	Description *string `mandatory:"true" json:"description"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m CreateCompartmentDetails) String() string {
	return common.PointerString(m)
}
