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

// AlarmHistoryEntry An alarm history entry indicating a description of the entry and the time that the entry occurred.
// If the entry corresponds to a state transition, such as OK to Firing, then the entry also includes a transition timestamp.
type AlarmHistoryEntry struct {

	// Description for this alarm history entry. Avoid entering confidential information.
	// Example 1 - alarm state history entry: `The alarm state is FIRING`
	// Example 2 - alarm state transition history entry: `State transitioned from OK to Firing`
	Summary *string `mandatory:"true" json:"summary"`

	// Timestamp for this alarm history entry. Format defined by RFC3339.
	// Example: `2019-02-01T01:02:29.600Z`
	Timestamp *common.SDKTime `mandatory:"true" json:"timestamp"`

	// Timestamp for the transition of the alarm state. For example, the time when the alarm transitioned from OK to Firing.
	// Available for state transition entries only. Note: A three-minute lag for this value accounts for any late-arriving metrics.
	// Example: `2019-02-01T0:59:00.789Z`
	TimestampTriggered *common.SDKTime `mandatory:"false" json:"timestampTriggered"`
}

func (m AlarmHistoryEntry) String() string {
	return common.PointerString(m)
}
