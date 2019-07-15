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

// BackendHealth The health status of the specified backend server as reported by the primary and standby load balancers.
type BackendHealth struct {

	// The general health status of the specified backend server as reported by the primary and standby load balancers.
	// *   **OK:** Both health checks returned `OK`.
	// *   **WARNING:** One health check returned `OK` and one did not.
	// *   **CRITICAL:** Neither health check returned `OK`.
	// *   **UNKNOWN:** One or both health checks returned `UNKNOWN`, or the system was unable to retrieve metrics at this time.
	Status BackendHealthStatusEnum `mandatory:"true" json:"status"`

	// A list of the most recent health check results returned for the specified backend server.
	HealthCheckResults []HealthCheckResult `mandatory:"true" json:"healthCheckResults"`
}

func (m BackendHealth) String() string {
	return common.PointerString(m)
}

// BackendHealthStatusEnum Enum with underlying type: string
type BackendHealthStatusEnum string

// Set of constants representing the allowable values for BackendHealthStatusEnum
const (
	BackendHealthStatusOk       BackendHealthStatusEnum = "OK"
	BackendHealthStatusWarning  BackendHealthStatusEnum = "WARNING"
	BackendHealthStatusCritical BackendHealthStatusEnum = "CRITICAL"
	BackendHealthStatusUnknown  BackendHealthStatusEnum = "UNKNOWN"
)

var mappingBackendHealthStatus = map[string]BackendHealthStatusEnum{
	"OK":       BackendHealthStatusOk,
	"WARNING":  BackendHealthStatusWarning,
	"CRITICAL": BackendHealthStatusCritical,
	"UNKNOWN":  BackendHealthStatusUnknown,
}

// GetBackendHealthStatusEnumValues Enumerates the set of values for BackendHealthStatusEnum
func GetBackendHealthStatusEnumValues() []BackendHealthStatusEnum {
	values := make([]BackendHealthStatusEnum, 0)
	for _, v := range mappingBackendHealthStatus {
		values = append(values, v)
	}
	return values
}
