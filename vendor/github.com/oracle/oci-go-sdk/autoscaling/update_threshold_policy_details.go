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

// UpdateThresholdPolicyDetails The representation of UpdateThresholdPolicyDetails
type UpdateThresholdPolicyDetails struct {

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The capacity requirements of the autoscaling policy.
	Capacity *Capacity `mandatory:"false" json:"capacity"`

	Rules []UpdateConditionDetails `mandatory:"false" json:"rules"`
}

//GetDisplayName returns DisplayName
func (m UpdateThresholdPolicyDetails) GetDisplayName() *string {
	return m.DisplayName
}

//GetCapacity returns Capacity
func (m UpdateThresholdPolicyDetails) GetCapacity() *Capacity {
	return m.Capacity
}

func (m UpdateThresholdPolicyDetails) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m UpdateThresholdPolicyDetails) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeUpdateThresholdPolicyDetails UpdateThresholdPolicyDetails
	s := struct {
		DiscriminatorParam string `json:"policyType"`
		MarshalTypeUpdateThresholdPolicyDetails
	}{
		"threshold",
		(MarshalTypeUpdateThresholdPolicyDetails)(m),
	}

	return json.Marshal(&s)
}
