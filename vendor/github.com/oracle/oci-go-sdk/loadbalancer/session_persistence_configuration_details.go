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

// SessionPersistenceConfigurationDetails The configuration details for implementing session persistence based on a user-specified cookie name (application
// cookie stickiness).
// Session persistence enables the Load Balancing service to direct any number of requests that originate from a single
// logical client to a single backend web server. For more information, see
// Session Persistence (https://docs.cloud.oracle.com/Content/Balance/Reference/sessionpersistence.htm).
// To disable application cookie stickiness on a running load balancer, use the
// UpdateBackendSet operation and specify `null` for the
// `SessionPersistenceConfigurationDetails` object.
// Example: `SessionPersistenceConfigurationDetails: null`
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type SessionPersistenceConfigurationDetails struct {

	// The name of the cookie used to detect a session initiated by the backend server. Use '*' to specify
	// that any cookie set by the backend causes the session to persist.
	// Example: `example_cookie`
	CookieName *string `mandatory:"true" json:"cookieName"`

	// Whether the load balancer is prevented from directing traffic from a persistent session client to
	// a different backend server if the original server is unavailable. Defaults to false.
	// Example: `false`
	DisableFallback *bool `mandatory:"false" json:"disableFallback"`
}

func (m SessionPersistenceConfigurationDetails) String() string {
	return common.PointerString(m)
}
