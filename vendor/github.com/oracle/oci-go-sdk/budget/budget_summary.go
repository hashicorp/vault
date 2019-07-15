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

// BudgetSummary A budget.
type BudgetSummary struct {

	// The OCID of the budget
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the compartment
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the compartment on which budget is applied
	TargetCompartmentId *string `mandatory:"true" json:"targetCompartmentId"`

	// The display name of the budget.
	DisplayName *string `mandatory:"true" json:"displayName"`

	// The amount of the budget expressed in the currency of the customer's rate card.
	Amount *float32 `mandatory:"true" json:"amount"`

	// The reset period for the budget.
	ResetPeriod BudgetSummaryResetPeriodEnum `mandatory:"true" json:"resetPeriod"`

	// The current state of the budget.
	LifecycleState BudgetSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Total number of alert rules in the budget
	AlertRuleCount *int `mandatory:"true" json:"alertRuleCount"`

	// Time budget was created
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Time budget was updated
	TimeUpdated *common.SDKTime `mandatory:"true" json:"timeUpdated"`

	// The description of the budget.
	Description *string `mandatory:"false" json:"description"`

	// Version of the budget. Starts from 1 and increments by 1.
	Version *int `mandatory:"false" json:"version"`

	// The actual spend in currency for the current budget cycle
	ActualSpend *float32 `mandatory:"false" json:"actualSpend"`

	// The forecasted spend in currency by the end of the current budget cycle
	ForecastedSpend *float32 `mandatory:"false" json:"forecastedSpend"`

	// Time budget spend was last computed
	TimeSpendComputed *common.SDKTime `mandatory:"false" json:"timeSpendComputed"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m BudgetSummary) String() string {
	return common.PointerString(m)
}

// BudgetSummaryResetPeriodEnum Enum with underlying type: string
type BudgetSummaryResetPeriodEnum string

// Set of constants representing the allowable values for BudgetSummaryResetPeriodEnum
const (
	BudgetSummaryResetPeriodMonthly BudgetSummaryResetPeriodEnum = "MONTHLY"
)

var mappingBudgetSummaryResetPeriod = map[string]BudgetSummaryResetPeriodEnum{
	"MONTHLY": BudgetSummaryResetPeriodMonthly,
}

// GetBudgetSummaryResetPeriodEnumValues Enumerates the set of values for BudgetSummaryResetPeriodEnum
func GetBudgetSummaryResetPeriodEnumValues() []BudgetSummaryResetPeriodEnum {
	values := make([]BudgetSummaryResetPeriodEnum, 0)
	for _, v := range mappingBudgetSummaryResetPeriod {
		values = append(values, v)
	}
	return values
}

// BudgetSummaryLifecycleStateEnum Enum with underlying type: string
type BudgetSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for BudgetSummaryLifecycleStateEnum
const (
	BudgetSummaryLifecycleStateActive   BudgetSummaryLifecycleStateEnum = "ACTIVE"
	BudgetSummaryLifecycleStateInactive BudgetSummaryLifecycleStateEnum = "INACTIVE"
)

var mappingBudgetSummaryLifecycleState = map[string]BudgetSummaryLifecycleStateEnum{
	"ACTIVE":   BudgetSummaryLifecycleStateActive,
	"INACTIVE": BudgetSummaryLifecycleStateInactive,
}

// GetBudgetSummaryLifecycleStateEnumValues Enumerates the set of values for BudgetSummaryLifecycleStateEnum
func GetBudgetSummaryLifecycleStateEnumValues() []BudgetSummaryLifecycleStateEnum {
	values := make([]BudgetSummaryLifecycleStateEnum, 0)
	for _, v := range mappingBudgetSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
