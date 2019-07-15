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

// Message A message in a stream.
type Message struct {

	// The name of the stream that the message belongs to.
	Stream *string `mandatory:"true" json:"stream"`

	// The ID of the partition where the message is stored.
	Partition *string `mandatory:"true" json:"partition"`

	// The key associated with the message, expressed as a byte array.
	Key []byte `mandatory:"true" json:"key"`

	// The value associated with the message, expressed as a byte array.
	Value []byte `mandatory:"true" json:"value"`

	// The offset of the message, which uniquely identifies it within the partition.
	Offset *int64 `mandatory:"true" json:"offset"`

	// The timestamp indicating when the server appended the message to the stream.
	Timestamp *common.SDKTime `mandatory:"true" json:"timestamp"`
}

func (m Message) String() string {
	return common.PointerString(m)
}
