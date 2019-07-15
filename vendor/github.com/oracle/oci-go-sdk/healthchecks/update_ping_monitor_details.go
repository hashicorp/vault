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

// UpdatePingMonitorDetails The request body used to update a ping monitor.
type UpdatePingMonitorDetails struct {
	Targets []string `mandatory:"false" json:"targets"`

	VantagePointNames []string `mandatory:"false" json:"vantagePointNames"`

	// The port on which to probe endpoints. If unspecified, probes will use the
	// default port of their protocol.
	Port *int `mandatory:"false" json:"port"`

	// The probe timeout in seconds. Valid values: 10, 20, 30, and 60.
	// The probe timeout must be less than or equal to `intervalInSeconds` for monitors.
	TimeoutInSeconds *int `mandatory:"false" json:"timeoutInSeconds"`

	Protocol UpdatePingMonitorDetailsProtocolEnum `mandatory:"false" json:"protocol,omitempty"`

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

func (m UpdatePingMonitorDetails) String() string {
	return common.PointerString(m)
}

// UpdatePingMonitorDetailsProtocolEnum Enum with underlying type: string
type UpdatePingMonitorDetailsProtocolEnum string

// Set of constants representing the allowable values for UpdatePingMonitorDetailsProtocolEnum
const (
	UpdatePingMonitorDetailsProtocolIcmp UpdatePingMonitorDetailsProtocolEnum = "ICMP"
	UpdatePingMonitorDetailsProtocolTcp  UpdatePingMonitorDetailsProtocolEnum = "TCP"
)

var mappingUpdatePingMonitorDetailsProtocol = map[string]UpdatePingMonitorDetailsProtocolEnum{
	"ICMP": UpdatePingMonitorDetailsProtocolIcmp,
	"TCP":  UpdatePingMonitorDetailsProtocolTcp,
}

// GetUpdatePingMonitorDetailsProtocolEnumValues Enumerates the set of values for UpdatePingMonitorDetailsProtocolEnum
func GetUpdatePingMonitorDetailsProtocolEnumValues() []UpdatePingMonitorDetailsProtocolEnum {
	values := make([]UpdatePingMonitorDetailsProtocolEnum, 0)
	for _, v := range mappingUpdatePingMonitorDetailsProtocol {
		values = append(values, v)
	}
	return values
}
