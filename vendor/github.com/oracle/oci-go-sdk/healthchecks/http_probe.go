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

// HttpProbe A summary that contains all of the mutable and immutable properties for an HTTP probe.
type HttpProbe struct {

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

	Protocol HttpProbeProtocolEnum `mandatory:"false" json:"protocol,omitempty"`

	Method HttpProbeMethodEnum `mandatory:"false" json:"method,omitempty"`

	// The optional URL path to probe, including query parameters.
	Path *string `mandatory:"false" json:"path"`

	// A dictionary of HTTP request headers.
	// *Note:* Monitors and probes do not support the use of the `Authorization` HTTP header.
	Headers map[string]string `mandatory:"false" json:"headers"`
}

func (m HttpProbe) String() string {
	return common.PointerString(m)
}

// HttpProbeProtocolEnum Enum with underlying type: string
type HttpProbeProtocolEnum string

// Set of constants representing the allowable values for HttpProbeProtocolEnum
const (
	HttpProbeProtocolHttp  HttpProbeProtocolEnum = "HTTP"
	HttpProbeProtocolHttps HttpProbeProtocolEnum = "HTTPS"
)

var mappingHttpProbeProtocol = map[string]HttpProbeProtocolEnum{
	"HTTP":  HttpProbeProtocolHttp,
	"HTTPS": HttpProbeProtocolHttps,
}

// GetHttpProbeProtocolEnumValues Enumerates the set of values for HttpProbeProtocolEnum
func GetHttpProbeProtocolEnumValues() []HttpProbeProtocolEnum {
	values := make([]HttpProbeProtocolEnum, 0)
	for _, v := range mappingHttpProbeProtocol {
		values = append(values, v)
	}
	return values
}

// HttpProbeMethodEnum Enum with underlying type: string
type HttpProbeMethodEnum string

// Set of constants representing the allowable values for HttpProbeMethodEnum
const (
	HttpProbeMethodGet  HttpProbeMethodEnum = "GET"
	HttpProbeMethodHead HttpProbeMethodEnum = "HEAD"
)

var mappingHttpProbeMethod = map[string]HttpProbeMethodEnum{
	"GET":  HttpProbeMethodGet,
	"HEAD": HttpProbeMethodHead,
}

// GetHttpProbeMethodEnumValues Enumerates the set of values for HttpProbeMethodEnum
func GetHttpProbeMethodEnumValues() []HttpProbeMethodEnum {
	values := make([]HttpProbeMethodEnum, 0)
	for _, v := range mappingHttpProbeMethod {
		values = append(values, v)
	}
	return values
}
