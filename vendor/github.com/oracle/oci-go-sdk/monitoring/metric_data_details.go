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

// MetricDataDetails A metric object containing raw metric data points to be posted to the Monitoring service.
type MetricDataDetails struct {

	// The source service or application emitting the metric.
	// A valid namespace value starts with an alphabetical character and includes only alphanumeric characters and underscores. The "oci_" prefix is reserved.
	// Avoid entering confidential information.
	// Example: `my_namespace`
	Namespace *string `mandatory:"true" json:"namespace"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment to use for metrics.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The name of the metric.
	// A valid name value starts with an alphabetical character and includes only alphanumeric characters, dots, underscores, hyphens, and dollar signs.
	// Avoid entering confidential information.
	// Example: `my_app.success_rate`
	Name *string `mandatory:"true" json:"name"`

	// A list of metric values with timestamps. At least one data point is required per call.
	Datapoints []Datapoint `mandatory:"true" json:"datapoints"`

	// Qualifiers provided in a metric definition. Available dimensions vary by metric namespace.
	// Each dimension takes the form of a key-value pair. A valid dimension key includes only printable ASCII, excluding periods (.) and spaces. A valid dimension value includes only Unicode characters.
	// Empty strings are not allowed for keys or values. Avoid entering confidential information.
	// Example: `"resourceId": "ocid1.instance.region1.phx.exampleuniqueID"`
	Dimensions map[string]string `mandatory:"false" json:"dimensions"`

	// Properties describing metrics. These are not part of the unique fields identifying the metric.
	// Example: `"unit": "bytes"`
	Metadata map[string]string `mandatory:"false" json:"metadata"`
}

func (m MetricDataDetails) String() string {
	return common.PointerString(m)
}
