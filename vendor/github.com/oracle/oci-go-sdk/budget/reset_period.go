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

// ResetPeriod The reset period for the budget. Valid value is MONTHLY.
type ResetPeriod struct {
}

func (m ResetPeriod) String() string {
	return common.PointerString(m)
}
