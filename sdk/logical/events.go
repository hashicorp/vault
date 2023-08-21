// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hashicorp/go-uuid"
	"google.golang.org/protobuf/types/known/structpb"
)

// common event metadata keys
const (
	EventMetadataSecretPath        = "secret_path"
	EventMetadataOperation         = "operation"
	EventMetadataSourcePluginMount = "source_plugin_mount"
)

var ErrOddMetadataStrings = errors.New("odd number of arguments to metadataPairs")

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
	if len(metadataPairs)%2 != 0 {
		return ErrOddMetadataStrings
	}
	for i := 0; i < len(metadataPairs); i += 2 {
		metadata[metadataPairs[i]] = metadataPairs[i+1]
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
