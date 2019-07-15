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

// PingProbeResultSummary The results returned by running a ping probe.  All times and durations are
// returned in milliseconds. All times are relative to the POSIX epoch
// (1970-01-01T00:00Z).
type PingProbeResultSummary struct {

	// A value identifying this specific probe result. The key is only unique within
	// the results of its probe configuration. The key may be reused after 90 days.
	Key *string `mandatory:"false" json:"key"`

	// The OCID of the monitor or on-demand probe responsible for creating this result.
	ProbeConfigurationId *string `mandatory:"false" json:"probeConfigurationId"`

	// The date and time the probe was executed, expressed in milliseconds since the
	// POSIX epoch. This field is defined by the PerformanceResourceTiming interface
	// of the W3C Resource Timing specification. For more information, see
	// Resource Timing (https://w3c.github.io/resource-timing/#sec-resource-timing).
	StartTime *float64 `mandatory:"false" json:"startTime"`

	// The target hostname or IP address of the probe.
	Target *string `mandatory:"false" json:"target"`

	// The name of the vantage point that executed the probe.
	VantagePointName *string `mandatory:"false" json:"vantagePointName"`

	// True if the probe did not complete before the configured `timeoutInSeconds` value.
	IsTimedOut *bool `mandatory:"false" json:"isTimedOut"`

	// True if the probe result is determined to be healthy based on probe
	// type-specific criteria.  For HTTP probes, a probe result is considered
	// healthy if the HTTP response code is greater than or equal to 200 and
	// less than 300.
	IsHealthy *bool `mandatory:"false" json:"isHealthy"`

	// The category of error if an error occurs executing the probe.
	// The `errorMessage` field provides a message with the error details.
	// * NONE - No error
	// * DNS - DNS errors
	// * TRANSPORT - Transport-related errors, for example a "TLS certificate expired" error.
	// * NETWORK - Network-related errors, for example a "network unreachable" error.
	// * SYSTEM - Internal system errors.
	ErrorCategory PingProbeResultSummaryErrorCategoryEnum `mandatory:"false" json:"errorCategory,omitempty"`

	// The error information indicating why a probe execution failed.
	ErrorMessage *string `mandatory:"false" json:"errorMessage"`

	Protocol PingProbeResultSummaryProtocolEnum `mandatory:"false" json:"protocol,omitempty"`

	Connection *Connection `mandatory:"false" json:"connection"`

	Dns *Dns `mandatory:"false" json:"dns"`

	// The time immediately before the vantage point starts the domain name lookup for
	// the resource.
	DomainLookupStart *float64 `mandatory:"false" json:"domainLookupStart"`

	// The time immediately before the vantage point finishes the domain name lookup for
	// the resource.
	DomainLookupEnd *float64 `mandatory:"false" json:"domainLookupEnd"`

	// The latency of the probe execution, in milliseconds.
	LatencyInMs *float64 `mandatory:"false" json:"latencyInMs"`

	// The ICMP code of the response message.  This field is not used when the protocol
	// is set to TCP.  For more information on ICMP codes, see
	// Internet Control Message Protocol (ICMP) Parameters (https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml).
	IcmpCode *int `mandatory:"false" json:"icmpCode"`
}

func (m PingProbeResultSummary) String() string {
	return common.PointerString(m)
}

// PingProbeResultSummaryErrorCategoryEnum Enum with underlying type: string
type PingProbeResultSummaryErrorCategoryEnum string

// Set of constants representing the allowable values for PingProbeResultSummaryErrorCategoryEnum
const (
	PingProbeResultSummaryErrorCategoryNone      PingProbeResultSummaryErrorCategoryEnum = "NONE"
	PingProbeResultSummaryErrorCategoryDns       PingProbeResultSummaryErrorCategoryEnum = "DNS"
	PingProbeResultSummaryErrorCategoryTransport PingProbeResultSummaryErrorCategoryEnum = "TRANSPORT"
	PingProbeResultSummaryErrorCategoryNetwork   PingProbeResultSummaryErrorCategoryEnum = "NETWORK"
	PingProbeResultSummaryErrorCategorySystem    PingProbeResultSummaryErrorCategoryEnum = "SYSTEM"
)

var mappingPingProbeResultSummaryErrorCategory = map[string]PingProbeResultSummaryErrorCategoryEnum{
	"NONE":      PingProbeResultSummaryErrorCategoryNone,
	"DNS":       PingProbeResultSummaryErrorCategoryDns,
	"TRANSPORT": PingProbeResultSummaryErrorCategoryTransport,
	"NETWORK":   PingProbeResultSummaryErrorCategoryNetwork,
	"SYSTEM":    PingProbeResultSummaryErrorCategorySystem,
}

// GetPingProbeResultSummaryErrorCategoryEnumValues Enumerates the set of values for PingProbeResultSummaryErrorCategoryEnum
func GetPingProbeResultSummaryErrorCategoryEnumValues() []PingProbeResultSummaryErrorCategoryEnum {
	values := make([]PingProbeResultSummaryErrorCategoryEnum, 0)
	for _, v := range mappingPingProbeResultSummaryErrorCategory {
		values = append(values, v)
	}
	return values
}

// PingProbeResultSummaryProtocolEnum Enum with underlying type: string
type PingProbeResultSummaryProtocolEnum string

// Set of constants representing the allowable values for PingProbeResultSummaryProtocolEnum
const (
	PingProbeResultSummaryProtocolIcmp PingProbeResultSummaryProtocolEnum = "ICMP"
	PingProbeResultSummaryProtocolTcp  PingProbeResultSummaryProtocolEnum = "TCP"
)

var mappingPingProbeResultSummaryProtocol = map[string]PingProbeResultSummaryProtocolEnum{
	"ICMP": PingProbeResultSummaryProtocolIcmp,
	"TCP":  PingProbeResultSummaryProtocolTcp,
}

// GetPingProbeResultSummaryProtocolEnumValues Enumerates the set of values for PingProbeResultSummaryProtocolEnum
func GetPingProbeResultSummaryProtocolEnumValues() []PingProbeResultSummaryProtocolEnum {
	values := make([]PingProbeResultSummaryProtocolEnum, 0)
	for _, v := range mappingPingProbeResultSummaryProtocol {
		values = append(values, v)
	}
	return values
}
