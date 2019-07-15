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

// CreatePingMonitorDetails The request body used to create a Ping monitor.
type CreatePingMonitorDetails struct {

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	Targets []string `mandatory:"true" json:"targets"`

	Protocol CreatePingMonitorDetailsProtocolEnum `mandatory:"true" json:"protocol"`

	// A user-friendly and mutable name suitable for display in a user interface.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The monitor interval in seconds. Valid values: 10, 30, and 60.
	IntervalInSeconds *int `mandatory:"true" json:"intervalInSeconds"`

	VantagePointNames []string `mandatory:"false" json:"vantagePointNames"`

	// The port on which to probe endpoints. If unspecified, probes will use the
	// default port of their protocol.
	Port *int `mandatory:"false" json:"port"`

	// The probe timeout in seconds. Valid values: 10, 20, 30, and 60.
	// The probe timeout must be less than or equal to `intervalInSeconds` for monitors.
	TimeoutInSeconds *int `mandatory:"false" json:"timeoutInSeconds"`

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

func (m CreatePingMonitorDetails) String() string {
	return common.PointerString(m)
}

// CreatePingMonitorDetailsProtocolEnum Enum with underlying type: string
type CreatePingMonitorDetailsProtocolEnum string

// Set of constants representing the allowable values for CreatePingMonitorDetailsProtocolEnum
const (
	CreatePingMonitorDetailsProtocolIcmp CreatePingMonitorDetailsProtocolEnum = "ICMP"
	CreatePingMonitorDetailsProtocolTcp  CreatePingMonitorDetailsProtocolEnum = "TCP"
)

var mappingCreatePingMonitorDetailsProtocol = map[string]CreatePingMonitorDetailsProtocolEnum{
	"ICMP": CreatePingMonitorDetailsProtocolIcmp,
	"TCP":  CreatePingMonitorDetailsProtocolTcp,
}

// GetCreatePingMonitorDetailsProtocolEnumValues Enumerates the set of values for CreatePingMonitorDetailsProtocolEnum
func GetCreatePingMonitorDetailsProtocolEnumValues() []CreatePingMonitorDetailsProtocolEnum {
	values := make([]CreatePingMonitorDetailsProtocolEnum, 0)
	for _, v := range mappingCreatePingMonitorDetailsProtocol {
		values = append(values, v)
	}
	return values
}
