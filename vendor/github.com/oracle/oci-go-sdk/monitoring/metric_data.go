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

// MetricData The set of aggregated data returned for a metric.
// For information about metrics, see Metrics Overview (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm#MetricsOverview).
type MetricData struct {

	// The reference provided in a metric definition to indicate the source service or
	// application that emitted the metric.
	// Example: `oci_computeagent`
	Namespace *string `mandatory:"true" json:"namespace"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment containing the
	// resources from which the aggregated data was returned.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The name of the metric.
	// Example: `CpuUtilization`
	Name *string `mandatory:"true" json:"name"`

	// Qualifiers provided in the definition of the returned metric.
	// Available dimensions vary by metric namespace. Each dimension takes the form of a key-value pair.
	// Example: `"resourceId": "ocid1.instance.region1.phx.exampleuniqueID"`
	Dimensions map[string]string `mandatory:"true" json:"dimensions"`

	// The list of timestamp-value pairs returned for the specified request. Metric values are rolled up to the start time specified in the request.
	AggregatedDatapoints []AggregatedDatapoint `mandatory:"true" json:"aggregatedDatapoints"`

	// The references provided in a metric definition to indicate extra information about the metric.
	// Example: `"unit": "bytes"`
	Metadata map[string]string `mandatory:"false" json:"metadata"`

	// The time between calculated aggregation windows. Use with the query interval to vary the
	// frequency at which aggregated data points are returned. For example, use a query interval of
	// 5 minutes with a resolution of 1 minute to retrieve five-minute aggregations at a one-minute
	// frequency. The resolution must be equal or less than the interval in the query. The default
	// resolution is 1m (one minute). Supported values: `1m`-`60m` (also `1h`).
	// Example: `5m`
	Resolution *string `mandatory:"false" json:"resolution"`
}

func (m MetricData) String() string {
	return common.PointerString(m)
}
