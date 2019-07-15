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

// HealthChecker The health check policy configuration.
// For more information, see Editing Health Check Policies (https://docs.cloud.oracle.com/Content/Balance/Tasks/editinghealthcheck.htm).
type HealthChecker struct {

	// The protocol the health check must use; either HTTP or TCP.
	// Example: `HTTP`
	Protocol *string `mandatory:"true" json:"protocol"`

	// The backend server port against which to run the health check. If the port is not specified, the load balancer uses the
	// port information from the `Backend` object.
	// Example: `8080`
	Port *int `mandatory:"true" json:"port"`

	// The status code a healthy backend server should return. If you configure the health check policy to use the HTTP protocol,
	// you can use common HTTP status codes such as "200".
	// Example: `200`
	ReturnCode *int `mandatory:"true" json:"returnCode"`

	// A regular expression for parsing the response body from the backend server.
	// Example: `^((?!false).|\s)*$`
	ResponseBodyRegex *string `mandatory:"true" json:"responseBodyRegex"`

	// The path against which to run the health check.
	// Example: `/healthcheck`
	UrlPath *string `mandatory:"false" json:"urlPath"`

	// The number of retries to attempt before a backend server is considered "unhealthy". Defaults to 3.
	// Example: `3`
	Retries *int `mandatory:"false" json:"retries"`

	// The maximum time, in milliseconds, to wait for a reply to a health check. A health check is successful only if a reply
	// returns within this timeout period. Defaults to 3000 (3 seconds).
	// Example: `3000`
	TimeoutInMillis *int `mandatory:"false" json:"timeoutInMillis"`

	// The interval between health checks, in milliseconds. The default is 10000 (10 seconds).
	// Example: `10000`
	IntervalInMillis *int `mandatory:"false" json:"intervalInMillis"`
}

func (m HealthChecker) String() string {
	return common.PointerString(m)
}
