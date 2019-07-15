// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Autoscaling API
//
// APIs for dynamically scaling Compute resources to meet application requirements.
// For information about the Compute service, see Overview of the Compute Service (https://docs.cloud.oracle.com/Content/Compute/Concepts/computeoverview.htm).
//

package autoscaling

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Condition A rule that defines a specific autoscaling action to take (scale in or scale out) and the metric that triggers that action.
type Condition struct {
	Action *Action `mandatory:"true" json:"action"`

	Metric *Metric `mandatory:"true" json:"metric"`

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// ID of the condition that is assigned after creation.
	Id *string `mandatory:"false" json:"id"`
}

func (m Condition) String() string {
	return common.PointerString(m)
}
