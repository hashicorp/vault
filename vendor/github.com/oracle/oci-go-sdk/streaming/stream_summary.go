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

// StreamSummary Summary representation of a stream.
type StreamSummary struct {

	// The name of the stream.
	// Example: `TelemetryEvents`
	Name *string `mandatory:"true" json:"name"`

	// The OCID of the stream.
	Id *string `mandatory:"true" json:"id"`

	// The number of partitions in the stream.
	Partitions *int `mandatory:"true" json:"partitions"`

	// The OCID of the compartment that contains the stream.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The current state of the stream.
	LifecycleState StreamSummaryLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The date and time the stream was created, expressed in RFC 3339 (https://tools.ietf.org/rfc/rfc3339) timestamp format.
	// Example: `2018-04-20T00:00:07.405Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The endpoint to use when creating the StreamClient to consume or publish messages in the stream.
	MessagesEndpoint *string `mandatory:"true" json:"messagesEndpoint"`

	// Free-form tags for this resource. Each tag is a simple key-value pair that is applied with no predefined name, type, or namespace. Exists for cross-compatibility only.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m StreamSummary) String() string {
	return common.PointerString(m)
}

// StreamSummaryLifecycleStateEnum Enum with underlying type: string
type StreamSummaryLifecycleStateEnum string

// Set of constants representing the allowable values for StreamSummaryLifecycleStateEnum
const (
	StreamSummaryLifecycleStateCreating StreamSummaryLifecycleStateEnum = "CREATING"
	StreamSummaryLifecycleStateActive   StreamSummaryLifecycleStateEnum = "ACTIVE"
	StreamSummaryLifecycleStateDeleting StreamSummaryLifecycleStateEnum = "DELETING"
	StreamSummaryLifecycleStateDeleted  StreamSummaryLifecycleStateEnum = "DELETED"
	StreamSummaryLifecycleStateFailed   StreamSummaryLifecycleStateEnum = "FAILED"
)

var mappingStreamSummaryLifecycleState = map[string]StreamSummaryLifecycleStateEnum{
	"CREATING": StreamSummaryLifecycleStateCreating,
	"ACTIVE":   StreamSummaryLifecycleStateActive,
	"DELETING": StreamSummaryLifecycleStateDeleting,
	"DELETED":  StreamSummaryLifecycleStateDeleted,
	"FAILED":   StreamSummaryLifecycleStateFailed,
}

// GetStreamSummaryLifecycleStateEnumValues Enumerates the set of values for StreamSummaryLifecycleStateEnum
func GetStreamSummaryLifecycleStateEnumValues() []StreamSummaryLifecycleStateEnum {
	values := make([]StreamSummaryLifecycleStateEnum, 0)
	for _, v := range mappingStreamSummaryLifecycleState {
		values = append(values, v)
	}
	return values
}
