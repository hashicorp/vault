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

// Stream Detailed representation of a stream, including all its partitions.
type Stream struct {

	// The name of the stream. Avoid entering confidential information.
	// Example: `TelemetryEvents`
	Name *string `mandatory:"true" json:"name"`

	// The OCID of the stream.
	Id *string `mandatory:"true" json:"id"`

	// The number of partitions in the stream.
	Partitions *int `mandatory:"true" json:"partitions"`

	// The retention period of the stream, in hours. This property is read-only.
	RetentionInHours *int `mandatory:"true" json:"retentionInHours"`

	// The OCID of the stream.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The current state of the stream.
	LifecycleState StreamLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the stream was created, expressed in in RFC 3339 (https://tools.ietf.org/rfc/rfc3339) timestamp format.
	// Example: `2018-04-20T00:00:07.405Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The endpoint to use when creating the StreamClient to consume or publish messages in the stream.
	MessagesEndpoint *string `mandatory:"true" json:"messagesEndpoint"`

	// Any additional details about the current state of the stream.
	LifecycleStateDetails *string `mandatory:"false" json:"lifecycleStateDetails"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace. Exists for cross-compatibility only.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}'
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m Stream) String() string {
	return common.PointerString(m)
}

// StreamLifecycleStateEnum Enum with underlying type: string
type StreamLifecycleStateEnum string

// Set of constants representing the allowable values for StreamLifecycleStateEnum
const (
	StreamLifecycleStateCreating StreamLifecycleStateEnum = "CREATING"
	StreamLifecycleStateActive   StreamLifecycleStateEnum = "ACTIVE"
	StreamLifecycleStateDeleting StreamLifecycleStateEnum = "DELETING"
	StreamLifecycleStateDeleted  StreamLifecycleStateEnum = "DELETED"
	StreamLifecycleStateFailed   StreamLifecycleStateEnum = "FAILED"
)

var mappingStreamLifecycleState = map[string]StreamLifecycleStateEnum{
	"CREATING": StreamLifecycleStateCreating,
	"ACTIVE":   StreamLifecycleStateActive,
	"DELETING": StreamLifecycleStateDeleting,
	"DELETED":  StreamLifecycleStateDeleted,
	"FAILED":   StreamLifecycleStateFailed,
}

// GetStreamLifecycleStateEnumValues Enumerates the set of values for StreamLifecycleStateEnum
func GetStreamLifecycleStateEnumValues() []StreamLifecycleStateEnum {
	values := make([]StreamLifecycleStateEnum, 0)
	for _, v := range mappingStreamLifecycleState {
		values = append(values, v)
	}
	return values
}
