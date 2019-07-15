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

// PutMessagesDetailsEntry Object that represents a message to emit to a stream.
type PutMessagesDetailsEntry struct {

	// The message, expressed as a byte array up to 1 MiB in size.
	Value []byte `mandatory:"true" json:"value"`

	// The key of the message, expressed as a byte array up to 256 bytes in size. Messages with the same key are stored in the same partition.
	Key []byte `mandatory:"false" json:"key"`
}

func (m PutMessagesDetailsEntry) String() string {
	return common.PointerString(m)
}
