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

// Metric The properties that define a metric.
// For information about metrics, see Metrics Overview (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm#MetricsOverview).
type Metric struct {

	// The name of the metric.
	// Example: `CpuUtilization`
	Name *string `mandatory:"false" json:"name"`

	// The source service or application emitting the metric.
	// Example: `oci_computeagent`
	Namespace *string `mandatory:"false" json:"namespace"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment containing
	// the resources monitored by the metric.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// Qualifiers provided in a metric definition. Available dimensions vary by metric namespace.
	// Each dimension takes the form of a key-value pair.
	// Example: `"resourceId": "ocid1.instance.region1.phx.exampleuniqueID"`
	Dimensions map[string]string `mandatory:"false" json:"dimensions"`
}

func (m Metric) String() string {
	return common.PointerString(m)
}
