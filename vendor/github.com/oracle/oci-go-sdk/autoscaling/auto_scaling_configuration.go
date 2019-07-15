// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Autoscaling API
//
// APIs for dynamically scaling Compute resources to meet application requirements.
// For information about the Compute service, see Overview of the Compute Service (https://docs.cloud.oracle.com/Content/Compute/Concepts/computeoverview.htm).
//

package autoscaling

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// AutoScalingConfiguration An autoscaling configuration allows you to dynamically scale the resources in a Compute instance pool.
// For more information, see Autoscaling (https://docs.cloud.oracle.com/iaas/Content/Compute/Tasks/autoscalinginstancepools.htm).
type AutoScalingConfiguration struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment containing the autoscaling configuration.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the autoscaling configuration.
	Id *string `mandatory:"true" json:"id"`

	Resource Resource `mandatory:"true" json:"resource"`

	// Autoscaling policy definitions for the autoscaling configuration. An autoscaling policy defines the criteria that
	// trigger autoscaling actions and the actions to take.
	// Each autoscaling configuration can have one autoscaling policy.
	Policies []AutoScalingPolicy `mandatory:"true" json:"policies"`

	// The date and time the AutoScalingConfiguration was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The minimum period of time to wait between scaling actions. The cooldown period gives the system time to stabilize
	// before rescaling. The minimum value is 300 seconds, which is also the default.
	CoolDownInSeconds *int `mandatory:"false" json:"coolDownInSeconds"`

	// Whether the autoscaling configuration is enabled.
	IsEnabled *bool `mandatory:"false" json:"isEnabled"`
}

func (m AutoScalingConfiguration) String() string {
	return common.PointerString(m)
}

// UnmarshalJSON unmarshals from json
func (m *AutoScalingConfiguration) UnmarshalJSON(data []byte) (e error) {
	model := struct {
		DefinedTags       map[string]map[string]interface{} `json:"definedTags"`
		DisplayName       *string                           `json:"displayName"`
		FreeformTags      map[string]string                 `json:"freeformTags"`
		CoolDownInSeconds *int                              `json:"coolDownInSeconds"`
		IsEnabled         *bool                             `json:"isEnabled"`
		CompartmentId     *string                           `json:"compartmentId"`
		Id                *string                           `json:"id"`
		Resource          resource                          `json:"resource"`
		Policies          []autoscalingpolicy               `json:"policies"`
		TimeCreated       *common.SDKTime                   `json:"timeCreated"`
	}{}

	e = json.Unmarshal(data, &model)
	if e != nil {
		return
	}
	m.DefinedTags = model.DefinedTags
	m.DisplayName = model.DisplayName
	m.FreeformTags = model.FreeformTags
	m.CoolDownInSeconds = model.CoolDownInSeconds
	m.IsEnabled = model.IsEnabled
	m.CompartmentId = model.CompartmentId
	m.Id = model.Id
	nn, e := model.Resource.UnmarshalPolymorphicJSON(model.Resource.JsonData)
	if e != nil {
		return
	}
	if nn != nil {
		m.Resource = nn.(Resource)
	} else {
		m.Resource = nil
	}
	m.Policies = make([]AutoScalingPolicy, len(model.Policies))
	for i, n := range model.Policies {
		nn, err := n.UnmarshalPolymorphicJSON(n.JsonData)
		if err != nil {
			return err
		}
		if nn != nil {
			m.Policies[i] = nn.(AutoScalingPolicy)
		} else {
			m.Policies[i] = nil
		}
	}
	m.TimeCreated = model.TimeCreated
	return
}
