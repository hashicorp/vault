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

// CreateConditionDetails Creation details for a condition in a threshold-based autoscaling policy.
type CreateConditionDetails struct {
	Action *Action `mandatory:"true" json:"action"`

	Metric *Metric `mandatory:"true" json:"metric"`

	// A user-friendly name. Does not have to be unique, and it's changeable. Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`
}

func (m CreateConditionDetails) String() string {
	return common.PointerString(m)
}
