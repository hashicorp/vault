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

// AlertType Valid values are ACTUAL (the alert will trigger based on actual usage) or
// FORECAST (the alert will trigger based on predicted usage).
type AlertType struct {
}

func (m AlertType) String() string {
	return common.PointerString(m)
}
