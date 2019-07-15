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

// PostMetricDataDetails An array of metric objects containing raw metric data points to be posted to the Monitoring service.
type PostMetricDataDetails struct {

	// A metric object containing raw metric data points to be posted to the Monitoring service.
	MetricData []MetricDataDetails `mandatory:"true" json:"metricData"`

	// Batch atomicity behavior. Requires either partial or full pass of input validation for
	// metric objects in PostMetricData requests. The default value of NON_ATOMIC requires a
	// partial pass: at least one metric object in the request must pass input validation, and
	// any objects that failed validation are identified in the returned summary, along with
	// their error messages. A value of ATOMIC requires a full pass: all metric objects in
	// the request must pass input validation.
	// Example: `NON_ATOMIC`
	BatchAtomicity PostMetricDataDetailsBatchAtomicityEnum `mandatory:"false" json:"batchAtomicity,omitempty"`
}

func (m PostMetricDataDetails) String() string {
	return common.PointerString(m)
}

// PostMetricDataDetailsBatchAtomicityEnum Enum with underlying type: string
type PostMetricDataDetailsBatchAtomicityEnum string

// Set of constants representing the allowable values for PostMetricDataDetailsBatchAtomicityEnum
const (
	PostMetricDataDetailsBatchAtomicityAtomic    PostMetricDataDetailsBatchAtomicityEnum = "ATOMIC"
	PostMetricDataDetailsBatchAtomicityNonAtomic PostMetricDataDetailsBatchAtomicityEnum = "NON_ATOMIC"
)

var mappingPostMetricDataDetailsBatchAtomicity = map[string]PostMetricDataDetailsBatchAtomicityEnum{
	"ATOMIC":     PostMetricDataDetailsBatchAtomicityAtomic,
	"NON_ATOMIC": PostMetricDataDetailsBatchAtomicityNonAtomic,
}

// GetPostMetricDataDetailsBatchAtomicityEnumValues Enumerates the set of values for PostMetricDataDetailsBatchAtomicityEnum
func GetPostMetricDataDetailsBatchAtomicityEnumValues() []PostMetricDataDetailsBatchAtomicityEnum {
	values := make([]PostMetricDataDetailsBatchAtomicityEnum, 0)
	for _, v := range mappingPostMetricDataDetailsBatchAtomicity {
		values = append(values, v)
	}
	return values
}
