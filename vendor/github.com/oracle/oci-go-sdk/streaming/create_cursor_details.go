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

// CreateCursorDetails Object used to create a cursor to consume messages in a stream.
type CreateCursorDetails struct {

	// The partition to get messages from.
	Partition *string `mandatory:"true" json:"partition"`

	// The type of cursor, which determines the starting point from which the stream will be consumed:
	// - `AFTER_OFFSET:` The partition position immediately following the offset you specify. (Offsets are assigned when you successfully append a message to a partition in a stream.)
	// - `AT_OFFSET:` The exact partition position indicated by the offset you specify.
	// - `AT_TIME:` A specific point in time.
	// - `LATEST:` The most recent message in the partition that was added after the cursor was created.
	// - `TRIM_HORIZON:` The oldest message in the partition that is within the retention period window.
	Type CreateCursorDetailsTypeEnum `mandatory:"true" json:"type"`

	// The offset to consume from if the cursor type is `AT_OFFSET` or `AFTER_OFFSET`.
	Offset *int64 `mandatory:"false" json:"offset"`

	// The time to consume from if the cursor type is `AT_TIME`, expressed in RFC 3339 (https://tools.ietf.org/rfc/rfc3339) timestamp format.
	Time *common.SDKTime `mandatory:"false" json:"time"`
}

func (m CreateCursorDetails) String() string {
	return common.PointerString(m)
}

// CreateCursorDetailsTypeEnum Enum with underlying type: string
type CreateCursorDetailsTypeEnum string

// Set of constants representing the allowable values for CreateCursorDetailsTypeEnum
const (
	CreateCursorDetailsTypeAfterOffset CreateCursorDetailsTypeEnum = "AFTER_OFFSET"
	CreateCursorDetailsTypeAtOffset    CreateCursorDetailsTypeEnum = "AT_OFFSET"
	CreateCursorDetailsTypeAtTime      CreateCursorDetailsTypeEnum = "AT_TIME"
	CreateCursorDetailsTypeLatest      CreateCursorDetailsTypeEnum = "LATEST"
	CreateCursorDetailsTypeTrimHorizon CreateCursorDetailsTypeEnum = "TRIM_HORIZON"
)

var mappingCreateCursorDetailsType = map[string]CreateCursorDetailsTypeEnum{
	"AFTER_OFFSET": CreateCursorDetailsTypeAfterOffset,
	"AT_OFFSET":    CreateCursorDetailsTypeAtOffset,
	"AT_TIME":      CreateCursorDetailsTypeAtTime,
	"LATEST":       CreateCursorDetailsTypeLatest,
	"TRIM_HORIZON": CreateCursorDetailsTypeTrimHorizon,
}

// GetCreateCursorDetailsTypeEnumValues Enumerates the set of values for CreateCursorDetailsTypeEnum
func GetCreateCursorDetailsTypeEnumValues() []CreateCursorDetailsTypeEnum {
	values := make([]CreateCursorDetailsTypeEnum, 0)
	for _, v := range mappingCreateCursorDetailsType {
		values = append(values, v)
	}
	return values
}
