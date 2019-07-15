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

// HttpMonitor This model contains all of the mutable and immutable properties for an HTTP monitor.
type HttpMonitor struct {

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

	Protocol HttpMonitorProtocolEnum `mandatory:"false" json:"protocol,omitempty"`

	Method HttpMonitorMethodEnum `mandatory:"false" json:"method,omitempty"`

	// The optional URL path to probe, including query parameters.
	Path *string `mandatory:"false" json:"path"`

	// A dictionary of HTTP request headers.
	// *Note:* Monitors and probes do not support the use of the `Authorization` HTTP header.
	Headers map[string]string `mandatory:"false" json:"headers"`

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

func (m HttpMonitor) String() string {
	return common.PointerString(m)
}

// HttpMonitorProtocolEnum Enum with underlying type: string
type HttpMonitorProtocolEnum string

// Set of constants representing the allowable values for HttpMonitorProtocolEnum
const (
	HttpMonitorProtocolHttp  HttpMonitorProtocolEnum = "HTTP"
	HttpMonitorProtocolHttps HttpMonitorProtocolEnum = "HTTPS"
)

var mappingHttpMonitorProtocol = map[string]HttpMonitorProtocolEnum{
	"HTTP":  HttpMonitorProtocolHttp,
	"HTTPS": HttpMonitorProtocolHttps,
}

// GetHttpMonitorProtocolEnumValues Enumerates the set of values for HttpMonitorProtocolEnum
func GetHttpMonitorProtocolEnumValues() []HttpMonitorProtocolEnum {
	values := make([]HttpMonitorProtocolEnum, 0)
	for _, v := range mappingHttpMonitorProtocol {
		values = append(values, v)
	}
	return values
}

// HttpMonitorMethodEnum Enum with underlying type: string
type HttpMonitorMethodEnum string

// Set of constants representing the allowable values for HttpMonitorMethodEnum
const (
	HttpMonitorMethodGet  HttpMonitorMethodEnum = "GET"
	HttpMonitorMethodHead HttpMonitorMethodEnum = "HEAD"
)

var mappingHttpMonitorMethod = map[string]HttpMonitorMethodEnum{
	"GET":  HttpMonitorMethodGet,
	"HEAD": HttpMonitorMethodHead,
}

// GetHttpMonitorMethodEnumValues Enumerates the set of values for HttpMonitorMethodEnum
func GetHttpMonitorMethodEnumValues() []HttpMonitorMethodEnum {
	values := make([]HttpMonitorMethodEnum, 0)
	for _, v := range mappingHttpMonitorMethod {
		values = append(values, v)
	}
	return values
}
