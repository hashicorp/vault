// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Autoscaling API
//
// APIs for dynamically scaling Compute resources to meet application requirements.
// For information about the Compute service, see Overview of the Compute Service (https://docs.cloud.oracle.com/Content/Compute/Concepts/computeoverview.htm).
//

package autoscaling

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Action The action to take when autoscaling is triggered.
type Action struct {

	// The type of action to take.
	Type ActionTypeEnum `mandatory:"true" json:"type"`

	// To scale out (increase the number of instances), provide a positive value. To scale in (decrease the number of
	// instances), provide a negative value.
	Value *int `mandatory:"true" json:"value"`
}

func (m Action) String() string {
	return common.PointerString(m)
}

// ActionTypeEnum Enum with underlying type: string
type ActionTypeEnum string

// Set of constants representing the allowable values for ActionTypeEnum
const (
	ActionTypeBy ActionTypeEnum = "CHANGE_COUNT_BY"
)

var mappingActionType = map[string]ActionTypeEnum{
	"CHANGE_COUNT_BY": ActionTypeBy,
}

// GetActionTypeEnumValues Enumerates the set of values for ActionTypeEnum
func GetActionTypeEnumValues() []ActionTypeEnum {
	values := make([]ActionTypeEnum, 0)
	for _, v := range mappingActionType {
		values = append(values, v)
	}
	return values
}
