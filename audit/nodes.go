// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/logical"
)

// processManual will attempt to create an (audit) event with the specified data
// and manually iterate over the supplied nodes calling Process on each until the
// event is nil (which indicates the pipeline has completed).
// Order of IDs in the NodeID slice determines the order they are processed.
// (Audit) Event will be of RequestType (as opposed to ResponseType).
// The last node must be a filter node (eventlogger.NodeTypeFilter) or
// sink node (eventlogger.NodeTypeSink).
func processManual(ctx context.Context, data *logical.LogInput, ids []eventlogger.NodeID, nodes map[eventlogger.NodeID]eventlogger.Node) error {
	switch {
	case data == nil:
		return errors.New("data cannot be nil")
	case len(ids) < 2:
		return errors.New("minimum of 2 ids are required")
	case nodes == nil:
		return errors.New("nodes cannot be nil")
	case len(nodes) == 0:
		return errors.New("nodes are required")
	}

	// Create an audit event.
	a, err := newEvent(RequestType)
	if err != nil {
		return err
	}

	// Insert the data into the audit event.
	a.Data = data

	// Create an eventlogger event with the audit event as the payload.
	e := &eventlogger.Event{
		Type:      event.AuditType.AsEventType(),
		CreatedAt: time.Now(),
		Formatted: make(map[string][]byte),
		Payload:   a,
	}

	var lastSeen eventlogger.NodeType

	// Process nodes in order, updating the event with the result.
	// This means we *should* do:
	// 1. filter (optional if configured)
	// 2. formatter (temporary)
	// 3. sink
	for _, id := range ids {
		// If the event is nil, we've completed processing the pipeline (hopefully
		// by either a filter node or a sink node).
		if e == nil {
			break
		}
		node, ok := nodes[id]
		if !ok {
			return fmt.Errorf("node not found: %v", id)
		}

		switch node.Type() {
		case eventlogger.NodeTypeFormatter:
			// Use a temporary formatter node  which doesn't persist its salt anywhere.
			if formatNode, ok := node.(*entryFormatter); ok && formatNode != nil {
				e, err = newTemporaryEntryFormatter(formatNode).Process(ctx, e)
			}
		default:
			e, err = node.Process(ctx, e)
		}

		if err != nil {
			return err
		}

		// Track the last node we have processed, as we should end with a filter or sink.
		lastSeen = node.Type()
	}

	switch lastSeen {
	case eventlogger.NodeTypeSink, eventlogger.NodeTypeFilter:
	default:
		return errors.New("last node must be a filter or sink")
	}

	return nil
}
