// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package event

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/eventlogger"
)

var _ eventlogger.Node = (*MetricsCounter)(nil)

// MetricsCounter offers a way for nodes to emit metrics which increment a label by 1.
type MetricsCounter struct {
	Name    string
	Node    eventlogger.Node
	labeler Labeler
}

// Labeler provides a way to inject the logic required to determine labels based
// on the state of the eventlogger.Event being returned and the error resulting
// from processing the by the underlying eventlogger.Node.
type Labeler interface {
	Labels(*eventlogger.Event, error) []string
}

// NewMetricsCounter should be used to create the MetricsCounter.
func NewMetricsCounter(name string, node eventlogger.Node, labeler Labeler) (*MetricsCounter, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name is required: %w", ErrInvalidParameter)
	}

	if node == nil || reflect.ValueOf(node).IsNil() {
		return nil, fmt.Errorf("node is required: %w", ErrInvalidParameter)
	}

	if labeler == nil || reflect.ValueOf(labeler).IsNil() {
		return nil, fmt.Errorf("labeler is required: %w", ErrInvalidParameter)
	}

	return &MetricsCounter{
		Name:    name,
		Node:    node,
		labeler: labeler,
	}, nil
}

// Process will process the event using the underlying eventlogger.Node, and then
// use the configured Labeler to provide a label which is used to increment a metric by 1.
func (m MetricsCounter) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	// NOTE: We don't provide an 'op' here, as we're just wrapping the underlying node.
	var err error

	// Process the node first
	e, err = m.Node.Process(ctx, e)

	// Provide the results to the Labeler.
	metrics.IncrCounter(m.labeler.Labels(e, err), 1)

	return e, err
}

// Reopen attempts to reopen the underlying eventlogger.Node.
func (m MetricsCounter) Reopen() error {
	return m.Node.Reopen()
}

// Type returns the type for the underlying eventlogger.Node.
func (m MetricsCounter) Type() eventlogger.NodeType {
	return m.Node.Type()
}
