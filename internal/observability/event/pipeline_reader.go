// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import "github.com/hashicorp/eventlogger"

// PipelineReader surfaces information required for pipeline registration.
type PipelineReader interface {
	// EventType should return the event type to be used for pipeline registration.
	EventType() eventlogger.EventType

	// HasFiltering should determine if filter nodes are used by this pipeline.
	HasFiltering() bool

	// Name for the pipeline which should be used for the eventlogger.PipelineID.
	Name() string

	// Nodes should return the nodes which should be used by the framework to process events.
	Nodes() map[eventlogger.NodeID]eventlogger.Node

	// NodeIDs should return the IDs of the nodes, in the order they are required.
	NodeIDs() []eventlogger.NodeID
}
