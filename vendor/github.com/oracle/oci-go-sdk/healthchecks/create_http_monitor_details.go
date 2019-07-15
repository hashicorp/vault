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

// CreateHttpMonitorDetails The request body used to create an HTTP monitor.
type CreateHttpMonitorDetails struct {

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	Targets []string `mandatory:"true" json:"targets"`

	Protocol CreateHttpMonitorDetailsProtocolEnum `mandatory:"true" json:"protocol"`

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

	Method CreateHttpMonitorDetailsMethodEnum `mandatory:"false" json:"method,omitempty"`

	// The optional URL path to probe, including query parameters.
	Path *string `mandatory:"false" json:"path"`

	// A dictionary of HTTP request headers.
	// *Note:* Monitors and probes do not support the use of the `Authorization` HTTP header.
	Headers map[string]string `mandatory:"false" json:"headers"`

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

func (m CreateHttpMonitorDetails) String() string {
	return common.PointerString(m)
}

// CreateHttpMonitorDetailsProtocolEnum Enum with underlying type: string
type CreateHttpMonitorDetailsProtocolEnum string

// Set of constants representing the allowable values for CreateHttpMonitorDetailsProtocolEnum
const (
	CreateHttpMonitorDetailsProtocolHttp  CreateHttpMonitorDetailsProtocolEnum = "HTTP"
	CreateHttpMonitorDetailsProtocolHttps CreateHttpMonitorDetailsProtocolEnum = "HTTPS"
)

var mappingCreateHttpMonitorDetailsProtocol = map[string]CreateHttpMonitorDetailsProtocolEnum{
	"HTTP":  CreateHttpMonitorDetailsProtocolHttp,
	"HTTPS": CreateHttpMonitorDetailsProtocolHttps,
}

// GetCreateHttpMonitorDetailsProtocolEnumValues Enumerates the set of values for CreateHttpMonitorDetailsProtocolEnum
func GetCreateHttpMonitorDetailsProtocolEnumValues() []CreateHttpMonitorDetailsProtocolEnum {
	values := make([]CreateHttpMonitorDetailsProtocolEnum, 0)
	for _, v := range mappingCreateHttpMonitorDetailsProtocol {
		values = append(values, v)
	}
	return values
}

// CreateHttpMonitorDetailsMethodEnum Enum with underlying type: string
type CreateHttpMonitorDetailsMethodEnum string

// Set of constants representing the allowable values for CreateHttpMonitorDetailsMethodEnum
const (
	CreateHttpMonitorDetailsMethodGet  CreateHttpMonitorDetailsMethodEnum = "GET"
	CreateHttpMonitorDetailsMethodHead CreateHttpMonitorDetailsMethodEnum = "HEAD"
)

var mappingCreateHttpMonitorDetailsMethod = map[string]CreateHttpMonitorDetailsMethodEnum{
	"GET":  CreateHttpMonitorDetailsMethodGet,
	"HEAD": CreateHttpMonitorDetailsMethodHead,
}

// GetCreateHttpMonitorDetailsMethodEnumValues Enumerates the set of values for CreateHttpMonitorDetailsMethodEnum
func GetCreateHttpMonitorDetailsMethodEnumValues() []CreateHttpMonitorDetailsMethodEnum {
	values := make([]CreateHttpMonitorDetailsMethodEnum, 0)
	for _, v := range mappingCreateHttpMonitorDetailsMethod {
		values = append(values, v)
	}
	return values
}
