// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Budgets API
//
// Use the Budgets API to manage budgets and budget alerts.
//

package budget

import (
	"github.com/oracle/oci-go-sdk/common"
)

// AlertRuleSummary The alert rule.
type AlertRuleSummary struct {

	// The OCID of the alert rule
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the budget
	BudgetId *string `mandatory:"true" json:"budgetId"`

	// The name of the alert rule.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// ACTUAL means the alert will trigger based on actual usage.
	// FORECAST means the alert will trigger based on predicted usage.
	Type AlertRuleSummaryTypeEnum `mandatory:"true" json:"type"`

	// The threshold for triggering the alert. If thresholdType is PERCENTAGE, the maximum value is 10000.
	Threshold *float32 `mandatory:"true" json:"threshold"`

	// The type of threshold.
	ThresholdType AlertRuleSummaryThresholdTypeEnum `mandatory:"true" json:"thresholdType"`

	// The current state of the alert rule.
	LifecycleState AlertRuleSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The audience that will received the alert when it triggers.
	Recipients *string `mandatory:"true" json:"recipients"`

	// Time when budget was created
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Time when budget was updated
	TimeUpdated *common.SDKTime `mandatory:"true" json:"timeUpdated"`

	// Custom message that will be sent when alert is triggered
	Message *string `mandatory:"false" json:"message"`

	// The description of the alert rule.
	Description *string `mandatory:"false" json:"description"`

	// Version of the alert rule. Starts from 1 and increments by 1.
	Version *int `mandatory:"false" json:"version"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m AlertRuleSummary) String() string {
	return common.PointerString(m)
}

// AlertRuleSummaryTypeEnum Enum with underlying type: string
type AlertRuleSummaryTypeEnum string

// Set of constants representing the allowable values for AlertRuleSummaryTypeEnum
const (
	AlertRuleSummaryTypeActual   AlertRuleSummaryTypeEnum = "ACTUAL"
	AlertRuleSummaryTypeForecast AlertRuleSummaryTypeEnum = "FORECAST"
)

var mappingAlertRuleSummaryType = map[string]AlertRuleSummaryTypeEnum{
	"ACTUAL":   AlertRuleSummaryTypeActual,
	"FORECAST": AlertRuleSummaryTypeForecast,
}

// GetAlertRuleSummaryTypeEnumValues Enumerates the set of values for AlertRuleSummaryTypeEnum
func GetAlertRuleSummaryTypeEnumValues() []AlertRuleSummaryTypeEnum {
	values := make([]AlertRuleSummaryTypeEnum, 0)
	for _, v := range mappingAlertRuleSummaryType {
		values = append(values, v)
	}
	return values
}

// AlertRuleSummaryThresholdTypeEnum Enum with underlying type: string
type AlertRuleSummaryThresholdTypeEnum string

// Set of constants representing the allowable values for AlertRuleSummaryThresholdTypeEnum
const (
	AlertRuleSummaryThresholdTypePercentage AlertRuleSummaryThresholdTypeEnum = "PERCENTAGE"
	AlertRuleSummaryThresholdTypeAbsolute   AlertRuleSummaryThresholdTypeEnum = "ABSOLUTE"
)

var mappingAlertRuleSummaryThresholdType = map[string]AlertRuleSummaryThresholdTypeEnum{
	"PERCENTAGE": AlertRuleSummaryThresholdTypePercentage,
	"ABSOLUTE":   AlertRuleSummaryThresholdTypeAbsolute,
}

// GetAlertRuleSummaryThresholdTypeEnumValues Enumerates the set of values for AlertRuleSummaryThresholdTypeEnum
func GetAlertRuleSummaryThresholdTypeEnumValues() []AlertRuleSummaryThresholdTypeEnum {
	values := make([]AlertRuleSummaryThresholdTypeEnum, 0)
	for _, v := range mappingAlertRuleSummaryThresholdType {
		values = append(values, v)
	}
	return values
}

// AlertRuleSummaryLifecycleStateEnum Enum with underlying type: string
type AlertRuleSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for AlertRuleSummaryLifecycleStateEnum
const (
	AlertRuleSummaryLifecycleStateActive   AlertRuleSummaryLifecycleStateEnum = "ACTIVE"
	AlertRuleSummaryLifecycleStateInactive AlertRuleSummaryLifecycleStateEnum = "INACTIVE"
)

var mappingAlertRuleSummaryLifecycleState = map[string]AlertRuleSummaryLifecycleStateEnum{
	"ACTIVE":   AlertRuleSummaryLifecycleStateActive,
	"INACTIVE": AlertRuleSummaryLifecycleStateInactive,
}

// GetAlertRuleSummaryLifecycleStateEnumValues Enumerates the set of values for AlertRuleSummaryLifecycleStateEnum
func GetAlertRuleSummaryLifecycleStateEnumValues() []AlertRuleSummaryLifecycleStateEnum {
	values := make([]AlertRuleSummaryLifecycleStateEnum, 0)
	for _, v := range mappingAlertRuleSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
