// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"fmt"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-uuid"
)

// EventType represents the event's type
type EventType string

const (
	AuditType EventType = "audit" // AuditType represents audit events
)

// Validate ensures that EventType is one of the set of allowed event types.
func (t EventType) Validate() error {
	const op = "event.(EventType).Validate"
	switch t {
	case AuditType:
		return nil
	default:
		return fmt.Errorf("%s: '%s' is not a valid event type: %w", op, t, ErrInvalidParameter)
	}
}

// GenerateNodeID generates a new UUID that it casts to the eventlogger.NodeID
// type.
func GenerateNodeID() (eventlogger.NodeID, error) {
	id, err := uuid.GenerateUUID()

	return eventlogger.NodeID(id), err
}

// String returns the string version of an EventType.
func (t EventType) String() string {
	return string(t)
}
