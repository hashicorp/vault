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

// TagNamespace A managed container for defined tags. A tag namespace is unique in a tenancy. A tag namespace can't be deleted.
// For more information, see Managing Tags and Tag Namespaces (https://docs.cloud.oracle.com/Content/Identity/Concepts/taggingoverview.htm).
type TagNamespace struct {

	// The OCID of the tag namespace.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment that contains the tag namespace.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The name of the tag namespace. It must be unique across all tag namespaces in the tenancy and cannot be changed.
	Name *string `mandatory:"true" json:"name"`

	// The description you assign to the tag namespace.
	Description *string `mandatory:"true" json:"description"`

	// Whether the tag namespace is retired.
	// See Retiring Key Definitions and Namespace Definitions (https://docs.cloud.oracle.com/Content/Identity/Concepts/taggingoverview.htm#Retiring).
	IsRetired *bool `mandatory:"true" json:"isRetired"`

	// Date and time the tagNamespace was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The tagnamespace's current state. After creating a tagnamespace, make sure its `lifecycleState` is ACTIVE before using it. After retiring a tagnamespace, make sure its `lifecycleState` is INACTIVE before using it.
	LifecycleState TagNamespaceLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`
}

func (m TagNamespace) String() string {
	return common.PointerString(m)
}

// TagNamespaceLifecycleStateEnum Enum with underlying type: string
type TagNamespaceLifecycleStateEnum string

// Set of constants representing the allowable values for TagNamespaceLifecycleStateEnum
const (
	TagNamespaceLifecycleStateActive   TagNamespaceLifecycleStateEnum = "ACTIVE"
	TagNamespaceLifecycleStateInactive TagNamespaceLifecycleStateEnum = "INACTIVE"
	TagNamespaceLifecycleStateDeleting TagNamespaceLifecycleStateEnum = "DELETING"
	TagNamespaceLifecycleStateDeleted  TagNamespaceLifecycleStateEnum = "DELETED"
)

var mappingTagNamespaceLifecycleState = map[string]TagNamespaceLifecycleStateEnum{
	"ACTIVE":   TagNamespaceLifecycleStateActive,
	"INACTIVE": TagNamespaceLifecycleStateInactive,
	"DELETING": TagNamespaceLifecycleStateDeleting,
	"DELETED":  TagNamespaceLifecycleStateDeleted,
}

// GetTagNamespaceLifecycleStateEnumValues Enumerates the set of values for TagNamespaceLifecycleStateEnum
func GetTagNamespaceLifecycleStateEnumValues() []TagNamespaceLifecycleStateEnum {
	values := make([]TagNamespaceLifecycleStateEnum, 0)
	for _, v := range mappingTagNamespaceLifecycleState {
		values = append(values, v)
	}
	return values
}
