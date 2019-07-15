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

// CreateOnDemandHttpProbeDetails The request body used to create an on-demand HTTP probe.
type CreateOnDemandHttpProbeDetails struct {

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	Targets []string `mandatory:"true" json:"targets"`

	Protocol CreateOnDemandHttpProbeDetailsProtocolEnum `mandatory:"true" json:"protocol"`

	VantagePointNames []string `mandatory:"false" json:"vantagePointNames"`

	// The port on which to probe endpoints. If unspecified, probes will use the
	// default port of their protocol.
	Port *int `mandatory:"false" json:"port"`

	// The probe timeout in seconds. Valid values: 10, 20, 30, and 60.
	// The probe timeout must be less than or equal to `intervalInSeconds` for monitors.
	TimeoutInSeconds *int `mandatory:"false" json:"timeoutInSeconds"`

	Method CreateOnDemandHttpProbeDetailsMethodEnum `mandatory:"false" json:"method,omitempty"`

	// The optional URL path to probe, including query parameters.
	Path *string `mandatory:"false" json:"path"`

	// A dictionary of HTTP request headers.
	// *Note:* Monitors and probes do not support the use of the `Authorization` HTTP header.
	Headers map[string]string `mandatory:"false" json:"headers"`
}

func (m CreateOnDemandHttpProbeDetails) String() string {
	return common.PointerString(m)
}

// CreateOnDemandHttpProbeDetailsProtocolEnum Enum with underlying type: string
type CreateOnDemandHttpProbeDetailsProtocolEnum string

// Set of constants representing the allowable values for CreateOnDemandHttpProbeDetailsProtocolEnum
const (
	CreateOnDemandHttpProbeDetailsProtocolHttp  CreateOnDemandHttpProbeDetailsProtocolEnum = "HTTP"
	CreateOnDemandHttpProbeDetailsProtocolHttps CreateOnDemandHttpProbeDetailsProtocolEnum = "HTTPS"
)

var mappingCreateOnDemandHttpProbeDetailsProtocol = map[string]CreateOnDemandHttpProbeDetailsProtocolEnum{
	"HTTP":  CreateOnDemandHttpProbeDetailsProtocolHttp,
	"HTTPS": CreateOnDemandHttpProbeDetailsProtocolHttps,
}

// GetCreateOnDemandHttpProbeDetailsProtocolEnumValues Enumerates the set of values for CreateOnDemandHttpProbeDetailsProtocolEnum
func GetCreateOnDemandHttpProbeDetailsProtocolEnumValues() []CreateOnDemandHttpProbeDetailsProtocolEnum {
	values := make([]CreateOnDemandHttpProbeDetailsProtocolEnum, 0)
	for _, v := range mappingCreateOnDemandHttpProbeDetailsProtocol {
		values = append(values, v)
	}
	return values
}

// CreateOnDemandHttpProbeDetailsMethodEnum Enum with underlying type: string
type CreateOnDemandHttpProbeDetailsMethodEnum string

// Set of constants representing the allowable values for CreateOnDemandHttpProbeDetailsMethodEnum
const (
	CreateOnDemandHttpProbeDetailsMethodGet  CreateOnDemandHttpProbeDetailsMethodEnum = "GET"
	CreateOnDemandHttpProbeDetailsMethodHead CreateOnDemandHttpProbeDetailsMethodEnum = "HEAD"
)

var mappingCreateOnDemandHttpProbeDetailsMethod = map[string]CreateOnDemandHttpProbeDetailsMethodEnum{
	"GET":  CreateOnDemandHttpProbeDetailsMethodGet,
	"HEAD": CreateOnDemandHttpProbeDetailsMethodHead,
}

// GetCreateOnDemandHttpProbeDetailsMethodEnumValues Enumerates the set of values for CreateOnDemandHttpProbeDetailsMethodEnum
func GetCreateOnDemandHttpProbeDetailsMethodEnumValues() []CreateOnDemandHttpProbeDetailsMethodEnum {
	values := make([]CreateOnDemandHttpProbeDetailsMethodEnum, 0)
	for _, v := range mappingCreateOnDemandHttpProbeDetailsMethod {
		values = append(values, v)
	}
	return values
}
