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

// AlarmStatusSummary A summary of properties for the specified alarm and its current evaluation status.
// For information about alarms, see Alarms Overview (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm#AlarmsOverview).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/iaas/Content/Identity/Concepts/policygetstarted.htm).
// For information about endpoints and signing API requests, see
// About the API (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm). For information about available SDKs and tools, see
// SDKS and Other Tools (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdks.htm).
type AlarmStatusSummary struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the alarm.
	Id *string `mandatory:"true" json:"id"`

	// The configured name of the alarm.
	// Example: `High CPU Utilization`
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The configured severity of the alarm.
	// Example: `CRITICAL`
	Severity AlarmStatusSummarySeverityEnum `mandatory:"true" json:"severity"`

	// Timestamp for the transition of the alarm state. For example, the time when the alarm transitioned from OK to Firing.
	// Example: `2019-02-01T01:02:29.600Z`
	TimestampTriggered *common.SDKTime `mandatory:"true" json:"timestampTriggered"`

	// The status of this alarm.
	// Example: `FIRING`
	Status AlarmStatusSummaryStatusEnum `mandatory:"true" json:"status"`

	// The configuration details for suppressing an alarm.
	Suppression *Suppression `mandatory:"false" json:"suppression"`
}

func (m AlarmStatusSummary) String() string {
	return common.PointerString(m)
}

// AlarmStatusSummarySeverityEnum Enum with underlying type: string
type AlarmStatusSummarySeverityEnum string

// Set of constants representing the allowable values for AlarmStatusSummarySeverityEnum
const (
	AlarmStatusSummarySeverityCritical AlarmStatusSummarySeverityEnum = "CRITICAL"
	AlarmStatusSummarySeverityError    AlarmStatusSummarySeverityEnum = "ERROR"
	AlarmStatusSummarySeverityWarning  AlarmStatusSummarySeverityEnum = "WARNING"
	AlarmStatusSummarySeverityInfo     AlarmStatusSummarySeverityEnum = "INFO"
)

var mappingAlarmStatusSummarySeverity = map[string]AlarmStatusSummarySeverityEnum{
	"CRITICAL": AlarmStatusSummarySeverityCritical,
	"ERROR":    AlarmStatusSummarySeverityError,
	"WARNING":  AlarmStatusSummarySeverityWarning,
	"INFO":     AlarmStatusSummarySeverityInfo,
}

// GetAlarmStatusSummarySeverityEnumValues Enumerates the set of values for AlarmStatusSummarySeverityEnum
func GetAlarmStatusSummarySeverityEnumValues() []AlarmStatusSummarySeverityEnum {
	values := make([]AlarmStatusSummarySeverityEnum, 0)
	for _, v := range mappingAlarmStatusSummarySeverity {
		values = append(values, v)
	}
	return values
}

// AlarmStatusSummaryStatusEnum Enum with underlying type: string
type AlarmStatusSummaryStatusEnum string

// Set of constants representing the allowable values for AlarmStatusSummaryStatusEnum
const (
	AlarmStatusSummaryStatusFiring    AlarmStatusSummaryStatusEnum = "FIRING"
	AlarmStatusSummaryStatusOk        AlarmStatusSummaryStatusEnum = "OK"
	AlarmStatusSummaryStatusSuspended AlarmStatusSummaryStatusEnum = "SUSPENDED"
)

var mappingAlarmStatusSummaryStatus = map[string]AlarmStatusSummaryStatusEnum{
	"FIRING":    AlarmStatusSummaryStatusFiring,
	"OK":        AlarmStatusSummaryStatusOk,
	"SUSPENDED": AlarmStatusSummaryStatusSuspended,
}

// GetAlarmStatusSummaryStatusEnumValues Enumerates the set of values for AlarmStatusSummaryStatusEnum
func GetAlarmStatusSummaryStatusEnumValues() []AlarmStatusSummaryStatusEnum {
	values := make([]AlarmStatusSummaryStatusEnum, 0)
	for _, v := range mappingAlarmStatusSummaryStatus {
		values = append(values, v)
	}
	return values
}
