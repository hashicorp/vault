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

// AlertRule The alert rule.
type AlertRule struct {

	// The OCID of the alert rule
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the budget
	BudgetId *string `mandatory:"true" json:"budgetId"`

	// The name of the alert rule.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The type of alert. Valid values are ACTUAL (the alert will trigger based on actual usage) or
	// FORECAST (the alert will trigger based on predicted usage).
	Type AlertRuleTypeEnum `mandatory:"true" json:"type"`

	// The threshold for triggering the alert. If thresholdType is PERCENTAGE, the maximum value is 10000.
	Threshold *float32 `mandatory:"true" json:"threshold"`

	// The type of threshold.
	ThresholdType AlertRuleThresholdTypeEnum `mandatory:"true" json:"thresholdType"`

	// The current state of the alert rule.
	LifecycleState AlertRuleLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Delimited list of email addresses to receive the alert when it triggers.
	// Delimiter character can be comma, space, TAB, or semicolon.
	Recipients *string `mandatory:"true" json:"recipients"`

	// Time budget was created
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Time budget was updated
	TimeUpdated *common.SDKTime `mandatory:"true" json:"timeUpdated"`

	// Custom message sent when alert is triggered
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

func (m AlertRule) String() string {
	return common.PointerString(m)
}

// AlertRuleTypeEnum Enum with underlying type: string
type AlertRuleTypeEnum string

// Set of constants representing the allowable values for AlertRuleTypeEnum
const (
	AlertRuleTypeActual   AlertRuleTypeEnum = "ACTUAL"
	AlertRuleTypeForecast AlertRuleTypeEnum = "FORECAST"
)

var mappingAlertRuleType = map[string]AlertRuleTypeEnum{
	"ACTUAL":   AlertRuleTypeActual,
	"FORECAST": AlertRuleTypeForecast,
}

// GetAlertRuleTypeEnumValues Enumerates the set of values for AlertRuleTypeEnum
func GetAlertRuleTypeEnumValues() []AlertRuleTypeEnum {
	values := make([]AlertRuleTypeEnum, 0)
	for _, v := range mappingAlertRuleType {
		values = append(values, v)
	}
	return values
}

// AlertRuleThresholdTypeEnum Enum with underlying type: string
type AlertRuleThresholdTypeEnum string

// Set of constants representing the allowable values for AlertRuleThresholdTypeEnum
const (
	AlertRuleThresholdTypePercentage AlertRuleThresholdTypeEnum = "PERCENTAGE"
	AlertRuleThresholdTypeAbsolute   AlertRuleThresholdTypeEnum = "ABSOLUTE"
)

var mappingAlertRuleThresholdType = map[string]AlertRuleThresholdTypeEnum{
	"PERCENTAGE": AlertRuleThresholdTypePercentage,
	"ABSOLUTE":   AlertRuleThresholdTypeAbsolute,
}

// GetAlertRuleThresholdTypeEnumValues Enumerates the set of values for AlertRuleThresholdTypeEnum
func GetAlertRuleThresholdTypeEnumValues() []AlertRuleThresholdTypeEnum {
	values := make([]AlertRuleThresholdTypeEnum, 0)
	for _, v := range mappingAlertRuleThresholdType {
		values = append(values, v)
	}
	return values
}

// AlertRuleLifecycleStateEnum Enum with underlying type: string
type AlertRuleLifecycleStateEnum string

// Set of constants representing the allowable values for AlertRuleLifecycleStateEnum
const (
	AlertRuleLifecycleStateActive   AlertRuleLifecycleStateEnum = "ACTIVE"
	AlertRuleLifecycleStateInactive AlertRuleLifecycleStateEnum = "INACTIVE"
)

var mappingAlertRuleLifecycleState = map[string]AlertRuleLifecycleStateEnum{
	"ACTIVE":   AlertRuleLifecycleStateActive,
	"INACTIVE": AlertRuleLifecycleStateInactive,
}

// GetAlertRuleLifecycleStateEnumValues Enumerates the set of values for AlertRuleLifecycleStateEnum
func GetAlertRuleLifecycleStateEnumValues() []AlertRuleLifecycleStateEnum {
	values := make([]AlertRuleLifecycleStateEnum, 0)
	for _, v := range mappingAlertRuleLifecycleState {
		values = append(values, v)
	}
	return values
}
