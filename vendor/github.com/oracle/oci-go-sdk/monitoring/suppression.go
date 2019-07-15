// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Monitoring API
//
// Use the Monitoring API to manage metric queries and alarms for assessing the health, capacity, and performance of your cloud resources.
// For information about monitoring, see Monitoring Overview (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm).
//

package monitoring

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Suppression The configuration details for suppressing an alarm.
// For information about alarms, see Alarms Overview (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm#AlarmsOverview).
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type Suppression struct {

	// The start date and time for the suppression to take place, inclusive. Format defined by RFC3339.
	// Example: `2019-02-01T01:02:29.600Z`
	TimeSuppressFrom *common.SDKTime `mandatory:"true" json:"timeSuppressFrom"`

	// The end date and time for the suppression to take place, inclusive. Format defined by RFC3339.
	// Example: `2019-02-01T02:02:29.600Z`
	TimeSuppressUntil *common.SDKTime `mandatory:"true" json:"timeSuppressUntil"`

	// Human-readable reason for suppressing alarm notifications.
	// It does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	// Oracle recommends including tracking information for the event or associated work,
	// such as a ticket number.
	// Example: `Planned outage due to change IT-1234.`
	Description *string `mandatory:"false" json:"description"`
}

func (m Suppression) String() string {
	return common.PointerString(m)
}
