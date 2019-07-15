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

// TagNamespaceSummary A container for defined tags.
type TagNamespaceSummary struct {

	// The OCID of the tag namespace.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the compartment that contains the tag namespace.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The name of the tag namespace. It must be unique across all tag namespaces in the tenancy and cannot be changed.
	Name *string `mandatory:"false" json:"name"`

	// The description you assign to the tag namespace.
	Description *string `mandatory:"false" json:"description"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// Whether the tag namespace is retired.
	// For more information, see Retiring Key Definitions and Namespace Definitions (https://docs.cloud.oracle.com/Content/Identity/Concepts/taggingoverview.htm#Retiring).
	IsRetired *bool `mandatory:"false" json:"isRetired"`

	// The tagnamespace's current state. After creating a tagnamespace, make sure its `lifecycleState` is ACTIVE before using it. After retiring a tagnamespace, make sure its `lifecycleState` is INACTIVE before using it.
	LifecycleState TagNamespaceLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// Date and time the tag namespace was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m TagNamespaceSummary) String() string {
	return common.PointerString(m)
}
