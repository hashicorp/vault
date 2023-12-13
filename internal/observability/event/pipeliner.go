// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import "github.com/hashicorp/eventlogger"

// The Pipeliner interface surfaces information required for pipeline registration.
type Pipeliner interface {
	// EventType returns the event type to be used for registration.
	EventType() eventlogger.EventType

	// IsFilteringPipeline determines whether the pipeline uses filtering.
	IsFilteringPipeline() bool

	// Name for the pipeline which should be used for the eventlogger.PipelineID.
	Name() string

	// Nodes returns the nodes which should be used by the framework to process events.
	Nodes() map[eventlogger.NodeID]eventlogger.Node

	// NodeIDs returns the IDs of the nodes, in the order they are required.
	NodeIDs() []eventlogger.NodeID
}
