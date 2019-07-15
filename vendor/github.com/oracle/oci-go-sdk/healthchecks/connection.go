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

// Connection The network connection results.
type Connection struct {

	// The connection IP address.
	Address *string `mandatory:"false" json:"address"`

	// The port.
	Port *int `mandatory:"false" json:"port"`
}

func (m Connection) String() string {
	return common.PointerString(m)
}
