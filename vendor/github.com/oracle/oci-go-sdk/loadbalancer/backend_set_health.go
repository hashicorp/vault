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

// BackendSetHealth The health status details for a backend set.
// This object does not explicitly enumerate backend servers with a status of `OK`. However, they are included in the
// `totalBackendCount` sum.
type BackendSetHealth struct {

	// Overall health status of the backend set.
	// *  **OK:** All backend servers in the backend set return a status of `OK`.
	// *  **WARNING:** Half or more of the backend set's backend servers return a status of `OK` and at least one backend
	// server returns a status of `WARNING`, `CRITICAL`, or `UNKNOWN`.
	// *  **CRITICAL:** Fewer than half of the backend set's backend servers return a status of `OK`.
	// *  **UNKNOWN:** More than half of the backend set's backend servers return a status of `UNKNOWN`, the system was
	// unable to retrieve metrics, or the backend set does not have a listener attached.
	Status BackendSetHealthStatusEnum `mandatory:"true" json:"status"`

	// A list of backend servers that are currently in the `WARNING` health state. The list identifies each backend server by
	// IP address and port.
	// Example: `10.0.0.3:8080`
	WarningStateBackendNames []string `mandatory:"true" json:"warningStateBackendNames"`

	// A list of backend servers that are currently in the `CRITICAL` health state. The list identifies each backend server by
	// IP address and port.
	// Example: `10.0.0.4:8080`
	CriticalStateBackendNames []string `mandatory:"true" json:"criticalStateBackendNames"`

	// A list of backend servers that are currently in the `UNKNOWN` health state. The list identifies each backend server by
	// IP address and port.
	// Example: `10.0.0.5:8080`
	UnknownStateBackendNames []string `mandatory:"true" json:"unknownStateBackendNames"`

	// The total number of backend servers in this backend set.
	// Example: `7`
	TotalBackendCount *int `mandatory:"true" json:"totalBackendCount"`
}

func (m BackendSetHealth) String() string {
	return common.PointerString(m)
}

// BackendSetHealthStatusEnum Enum with underlying type: string
type BackendSetHealthStatusEnum string

// Set of constants representing the allowable values for BackendSetHealthStatusEnum
const (
	BackendSetHealthStatusOk       BackendSetHealthStatusEnum = "OK"
	BackendSetHealthStatusWarning  BackendSetHealthStatusEnum = "WARNING"
	BackendSetHealthStatusCritical BackendSetHealthStatusEnum = "CRITICAL"
	BackendSetHealthStatusUnknown  BackendSetHealthStatusEnum = "UNKNOWN"
)

var mappingBackendSetHealthStatus = map[string]BackendSetHealthStatusEnum{
	"OK":       BackendSetHealthStatusOk,
	"WARNING":  BackendSetHealthStatusWarning,
	"CRITICAL": BackendSetHealthStatusCritical,
	"UNKNOWN":  BackendSetHealthStatusUnknown,
}

// GetBackendSetHealthStatusEnumValues Enumerates the set of values for BackendSetHealthStatusEnum
func GetBackendSetHealthStatusEnumValues() []BackendSetHealthStatusEnum {
	values := make([]BackendSetHealthStatusEnum, 0)
	for _, v := range mappingBackendSetHealthStatus {
		values = append(values, v)
	}
	return values
}
