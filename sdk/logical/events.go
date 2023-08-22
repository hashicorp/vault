// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/go-uuid"
	"google.golang.org/protobuf/types/known/structpb"
)

// common event metadata keys
const (
	// EventMetadataApiPath is used in event metadata to show the API path that can be used to fetch any underlying
	// data. For example, the KV plugin would set this to `data/mysecret`. The event system will automatically prepend
	// the plugin mount to this path, if present, so it would be come `secret/data/mysecret`, for example.
	// If this is an auth plugin event, this will additionally be prepended with `auth/`.
	EventMetadataApiPath = "api_path"
	// EventMetadataOperation is used in event metadata to express what operation was performed that generated the
	// event, e.g., `read` or `write`.
	EventMetadataOperation = "operation"

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
	metadata := map[string]string{}
	if len(metadataPairs) >= 2 {
		for i := 0; i < len(metadataPairs)-1; i += 2 {
			metadata[metadataPairs[i]] = metadataPairs[i+1]
		}
	}
	if len(metadataPairs)%2 != 0 {
		metadata[extraMetadataArgument] = metadataPairs[len(metadataPairs)-1]
	}
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	ev.Metadata = &structpb.Struct{}
	if err := ev.Metadata.UnmarshalJSON(metadataBytes); err != nil {
		return err
	}
	return sender.SendEvent(ctx, EventType(eventType), ev)
}
