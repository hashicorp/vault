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

// Policy A document that specifies the type of access a group has to the resources in a compartment. For information about
// policies and other IAM Service components, see
// Overview of the IAM Service (https://docs.cloud.oracle.com/Content/Identity/Concepts/overview.htm). If you're new to policies, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// The word "policy" is used by people in different ways:
//   * An individual statement written in the policy language
//   * A collection of statements in a single, named "policy" document (which has an Oracle Cloud ID (OCID) assigned to it)
//   * The overall body of policies your organization uses to control access to resources
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator.
type Policy struct {

	// The OCID of the policy.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment containing the policy (either the tenancy or another compartment).
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The name you assign to the policy during creation. The name must be unique across all policies
	// in the tenancy and cannot be changed.
	Name *string `mandatory:"true" json:"name"`

	// An array of one or more policy statements written in the policy language.
	Statements []string `mandatory:"true" json:"statements"`

	// The description you assign to the policy. Does not have to be unique, and it's changeable.
	Description *string `mandatory:"true" json:"description"`

	// Date and time the policy was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The policy's current state. After creating a policy, make sure its `lifecycleState` changes from CREATING to
	// ACTIVE before using it.
	LifecycleState PolicyLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The detailed status of INACTIVE lifecycleState.
	InactiveStatus *int64 `mandatory:"false" json:"inactiveStatus"`

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

func (m Policy) String() string {
	return common.PointerString(m)
}

// PolicyLifecycleStateEnum Enum with underlying type: string
type PolicyLifecycleStateEnum string

// Set of constants representing the allowable values for PolicyLifecycleStateEnum
const (
	PolicyLifecycleStateCreating PolicyLifecycleStateEnum = "CREATING"
	PolicyLifecycleStateActive   PolicyLifecycleStateEnum = "ACTIVE"
	PolicyLifecycleStateInactive PolicyLifecycleStateEnum = "INACTIVE"
	PolicyLifecycleStateDeleting PolicyLifecycleStateEnum = "DELETING"
	PolicyLifecycleStateDeleted  PolicyLifecycleStateEnum = "DELETED"
)

var mappingPolicyLifecycleState = map[string]PolicyLifecycleStateEnum{
	"CREATING": PolicyLifecycleStateCreating,
	"ACTIVE":   PolicyLifecycleStateActive,
	"INACTIVE": PolicyLifecycleStateInactive,
	"DELETING": PolicyLifecycleStateDeleting,
	"DELETED":  PolicyLifecycleStateDeleted,
}

// GetPolicyLifecycleStateEnumValues Enumerates the set of values for PolicyLifecycleStateEnum
func GetPolicyLifecycleStateEnumValues() []PolicyLifecycleStateEnum {
	values := make([]PolicyLifecycleStateEnum, 0)
	for _, v := range mappingPolicyLifecycleState {
		values = append(values, v)
	}
	return values
}
