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

// Capacity Capacity limits for the instance pool.
type Capacity struct {

	// The maximum number of instances the instance pool is allowed to increase to (scale out).
	Max *int `mandatory:"true" json:"max"`

	// The minimum number of instances the instance pool is allowed to decrease to (scale in).
	Min *int `mandatory:"true" json:"min"`

	// The initial number of instances to launch in the instance pool immediately after autoscaling is
	// enabled. After autoscaling retrieves performance metrics, the number of instances is automatically adjusted from this
	// initial number to a number that is based on the limits that you set.
	Initial *int `mandatory:"true" json:"initial"`
}

func (m Capacity) String() string {
	return common.PointerString(m)
}
