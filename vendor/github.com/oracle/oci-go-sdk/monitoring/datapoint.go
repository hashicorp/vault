// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Monitoring API
//
// Use the Monitoring API to manage metric queries and alarms for assessing the health, capacity, and performance of your cloud resources.
// For information about monitoring, see Monitoring Overview (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm).
//

package monitoring

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Datapoint Metric value for a specific timestamp.
type Datapoint struct {

	// Timestamp for this metric value. Format defined by RFC3339.
	// Example: `2019-02-01T01:02:29.600Z`
	Timestamp *common.SDKTime `mandatory:"true" json:"timestamp"`

	// Numeric value of the metric.
	// Example: `10.23`
	Value *float64 `mandatory:"true" json:"value"`

	// The number of occurrences of the associated value in the set of data.
	// Optional. Default is 1.
	Count *int `mandatory:"false" json:"count"`
}

func (m Datapoint) String() string {
	return common.PointerString(m)
}
