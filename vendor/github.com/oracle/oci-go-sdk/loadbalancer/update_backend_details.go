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

// UpdateBackendDetails The configuration details for updating a backend server.
type UpdateBackendDetails struct {

	// The load balancing policy weight assigned to the server. Backend servers with a higher weight receive a larger
	// proportion of incoming traffic. For example, a server weighted '3' receives 3 times the number of new connections
	// as a server weighted '1'.
	// For more information on load balancing policies, see
	// How Load Balancing Policies Work (https://docs.cloud.oracle.com/Content/Balance/Reference/lbpolicies.htm).
	// Example: `3`
	Weight *int `mandatory:"true" json:"weight"`

	// Whether the load balancer should treat this server as a backup unit. If `true`, the load balancer forwards no ingress
	// traffic to this backend server unless all other backend servers not marked as "backup" fail the health check policy.
	// Example: `false`
	Backup *bool `mandatory:"true" json:"backup"`

	// Whether the load balancer should drain this server. Servers marked "drain" receive no new
	// incoming traffic.
	// Example: `false`
	Drain *bool `mandatory:"true" json:"drain"`

	// Whether the load balancer should treat this server as offline. Offline servers receive no incoming
	// traffic.
	// Example: `false`
	Offline *bool `mandatory:"true" json:"offline"`
}

func (m UpdateBackendDetails) String() string {
	return common.PointerString(m)
}
