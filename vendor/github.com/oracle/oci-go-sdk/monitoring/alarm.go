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

// Alarm The properties that define an alarm.
// For information about alarms, see Alarms Overview (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm#AlarmsOverview).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/iaas/Content/Identity/Concepts/policygetstarted.htm).
// For information about endpoints and signing API requests, see
// About the API (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm). For information about available SDKs and tools, see
// SDKS and Other Tools (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdks.htm).
type Alarm struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the alarm.
	Id *string `mandatory:"true" json:"id"`

	// A user-friendly name for the alarm. It does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	// This name is sent as the title for notifications related to this alarm.
	// Example: `High CPU Utilization`
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment containing the alarm.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment containing the metric
	// being evaluated by the alarm.
	MetricCompartmentId *string `mandatory:"true" json:"metricCompartmentId"`

	// The source service or application emitting the metric that is evaluated by the alarm.
	// Example: `oci_computeagent`
	Namespace *string `mandatory:"true" json:"namespace"`

	// The Monitoring Query Language (MQL) expression to evaluate for the alarm. The Alarms feature of
	// the Monitoring service interprets results for each returned time series as Boolean values,
	// where zero represents false and a non-zero value represents true. A true value means that the trigger
	// rule condition has been met. The query must specify a metric, statistic, interval, and trigger
	// rule (threshold or absence). Supported values for interval: `1m`-`60m` (also `1h`). You can optionally
	// specify dimensions and grouping functions. Supported grouping functions: `grouping()`, `groupBy()`.
	// For details about Monitoring Query Language (MQL), see Monitoring Query Language (MQL) Reference (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Reference/mql.htm).
	// For available dimensions, review the metric definition for the supported service.
	// See Supported Services (https://docs.cloud.oracle.com/iaas/Content/Monitoring/Concepts/monitoringoverview.htm#SupportedServices).
	// Example of threshold alarm:
	//   -----
	//     CpuUtilization[1m]{availabilityDomain="cumS:PHX-AD-1"}.groupBy(availabilityDomain).percentile(0.9) > 85
	//   -----
	// Example of absence alarm:
	//   -----
	//     CpuUtilization[1m]{availabilityDomain="cumS:PHX-AD-1"}.absent()
	//   -----
	Query *string `mandatory:"true" json:"query"`

	// The perceived type of response required when the alarm is in the "FIRING" state.
	// Example: `CRITICAL`
	Severity AlarmSeverityEnum `mandatory:"true" json:"severity"`

	// An array of OCIDs (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) to which the notifications for
	// this alarm will be delivered. An example destination is an OCID for a topic managed by the
	// Oracle Cloud Infrastructure Notification service.
	Destinations []string `mandatory:"true" json:"destinations"`

	// Whether the alarm is enabled.
	// Example: `true`
	IsEnabled *bool `mandatory:"true" json:"isEnabled"`

	// The current lifecycle state of the alarm.
	// Example: `DELETED`
	LifecycleState AlarmLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the alarm was created. Format defined by RFC3339.
	// Example: `2019-02-01T01:02:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The date and time the alarm was last updated. Format defined by RFC3339.
	// Example: `2019-02-03T01:02:29.600Z`
	TimeUpdated *common.SDKTime `mandatory:"true" json:"timeUpdated"`

	// When true, the alarm evaluates metrics from all compartments and subcompartments. The parameter can
	// only be set to true when metricCompartmentId is the tenancy OCID (the tenancy is the root compartment).
	// A true value requires the user to have tenancy-level permissions. If this requirement is not met,
	// then the call is rejected. When false, the alarm evaluates metrics from only the compartment specified
	// in metricCompartmentId. Default is false.
	// Example: `true`
	MetricCompartmentIdInSubtree *bool `mandatory:"false" json:"metricCompartmentIdInSubtree"`

	// The time between calculated aggregation windows for the alarm. Supported value: `1m`
	Resolution *string `mandatory:"false" json:"resolution"`

	// The period of time that the condition defined in the alarm must persist before the alarm state
	// changes from "OK" to "FIRING" or vice versa. For example, a value of 5 minutes means that the
	// alarm must persist in breaching the condition for five minutes before the alarm updates its
	// state to "FIRING"; likewise, the alarm must persist in not breaching the condition for five
	// minutes before the alarm updates its state to "OK."
	// The duration is specified as a string in ISO 8601 format (`PT10M` for ten minutes or `PT1H`
	// for one hour). Minimum: PT1M. Maximum: PT1H. Default: PT1M.
	// Under the default value of PT1M, the first evaluation that breaches the alarm updates the
	// state to "FIRING" and the first evaluation that does not breach the alarm updates the state
	// to "OK".
	// Example: `PT5M`
	PendingDuration *string `mandatory:"false" json:"pendingDuration"`

	// The human-readable content of the notification delivered. Oracle recommends providing guidance
	// to operators for resolving the alarm condition. Consider adding links to standard runbook
	// practices. Avoid entering confidential information.
	// Example: `High CPU usage alert. Follow runbook instructions for resolution.`
	Body *string `mandatory:"false" json:"body"`

	// The frequency at which notifications are re-submitted, if the alarm keeps firing without
	// interruption. Format defined by ISO 8601. For example, `PT4H` indicates four hours.
	// Minimum: PT1M. Maximum: P30D.
	// Default value: null (notifications are not re-submitted).
	// Example: `PT2H`
	RepeatNotificationDuration *string `mandatory:"false" json:"repeatNotificationDuration"`

	// The configuration details for suppressing an alarm.
	Suppression *Suppression `mandatory:"false" json:"suppression"`

	// Simple key-value pair that is applied without any predefined name, type or scope. Exists for cross-compatibility only.
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Usage of predefined tag keys. These predefined keys are scoped to namespaces.
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m Alarm) String() string {
	return common.PointerString(m)
}

// AlarmSeverityEnum Enum with underlying type: string
type AlarmSeverityEnum string

// Set of constants representing the allowable values for AlarmSeverityEnum
const (
	AlarmSeverityCritical AlarmSeverityEnum = "CRITICAL"
	AlarmSeverityError    AlarmSeverityEnum = "ERROR"
	AlarmSeverityWarning  AlarmSeverityEnum = "WARNING"
	AlarmSeverityInfo     AlarmSeverityEnum = "INFO"
)

var mappingAlarmSeverity = map[string]AlarmSeverityEnum{
	"CRITICAL": AlarmSeverityCritical,
	"ERROR":    AlarmSeverityError,
	"WARNING":  AlarmSeverityWarning,
	"INFO":     AlarmSeverityInfo,
}

// GetAlarmSeverityEnumValues Enumerates the set of values for AlarmSeverityEnum
func GetAlarmSeverityEnumValues() []AlarmSeverityEnum {
	values := make([]AlarmSeverityEnum, 0)
	for _, v := range mappingAlarmSeverity {
		values = append(values, v)
	}
	return values
}

// AlarmLifecycleStateEnum Enum with underlying type: string
type AlarmLifecycleStateEnum string

// Set of constants representing the allowable values for AlarmLifecycleStateEnum
const (
	AlarmLifecycleStateActive   AlarmLifecycleStateEnum = "ACTIVE"
	AlarmLifecycleStateDeleting AlarmLifecycleStateEnum = "DELETING"
	AlarmLifecycleStateDeleted  AlarmLifecycleStateEnum = "DELETED"
)

var mappingAlarmLifecycleState = map[string]AlarmLifecycleStateEnum{
	"ACTIVE":   AlarmLifecycleStateActive,
	"DELETING": AlarmLifecycleStateDeleting,
	"DELETED":  AlarmLifecycleStateDeleted,
}

// GetAlarmLifecycleStateEnumValues Enumerates the set of values for AlarmLifecycleStateEnum
func GetAlarmLifecycleStateEnumValues() []AlarmLifecycleStateEnum {
	values := make([]AlarmLifecycleStateEnum, 0)
	for _, v := range mappingAlarmLifecycleState {
		values = append(values, v)
	}
	return values
}
