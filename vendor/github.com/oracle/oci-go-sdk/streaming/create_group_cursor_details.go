// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Streaming Service API
//
// The API for the Streaming Service.
//

package streaming

import (
	"github.com/oracle/oci-go-sdk/common"
)

// CreateGroupCursorDetails Object used to create a group cursor.
type CreateGroupCursorDetails struct {

	// The type of the cursor. This value is only used when the group is created.
	Type CreateGroupCursorDetailsTypeEnum `mandatory:"true" json:"type"`

	// Name of the consumer group.
	GroupName *string `mandatory:"true" json:"groupName"`

	// The time to consume from if type is AT_TIME.
	Time *common.SDKTime `mandatory:"false" json:"time"`

	// A unique identifier for the instance joining the consumer group. If an instanceName is not provided, a UUID will be generated and used.
	InstanceName *string `mandatory:"false" json:"instanceName"`

	// The amount of a consumer instance inactivity time, before partition reservations are released.
	TimeoutInMs *int `mandatory:"false" json:"timeoutInMs"`

	// When using consumer-groups, the default commit-on-get behaviour can be overriden by setting this value to false.
	// If disabled, a consumer must manually commit their cursors.
	CommitOnGet *bool `mandatory:"false" json:"commitOnGet"`
}

func (m CreateGroupCursorDetails) String() string {
	return common.PointerString(m)
}

// CreateGroupCursorDetailsTypeEnum Enum with underlying type: string
type CreateGroupCursorDetailsTypeEnum string

// Set of constants representing the allowable values for CreateGroupCursorDetailsTypeEnum
const (
	CreateGroupCursorDetailsTypeAtTime      CreateGroupCursorDetailsTypeEnum = "AT_TIME"
	CreateGroupCursorDetailsTypeLatest      CreateGroupCursorDetailsTypeEnum = "LATEST"
	CreateGroupCursorDetailsTypeTrimHorizon CreateGroupCursorDetailsTypeEnum = "TRIM_HORIZON"
)

var mappingCreateGroupCursorDetailsType = map[string]CreateGroupCursorDetailsTypeEnum{
	"AT_TIME":      CreateGroupCursorDetailsTypeAtTime,
	"LATEST":       CreateGroupCursorDetailsTypeLatest,
	"TRIM_HORIZON": CreateGroupCursorDetailsTypeTrimHorizon,
}

// GetCreateGroupCursorDetailsTypeEnumValues Enumerates the set of values for CreateGroupCursorDetailsTypeEnum
func GetCreateGroupCursorDetailsTypeEnumValues() []CreateGroupCursorDetailsTypeEnum {
	values := make([]CreateGroupCursorDetailsTypeEnum, 0)
	for _, v := range mappingCreateGroupCursorDetailsType {
		values = append(values, v)
	}
	return values
}
