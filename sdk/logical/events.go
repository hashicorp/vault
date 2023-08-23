// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"google.golang.org/protobuf/types/known/structpb"
)

// common event metadata keys
const (
	// EventMetadataDataPath is used in event metadata to show the API path that can be used to fetch any underlying
	// data. For example, the KV plugin would set this to `data/mysecret`. The event system will automatically prepend
	// the plugin mount to this path, if present, so it would become `secret/data/mysecret`, for example.
	// If this is an auth plugin event, this will additionally be prepended with `auth/`.
	EventMetadataDataPath = "data_path"
	// EventMetadataOperation is used in event metadata to express what operation was performed that generated the
	// event, e.g., `read` or `write`.
	EventMetadataOperation = "operation"
	// EventMetadataModified is used in event metadata when the event attests that the underlying data has been modified
	// and might need to be re-fetched (at the EventMetadataDataPath).
	EventMetadataModified = "modified"

	extraMetadataArgument = "EXTRA_VALUE_AT_END"
)

// ID is an alias to GetId() for CloudEvents compatibility.
func (x *EventReceived) ID() string {
	return x.Event.GetId()
}

// NewEvent returns an event with a new, random EID.
func NewEvent() (*EventData, error) {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	return &EventData{
		Id: id,
	}, nil
}

// EventType represents a topic, and is a wrapper around eventlogger.EventType.
type EventType string

// EventSender sends events to the common event bus.
type EventSender interface {
	SendEvent(ctx context.Context, eventType EventType, event *EventData) error
}

// SendEvent is a convenience method for plugins events to an EventSender, converting the
// metadataPairs to the EventData structure.
func SendEvent(ctx context.Context, sender EventSender, eventType string, metadataPairs ...string) error {
	ev, err := NewEvent()
	if err != nil {
		return err
	}
	ev.Metadata = &structpb.Struct{Fields: make(map[string]*structpb.Value, (len(metadataPairs)+1)/2)}
	for i := 0; i < len(metadataPairs)-1; i += 2 {
		ev.Metadata.Fields[metadataPairs[i]] = structpb.NewStringValue(metadataPairs[i+1])
	}
	if len(metadataPairs)%2 != 0 {
		ev.Metadata.Fields[extraMetadataArgument] = structpb.NewStringValue(metadataPairs[len(metadataPairs)-1])
	}
	return sender.SendEvent(ctx, EventType(eventType), ev)
}
