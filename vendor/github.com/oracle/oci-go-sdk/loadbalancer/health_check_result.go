// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Load Balancing API
//
// API for the Load Balancing service. Use this API to manage load balancers, backend sets, and related items. For more
// information, see Overview of Load Balancing (https://docs.cloud.oracle.com/iaas/Content/Balance/Concepts/balanceoverview.htm).
//

package loadbalancer

import (
	"github.com/oracle/oci-go-sdk/common"
)

// HealthCheckResult Information about a single backend server health check result reported by a load balancer.
type HealthCheckResult struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the subnet hosting the load balancer that reported this health check status.
	SubnetId *string `mandatory:"true" json:"subnetId"`

	// The IP address of the health check status report provider. This identifier helps you differentiate same-subnet
	// load balancers that report health check status.
	// Example: `10.0.0.7`
	SourceIpAddress *string `mandatory:"true" json:"sourceIpAddress"`

	// The date and time the data was retrieved, in the format defined by RFC3339.
	// Example: `2017-06-02T18:28:11+00:00`
	Timestamp *common.SDKTime `mandatory:"true" json:"timestamp"`

	// The result of the most recent health check.
	HealthCheckStatus HealthCheckResultHealthCheckStatusEnum `mandatory:"true" json:"healthCheckStatus"`
}

func (m HealthCheckResult) String() string {
	return common.PointerString(m)
}

// HealthCheckResultHealthCheckStatusEnum Enum with underlying type: string
type HealthCheckResultHealthCheckStatusEnum string

// Set of constants representing the allowable values for HealthCheckResultHealthCheckStatusEnum
const (
	HealthCheckResultHealthCheckStatusOk                HealthCheckResultHealthCheckStatusEnum = "OK"
	HealthCheckResultHealthCheckStatusInvalidStatusCode HealthCheckResultHealthCheckStatusEnum = "INVALID_STATUS_CODE"
	HealthCheckResultHealthCheckStatusTimedOut          HealthCheckResultHealthCheckStatusEnum = "TIMED_OUT"
	HealthCheckResultHealthCheckStatusRegexMismatch     HealthCheckResultHealthCheckStatusEnum = "REGEX_MISMATCH"
	HealthCheckResultHealthCheckStatusConnectFailed     HealthCheckResultHealthCheckStatusEnum = "CONNECT_FAILED"
	HealthCheckResultHealthCheckStatusIoError           HealthCheckResultHealthCheckStatusEnum = "IO_ERROR"
	HealthCheckResultHealthCheckStatusOffline           HealthCheckResultHealthCheckStatusEnum = "OFFLINE"
	HealthCheckResultHealthCheckStatusUnknown           HealthCheckResultHealthCheckStatusEnum = "UNKNOWN"
)

var mappingHealthCheckResultHealthCheckStatus = map[string]HealthCheckResultHealthCheckStatusEnum{
	"OK": HealthCheckResultHealthCheckStatusOk,
	"INVALID_STATUS_CODE": HealthCheckResultHealthCheckStatusInvalidStatusCode,
	"TIMED_OUT":           HealthCheckResultHealthCheckStatusTimedOut,
	"REGEX_MISMATCH":      HealthCheckResultHealthCheckStatusRegexMismatch,
	"CONNECT_FAILED":      HealthCheckResultHealthCheckStatusConnectFailed,
	"IO_ERROR":            HealthCheckResultHealthCheckStatusIoError,
	"OFFLINE":             HealthCheckResultHealthCheckStatusOffline,
	"UNKNOWN":             HealthCheckResultHealthCheckStatusUnknown,
}

// GetHealthCheckResultHealthCheckStatusEnumValues Enumerates the set of values for HealthCheckResultHealthCheckStatusEnum
func GetHealthCheckResultHealthCheckStatusEnumValues() []HealthCheckResultHealthCheckStatusEnum {
	values := make([]HealthCheckResultHealthCheckStatusEnum, 0)
	for _, v := range mappingHealthCheckResultHealthCheckStatus {
		values = append(values, v)
	}
	return values
}
