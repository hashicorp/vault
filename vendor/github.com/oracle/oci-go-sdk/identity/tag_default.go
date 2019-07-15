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

// TagDefault Tag defaults let you specify a default tag (tagnamespace.tag="value") to apply to all resource types
// in a specified compartment. The tag default is applied at the time the resource is created. Resources
// that exist in the compartment before you create the tag default are not tagged. The `TagDefault` object
// specifies the tag and compartment details.
// Tag defaults are inherited by child compartments. This means that if you set a tag default on the root compartment
// for a tenancy, all resources that are created in the tenancy are tagged. For more information about
// using tag defaults, see Managing Tag Defaults (https://docs.cloud.oracle.com/Content/Identity/Tasks/managingtagdefaults.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator.
type TagDefault struct {

	// The OCID of the tag default.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment. The tag default applies to all new resources that get created in the
	// compartment. Resources that existed before the tag default was created are not tagged.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the tag namespace that contains the tag definition.
	TagNamespaceId *string `mandatory:"true" json:"tagNamespaceId"`

	// The OCID of the tag definition. The tag default will always assign a default value for this tag definition.
	TagDefinitionId *string `mandatory:"true" json:"tagDefinitionId"`

	// The name used in the tag definition. This field is informational in the context of the tag default.
	TagDefinitionName *string `mandatory:"true" json:"tagDefinitionName"`

	// The default value for the tag definition. This will be applied to all resources created in the compartment.
	Value *string `mandatory:"true" json:"value"`

	// Date and time the `TagDefault` object was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The tag default's current state. After creating a `TagDefault`, make sure its `lifecycleState` is ACTIVE before using it.
	LifecycleState TagDefaultLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`
}

func (m TagDefault) String() string {
	return common.PointerString(m)
}

// TagDefaultLifecycleStateEnum Enum with underlying type: string
type TagDefaultLifecycleStateEnum string

// Set of constants representing the allowable values for TagDefaultLifecycleStateEnum
const (
	TagDefaultLifecycleStateActive TagDefaultLifecycleStateEnum = "ACTIVE"
)

var mappingTagDefaultLifecycleState = map[string]TagDefaultLifecycleStateEnum{
	"ACTIVE": TagDefaultLifecycleStateActive,
}

// GetTagDefaultLifecycleStateEnumValues Enumerates the set of values for TagDefaultLifecycleStateEnum
func GetTagDefaultLifecycleStateEnumValues() []TagDefaultLifecycleStateEnum {
	values := make([]TagDefaultLifecycleStateEnum, 0)
	for _, v := range mappingTagDefaultLifecycleState {
		values = append(values, v)
	}
	return values
}
