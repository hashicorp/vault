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

// UpdateHostnameDetails The configuration details for updating a virtual hostname.
// For more information on virtual hostnames, see
// Managing Request Routing (https://docs.cloud.oracle.com/Content/Balance/Tasks/managingrequest.htm).
type UpdateHostnameDetails struct {

	// The virtual hostname to update. For more information about virtual hostname string construction, see
	// Managing Request Routing (https://docs.cloud.oracle.com/Content/Balance/Tasks/managingrequest.htm#routing).
	// Example: `app.example.com`
	Hostname *string `mandatory:"false" json:"hostname"`
}

func (m UpdateHostnameDetails) String() string {
	return common.PointerString(m)
}
