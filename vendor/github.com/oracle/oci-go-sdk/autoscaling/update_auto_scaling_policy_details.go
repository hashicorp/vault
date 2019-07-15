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

// UpdateAutoScalingPolicyDetails The representation of UpdateAutoScalingPolicyDetails
type UpdateAutoScalingPolicyDetails interface {

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
	GetDisplayName() *string

	// The capacity requirements of the autoscaling policy.
	GetCapacity() *Capacity
}

type updateautoscalingpolicydetails struct {
	JsonData    []byte
	DisplayName *string   `mandatory:"false" json:"displayName"`
	Capacity    *Capacity `mandatory:"false" json:"capacity"`
	PolicyType  string    `json:"policyType"`
}

// UnmarshalJSON unmarshals json
func (m *updateautoscalingpolicydetails) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalerupdateautoscalingpolicydetails updateautoscalingpolicydetails
	s := struct {
		Model Unmarshalerupdateautoscalingpolicydetails
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.DisplayName = s.Model.DisplayName
	m.Capacity = s.Model.Capacity
	m.PolicyType = s.Model.PolicyType

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *updateautoscalingpolicydetails) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.PolicyType {
	case "threshold":
		mm := UpdateThresholdPolicyDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

//GetDisplayName returns DisplayName
func (m updateautoscalingpolicydetails) GetDisplayName() *string {
	return m.DisplayName
}

//GetCapacity returns Capacity
func (m updateautoscalingpolicydetails) GetCapacity() *Capacity {
	return m.Capacity
}

func (m updateautoscalingpolicydetails) String() string {
	return common.PointerString(m)
}
