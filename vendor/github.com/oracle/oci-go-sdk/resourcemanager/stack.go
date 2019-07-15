// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Resource Manager API
//
// API for the Resource Manager service. Use this API to install, configure, and manage resources via the "infrastructure-as-code" model. For more information, see Overview of Resource Manager (https://docs.cloud.oracle.com/iaas/Content/ResourceManager/Concepts/resourcemanager.htm).
//

package resourcemanager

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// Stack The stack object. Stacks represent definitions of groups of Oracle Cloud Infrastructure
// resources that you can act upon as a group. You take action on stacks by using jobs.
type Stack struct {

	// Unique identifier (OCID) for the stack.
	Id *string `mandatory:"false" json:"id"`

	// Unique identifier (OCID) for the compartment where the stack is located.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// Human-readable name of the stack.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Description of the stack.
	Description *string `mandatory:"false" json:"description"`

	// The date and time at which the stack was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The current lifecycle state of the stack.
	LifecycleState StackLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// Specifies the `configSourceType` for uploading the Terraform configuration.
	// Presently, the .zip file type (`ZIP_UPLOAD`) is the only supported `configSourceType`.
	ConfigSource ConfigSource `mandatory:"false" json:"configSource"`

	// Terraform variables associated with this resource.
	// Maximum number of variables supported is 100.
	// The maximum size of each variable, including both name and value, is 4096 bytes.
	// Example: `{"CompartmentId": "compartment-id-value"}`
	Variables map[string]string `mandatory:"false" json:"variables"`

	// Free-form tags associated with the resource. Each tag is a key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m Stack) String() string {
	return common.PointerString(m)
}

// UnmarshalJSON unmarshals from json
func (m *Stack) UnmarshalJSON(data []byte) (e error) {
	model := struct {
		Id             *string                           `json:"id"`
		CompartmentId  *string                           `json:"compartmentId"`
		DisplayName    *string                           `json:"displayName"`
		Description    *string                           `json:"description"`
		TimeCreated    *common.SDKTime                   `json:"timeCreated"`
		LifecycleState StackLifecycleStateEnum           `json:"lifecycleState"`
		ConfigSource   configsource                      `json:"configSource"`
		Variables      map[string]string                 `json:"variables"`
		FreeformTags   map[string]string                 `json:"freeformTags"`
		DefinedTags    map[string]map[string]interface{} `json:"definedTags"`
	}{}

	e = json.Unmarshal(data, &model)
	if e != nil {
		return
	}
	m.Id = model.Id
	m.CompartmentId = model.CompartmentId
	m.DisplayName = model.DisplayName
	m.Description = model.Description
	m.TimeCreated = model.TimeCreated
	m.LifecycleState = model.LifecycleState
	nn, e := model.ConfigSource.UnmarshalPolymorphicJSON(model.ConfigSource.JsonData)
	if e != nil {
		return
	}
	if nn != nil {
		m.ConfigSource = nn.(ConfigSource)
	} else {
		m.ConfigSource = nil
	}
	m.Variables = model.Variables
	m.FreeformTags = model.FreeformTags
	m.DefinedTags = model.DefinedTags
	return
}

// StackLifecycleStateEnum Enum with underlying type: string
type StackLifecycleStateEnum string

// Set of constants representing the allowable values for StackLifecycleStateEnum
const (
	StackLifecycleStateCreating StackLifecycleStateEnum = "CREATING"
	StackLifecycleStateActive   StackLifecycleStateEnum = "ACTIVE"
	StackLifecycleStateDeleting StackLifecycleStateEnum = "DELETING"
	StackLifecycleStateDeleted  StackLifecycleStateEnum = "DELETED"
)

var mappingStackLifecycleState = map[string]StackLifecycleStateEnum{
	"CREATING": StackLifecycleStateCreating,
	"ACTIVE":   StackLifecycleStateActive,
	"DELETING": StackLifecycleStateDeleting,
	"DELETED":  StackLifecycleStateDeleted,
}

// GetStackLifecycleStateEnumValues Enumerates the set of values for StackLifecycleStateEnum
func GetStackLifecycleStateEnumValues() []StackLifecycleStateEnum {
	values := make([]StackLifecycleStateEnum, 0)
	for _, v := range mappingStackLifecycleState {
		values = append(values, v)
	}
	return values
}
