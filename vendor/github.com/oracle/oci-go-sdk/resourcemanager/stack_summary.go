// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Resource Manager API
//
// API for the Resource Manager service. Use this API to install, configure, and manage resources via the "infrastructure-as-code" model. For more information, see Overview of Resource Manager (https://docs.cloud.oracle.com/iaas/Content/ResourceManager/Concepts/resourcemanager.htm).
//

package resourcemanager

import (
	"github.com/oracle/oci-go-sdk/common"
)

// StackSummary Returns a list of properties and the defining property values for the specified stack.
type StackSummary struct {

	// Unique identifier of the specified stack.
	Id *string `mandatory:"false" json:"id"`

	// Unique identifier of the compartment in which the stack resides.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// Human-readable display name for the stack.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// General description of the stack.
	Description *string `mandatory:"false" json:"description"`

	// Date and time at which the stack was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	LifecycleState StackLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// Free-form tags associated with this resource. Each tag is a key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m StackSummary) String() string {
	return common.PointerString(m)
}
