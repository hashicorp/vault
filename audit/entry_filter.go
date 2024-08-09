// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-bexpr"
	nshelper "github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

var _ eventlogger.Node = (*entryFilter)(nil)

// entryFilter should be used to filter audit requests and responses which should
// make it to a sink.
type entryFilter struct {
	// the evaluator for the bexpr expression that should be applied by the node.
	evaluator *bexpr.Evaluator
}

// newEntryFilter should be used to create an entryFilter node.
// The filter supplied should be in bexpr format and reference fields from logical.LogInputBexpr.
func newEntryFilter(filter string) (*entryFilter, error) {
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return nil, fmt.Errorf("cannot create new audit filter with empty filter expression: %w", ErrExternalOptions)
	}

	eval, err := bexpr.CreateEvaluator(filter)
	if err != nil {
		return nil, fmt.Errorf("cannot create new audit filter: %w: %w", ErrExternalOptions, err)
	}

	// Validate the filter by attempting to evaluate it with an empty input.
	// This prevents users providing a filter with a field that would error during
	// matching, and block all auditable requests to Vault.
	li := logical.LogInputBexpr{}
	_, err = eval.Evaluate(li)
	if err != nil {
		return nil, fmt.Errorf("filter references an unsupported field: %s: %w", filter, ErrExternalOptions)
	}

	return &entryFilter{evaluator: eval}, nil
}

// Reopen is a no-op for the filter node.
func (*entryFilter) Reopen() error {
	return nil
}

// Type describes the type of this node (filter).
func (*entryFilter) Type() eventlogger.NodeType {
	return eventlogger.NodeTypeFilter
}

// Process will attempt to parse the incoming event data and decide whether it
// should be filtered or remain in the pipeline and passed to the next node.
func (f *entryFilter) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if e == nil {
		return nil, fmt.Errorf("event is nil: %w", ErrInvalidParameter)
	}

	a, ok := e.Payload.(*Event)
	if !ok {
		return nil, fmt.Errorf("cannot parse event payload: %w", ErrInvalidParameter)
	}

	// If we don't have data to process, then we're done.
	if a.Data == nil {
		return nil, nil
	}

	ns, err := nshelper.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot obtain namespace: %w", err)
	}

	datum := a.Data.BexprDatum(ns.Path)

	result, err := f.evaluator.Evaluate(datum)
	if err != nil {
		return nil, fmt.Errorf("unable to evaluate filter: %w", err)
	}

	if result {
		// Allow this event to carry on through the pipeline.
		return e, nil
	}

	// End process of this pipeline.
	return nil, nil
}
