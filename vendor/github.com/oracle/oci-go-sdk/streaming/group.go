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

// Group Represents the current state of a consumer group, including partition reservations and committed offsets.
type Group struct {

	// The streamId for which the group exists.
	StreamId *string `mandatory:"false" json:"streamId"`

	// The name of the consumer group.
	GroupName *string `mandatory:"false" json:"groupName"`

	// An array of the partition reservations of a group.
	Reservations []PartitionReservation `mandatory:"false" json:"reservations"`
}

func (m Group) String() string {
	return common.PointerString(m)
}
