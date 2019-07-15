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

// PutMessagesResult The response to a PutMessages request. It indicates the number
// of failed messages as well as an array of results for successful and failed messages.
type PutMessagesResult struct {

	// The number of messages that failed to be added to the stream.
	Failures *int `mandatory:"true" json:"failures"`

	// An array of items representing the result of each message.
	// The order is guaranteed to be the same as in the `PutMessagesDetails` object.
	// If a message was successfully appended to the stream, the entry includes the `offset`, `partition`, and `timestamp`.
	// If a message failed to be appended to the stream, the entry includes the `error` and `errorMessage`.
	Entries []PutMessagesResultEntry `mandatory:"false" json:"entries"`
}

func (m PutMessagesResult) String() string {
	return common.PointerString(m)
}
