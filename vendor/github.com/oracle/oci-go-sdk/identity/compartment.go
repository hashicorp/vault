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

// Compartment A collection of related resources. Compartments are a fundamental component of Oracle Cloud Infrastructure
// for organizing and isolating your cloud resources. You use them to clearly separate resources for the purposes
// of measuring usage and billing, access (through the use of IAM Service policies), and isolation (separating the
// resources for one project or business unit from another). A common approach is to create a compartment for each
// major part of your organization. For more information, see
// Overview of the IAM Service (https://docs.cloud.oracle.com/Content/Identity/Concepts/overview.htm) and also
// Setting Up Your Tenancy (https://docs.cloud.oracle.com/Content/GSG/Concepts/settinguptenancy.htm).
// To place a resource in a compartment, simply specify the compartment ID in the "Create" request object when
// initially creating the resource. For example, to launch an instance into a particular compartment, specify
// that compartment's OCID in the `LaunchInstance` request. You can't move an existing resource from one
// compartment to another.
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access,
// see Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type Compartment struct {

	// The OCID of the compartment.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the parent compartment containing the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The name you assign to the compartment during creation. The name must be unique across all
	// compartments in the parent. Avoid entering confidential information.
	Name *string `mandatory:"true" json:"name"`

	// The description you assign to the compartment. Does not have to be unique, and it's changeable.
	Description *string `mandatory:"true" json:"description"`

	// Date and time the compartment was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The compartment's current state. After creating a compartment, make sure its `lifecycleState` changes from
	// CREATING to ACTIVE before using it.
	LifecycleState CompartmentLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`

	// Indicates whether or not the compartment is accessible for the user making the request.
	// Returns true when the user has INSPECT permissions directly on a resource in the
	// compartment or indirectly (permissions can be on a resource in a subcompartment).
	IsAccessible *bool `mandatory:"false" json:"isAccessible"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m Compartment) String() string {
	return common.PointerString(m)
}

// CompartmentLifecycleStateEnum Enum with underlying type: string
type CompartmentLifecycleStateEnum string

// Set of constants representing the allowable values for CompartmentLifecycleStateEnum
const (
	CompartmentLifecycleStateCreating CompartmentLifecycleStateEnum = "CREATING"
	CompartmentLifecycleStateActive   CompartmentLifecycleStateEnum = "ACTIVE"
	CompartmentLifecycleStateInactive CompartmentLifecycleStateEnum = "INACTIVE"
	CompartmentLifecycleStateDeleting CompartmentLifecycleStateEnum = "DELETING"
	CompartmentLifecycleStateDeleted  CompartmentLifecycleStateEnum = "DELETED"
)

var mappingCompartmentLifecycleState = map[string]CompartmentLifecycleStateEnum{
	"CREATING": CompartmentLifecycleStateCreating,
	"ACTIVE":   CompartmentLifecycleStateActive,
	"INACTIVE": CompartmentLifecycleStateInactive,
	"DELETING": CompartmentLifecycleStateDeleting,
	"DELETED":  CompartmentLifecycleStateDeleted,
}

// GetCompartmentLifecycleStateEnumValues Enumerates the set of values for CompartmentLifecycleStateEnum
func GetCompartmentLifecycleStateEnumValues() []CompartmentLifecycleStateEnum {
	values := make([]CompartmentLifecycleStateEnum, 0)
	for _, v := range mappingCompartmentLifecycleState {
		values = append(values, v)
	}
	return values
}
