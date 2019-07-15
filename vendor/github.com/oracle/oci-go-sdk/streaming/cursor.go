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

// Cursor A cursor that indicates the position in the stream from which you want to begin consuming messages and which is required by the GetMessages operation.
type Cursor struct {

	// The cursor to pass to the `GetMessages` operation.
	Value *string `mandatory:"true" json:"value"`
}

func (m Cursor) String() string {
	return common.PointerString(m)
}
