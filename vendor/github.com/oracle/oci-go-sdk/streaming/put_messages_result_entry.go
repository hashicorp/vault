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

// PutMessagesResultEntry Represents the result of a PutMessages request, whether it was successful or not.
// If a message was successfully appended to the stream, the entry includes the `offset`, `partition`, and `timestamp`.
// If the message failed to be appended to the stream, the entry includes the `error` and `errorMessage`.
type PutMessagesResultEntry struct {

	// The ID of the partition where the message was stored.
	Partition *string `mandatory:"false" json:"partition"`

	// The offset of the message in the partition.
	Offset *int64 `mandatory:"false" json:"offset"`

	// The timestamp indicating when the server appended the message to the stream.
	Timestamp *common.SDKTime `mandatory:"false" json:"timestamp"`

	// The error code, in case the message was not successfully appended to the stream.
	Error *string `mandatory:"false" json:"error"`

	// A human-readable error message associated with the error code.
	ErrorMessage *string `mandatory:"false" json:"errorMessage"`
}

func (m PutMessagesResultEntry) String() string {
	return common.PointerString(m)
}
