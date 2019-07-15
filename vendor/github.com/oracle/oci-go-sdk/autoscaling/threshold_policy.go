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

// ThresholdPolicy An autoscaling policy that defines threshold-based rules for an autoscaling configuration.
type ThresholdPolicy struct {

	// The capacity requirements of the autoscaling policy.
	Capacity *Capacity `mandatory:"true" json:"capacity"`

	// The date and time the autoscaling configuration was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	Rules []Condition `mandatory:"true" json:"rules"`

	// The ID of the autoscaling policy that is assigned after creation.
	Id *string `mandatory:"false" json:"id"`

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`
}

//GetCapacity returns Capacity
func (m ThresholdPolicy) GetCapacity() *Capacity {
	return m.Capacity
}

//GetId returns Id
func (m ThresholdPolicy) GetId() *string {
	return m.Id
}

//GetDisplayName returns DisplayName
func (m ThresholdPolicy) GetDisplayName() *string {
	return m.DisplayName
}

//GetTimeCreated returns TimeCreated
func (m ThresholdPolicy) GetTimeCreated() *common.SDKTime {
	return m.TimeCreated
}

func (m ThresholdPolicy) String() string {
	return common.PointerString(m)
}

// MarshalJSON marshals to json representation
func (m ThresholdPolicy) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeThresholdPolicy ThresholdPolicy
	s := struct {
		DiscriminatorParam string `json:"policyType"`
		MarshalTypeThresholdPolicy
	}{
		"threshold",
		(MarshalTypeThresholdPolicy)(m),
	}

	return json.Marshal(&s)
}
