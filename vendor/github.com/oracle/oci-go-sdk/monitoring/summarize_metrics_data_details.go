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

// SummarizeMetricsDataDetails The request details for retrieving aggregated data.
// Use the query and optional properties to filter the returned results.
type SummarizeMetricsDataDetails struct {

	// The source service or application to use when searching for metric data points to aggregate.
	// Example: `oci_computeagent`
	Namespace *string `mandatory:"true" json:"namespace"`

	// The Monitoring Query Language (MQL) expression to use when searching for metric data points to
	// aggregate. The query must specify a metric, statistic, and interval. Supported values for
	// interval: `1m`-`60m` (also `1h`). You can optionally specify dimensions and grouping functions.
	// Supported grouping functions: `grouping()`, `groupBy()`.
	// For details about Monitoring Query Language (MQL), see
	// Monitoring Query Language (MQL) Reference (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Reference/mql.htm).
	// For available dimensions, review the metric definition for the supported service.
	// See Supported Services (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm#SupportedServices).
	// Example: `CpuUtilization[1m].sum()`
	Query *string `mandatory:"true" json:"query"`

	// The beginning of the time range to use when searching for metric data points.
	// Format is defined by RFC3339. The response includes metric data points for the startTime.
	// Default value: the timestamp 3 hours before the call was sent.
	// Example: `2019-02-01T01:02:29.600Z`
	StartTime *common.SDKTime `mandatory:"false" json:"startTime"`

	// The end of the time range to use when searching for metric data points.
	// Format is defined by RFC3339. The response excludes metric data points for the endTime.
	// Default value: the timestamp representing when the call was sent.
	// Example: `2019-02-01T02:02:29.600Z`
	EndTime *common.SDKTime `mandatory:"false" json:"endTime"`

	// The time between calculated aggregation windows. Use with the query interval to vary the
	// frequency at which aggregated data points are returned. For example, use a query interval of
	// 5 minutes with a resolution of 1 minute to retrieve five-minute aggregations at a one-minute
	// frequency. The resolution must be equal or less than the interval in the query. The default
	// resolution is 1m (one minute). Supported values: `1m`-`60m` (also `1h`).
	// Example: `5m`
	Resolution *string `mandatory:"false" json:"resolution"`
}

func (m SummarizeMetricsDataDetails) String() string {
	return common.PointerString(m)
}
