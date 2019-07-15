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

// CreateStackDetails Properties provided for creating a stack.
type CreateStackDetails struct {

	// Unique identifier (OCID) of the compartment in which the stack resides.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	ConfigSource CreateConfigSourceDetails `mandatory:"true" json:"configSource"`

	// The stack's display name.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Description of the stack.
	Description *string `mandatory:"false" json:"description"`

	// Terraform variables associated with this resource.
	// Maximum number of variables supported is 100.
	// The maximum size of each variable, including both name and value, is 4096 bytes.
	// Example: `{"CompartmentId": "compartment-id-value"}`
	Variables map[string]string `mandatory:"false" json:"variables"`

	// Free-form tags associated with this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags associated with this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m CreateStackDetails) String() string {
	return common.PointerString(m)
}

// UnmarshalJSON unmarshals from json
func (m *CreateStackDetails) UnmarshalJSON(data []byte) (e error) {
	model := struct {
		DisplayName   *string                           `json:"displayName"`
		Description   *string                           `json:"description"`
		Variables     map[string]string                 `json:"variables"`
		FreeformTags  map[string]string                 `json:"freeformTags"`
		DefinedTags   map[string]map[string]interface{} `json:"definedTags"`
		CompartmentId *string                           `json:"compartmentId"`
		ConfigSource  createconfigsourcedetails         `json:"configSource"`
	}{}

	e = json.Unmarshal(data, &model)
	if e != nil {
		return
	}
	m.DisplayName = model.DisplayName
	m.Description = model.Description
	m.Variables = model.Variables
	m.FreeformTags = model.FreeformTags
	m.DefinedTags = model.DefinedTags
	m.CompartmentId = model.CompartmentId
	nn, e := model.ConfigSource.UnmarshalPolymorphicJSON(model.ConfigSource.JsonData)
	if e != nil {
		return
	}
	if nn != nil {
		m.ConfigSource = nn.(CreateConfigSourceDetails)
	} else {
		m.ConfigSource = nil
	}
	return
}
