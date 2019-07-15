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

// TagDefaultSummary Summary information for the specified tag default.
type TagDefaultSummary struct {

	// The OCID of the tag default.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment. The tag default will apply to all new resources that are created in the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the tag namespace that contains the tag definition.
	TagNamespaceId *string `mandatory:"true" json:"tagNamespaceId"`

	// The OCID of the tag definition. The tag default will always assign a default value for this tag definition.
	TagDefinitionId *string `mandatory:"true" json:"tagDefinitionId"`

	// The name used in the tag definition. This field is informational in the context of the tag default.
	TagDefinitionName *string `mandatory:"true" json:"tagDefinitionName"`

	// The default value for the tag definition. This will be applied to all new resources created in the compartment.
	Value *string `mandatory:"true" json:"value"`

	// Date and time the `TagDefault` object was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The tag default's current state. After creating a `TagDefault`, make sure its `lifecycleState` is ACTIVE before using it.
	LifecycleState TagDefaultSummaryLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`
}

func (m TagDefaultSummary) String() string {
	return common.PointerString(m)
}

// TagDefaultSummaryLifecycleStateEnum Enum with underlying type: string
type TagDefaultSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for TagDefaultSummaryLifecycleStateEnum
const (
	TagDefaultSummaryLifecycleStateActive TagDefaultSummaryLifecycleStateEnum = "ACTIVE"
)

var mappingTagDefaultSummaryLifecycleState = map[string]TagDefaultSummaryLifecycleStateEnum{
	"ACTIVE": TagDefaultSummaryLifecycleStateActive,
}

// GetTagDefaultSummaryLifecycleStateEnumValues Enumerates the set of values for TagDefaultSummaryLifecycleStateEnum
func GetTagDefaultSummaryLifecycleStateEnumValues() []TagDefaultSummaryLifecycleStateEnum {
	values := make([]TagDefaultSummaryLifecycleStateEnum, 0)
	for _, v := range mappingTagDefaultSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
