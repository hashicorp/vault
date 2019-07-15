// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ConsoleHistory An instance's serial console data. It includes configuration messages that occur when the
// instance boots, such as kernel and BIOS messages, and is useful for checking the status of
// the instance or diagnosing problems. The console data is minimally formatted ASCII text.
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type ConsoleHistory struct {

	// The availability domain of an instance.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the console history metadata object.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the instance this console history was fetched from.
	InstanceId *string `mandatory:"true" json:"instanceId"`

	// The current state of the console history.
	LifecycleState ConsoleHistoryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the history was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	// Example: `My console history metadata`
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`
}

func (m ConsoleHistory) String() string {
	return common.PointerString(m)
}

// ConsoleHistoryLifecycleStateEnum Enum with underlying type: string
type ConsoleHistoryLifecycleStateEnum string

// Set of constants representing the allowable values for ConsoleHistoryLifecycleStateEnum
const (
	ConsoleHistoryLifecycleStateRequested      ConsoleHistoryLifecycleStateEnum = "REQUESTED"
	ConsoleHistoryLifecycleStateGettingHistory ConsoleHistoryLifecycleStateEnum = "GETTING-HISTORY"
	ConsoleHistoryLifecycleStateSucceeded      ConsoleHistoryLifecycleStateEnum = "SUCCEEDED"
	ConsoleHistoryLifecycleStateFailed         ConsoleHistoryLifecycleStateEnum = "FAILED"
)

var mappingConsoleHistoryLifecycleState = map[string]ConsoleHistoryLifecycleStateEnum{
	"REQUESTED":       ConsoleHistoryLifecycleStateRequested,
	"GETTING-HISTORY": ConsoleHistoryLifecycleStateGettingHistory,
	"SUCCEEDED":       ConsoleHistoryLifecycleStateSucceeded,
	"FAILED":          ConsoleHistoryLifecycleStateFailed,
}

// GetConsoleHistoryLifecycleStateEnumValues Enumerates the set of values for ConsoleHistoryLifecycleStateEnum
func GetConsoleHistoryLifecycleStateEnumValues() []ConsoleHistoryLifecycleStateEnum {
	values := make([]ConsoleHistoryLifecycleStateEnum, 0)
	for _, v := range mappingConsoleHistoryLifecycleState {
		values = append(values, v)
	}
	return values
}
