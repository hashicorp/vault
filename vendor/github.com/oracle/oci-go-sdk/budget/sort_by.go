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

// SortBy The field to sort by. If not specified, the default is timeCreated.
// The default sort order for timeCreated is DESC.
// The default sort order for displayName is ASC in alphanumeric order.
type SortBy struct {
}

func (m SortBy) String() string {
	return common.PointerString(m)
}
