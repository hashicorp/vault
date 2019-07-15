// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Health Checks API
//
// API for the Health Checks service. Use this API to manage endpoint probes and monitors.
// For more information, see
// Overview of the Health Checks Service (https://docs.cloud.oracle.com/iaas/Content/HealthChecks/Concepts/healthchecks.htm).
//

package healthchecks

import (
	"github.com/oracle/oci-go-sdk/common"
)

// PingMonitor A summary containing all of the mutable and immutable properties for a ping monitor.
type PingMonitor struct {

	// The OCID of the resource.
	Id *string `mandatory:"false" json:"id"`

	// A URL for fetching the probe results.
	ResultsUrl *string `mandatory:"false" json:"resultsUrl"`

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	Targets []string `mandatory:"false" json:"targets"`

	VantagePointNames []string `mandatory:"false" json:"vantagePointNames"`

	// The port on which to probe endpoints. If unspecified, probes will use the
	// default port of their protocol.
	Port *int `mandatory:"false" json:"port"`

	// The probe timeout in seconds. Valid values: 10, 20, 30, and 60.
	// The probe timeout must be less than or equal to `intervalInSeconds` for monitors.
	TimeoutInSeconds *int `mandatory:"false" json:"timeoutInSeconds"`

	Protocol PingMonitorProtocolEnum `mandatory:"false" json:"protocol,omitempty"`

	// A user-friendly and mutable name suitable for display in a user interface.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The monitor interval in seconds. Valid values: 10, 30, and 60.
	IntervalInSeconds *int `mandatory:"false" json:"intervalInSeconds"`

	// Enables or disables the monitor. Set to 'true' to launch monitoring.
	IsEnabled *bool `mandatory:"false" json:"isEnabled"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace.  For more information,
	// see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m PingMonitor) String() string {
	return common.PointerString(m)
}

// PingMonitorProtocolEnum Enum with underlying type: string
type PingMonitorProtocolEnum string

// Set of constants representing the allowable values for PingMonitorProtocolEnum
const (
	PingMonitorProtocolIcmp PingMonitorProtocolEnum = "ICMP"
	PingMonitorProtocolTcp  PingMonitorProtocolEnum = "TCP"
)

var mappingPingMonitorProtocol = map[string]PingMonitorProtocolEnum{
	"ICMP": PingMonitorProtocolIcmp,
	"TCP":  PingMonitorProtocolTcp,
}

// GetPingMonitorProtocolEnumValues Enumerates the set of values for PingMonitorProtocolEnum
func GetPingMonitorProtocolEnumValues() []PingMonitorProtocolEnum {
	values := make([]PingMonitorProtocolEnum, 0)
	for _, v := range mappingPingMonitorProtocol {
		values = append(values, v)
	}
	return values
}
