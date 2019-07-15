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

// BackendSetDetails The configuration details for a load balancer backend set.
// For more information on backend set configuration, see
// Managing Backend Sets (https://docs.cloud.oracle.com/Content/Balance/Tasks/managingbackendsets.htm).
type BackendSetDetails struct {

	// The load balancer policy for the backend set. To get a list of available policies, use the
	// ListPolicies operation.
	// Example: `LEAST_CONNECTIONS`
	Policy *string `mandatory:"true" json:"policy"`

	HealthChecker *HealthCheckerDetails `mandatory:"true" json:"healthChecker"`

	Backends []BackendDetails `mandatory:"false" json:"backends"`

	SslConfiguration *SslConfigurationDetails `mandatory:"false" json:"sslConfiguration"`

	SessionPersistenceConfiguration *SessionPersistenceConfigurationDetails `mandatory:"false" json:"sessionPersistenceConfiguration"`
}

func (m BackendSetDetails) String() string {
	return common.PointerString(m)
}
