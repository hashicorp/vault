// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventlogger

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// graph
type graph struct {

	// roots maps PipelineIDs to pipelineRegistrations.
	// A registeredPipeline includes the root Node for a pipeline.
	roots graphMap

	// successThreshold specifies how many pipelines must successfully process
	// an event for Process to not return an error.  This means that a filter
	// could of course filter an event before it reaches the pipeline's sink,
	// but it would still count as success when it comes to meeting this threshold
	successThreshold int

	// successThresholdSinks specifies how many sinks must successfully process
	// an event for Process to not return an error.
	successThresholdSinks int
}

// Process the Event by routing it through all of the graph's nodes,
// starting with the root node.
func (g *graph) process(ctx context.Context, e *Event) (Status, error) {
	statusChan := make(chan Status)
	var wg sync.WaitGroup
	go func() {
		g.roots.Range(func(_ PipelineID, pipeline *registeredPipeline) bool {
			select {
			// Don't continue to start root nodes if our context is already done.
			// We would just process the node and then drop the status, and no
			// other linked nodes would be processed.
			case <-ctx.Done():
				return false
			default:
			}

			wg.Add(1)
			g.doProcess(ctx, pipeline.rootNode, e, statusChan, &wg)
			return true
		})
		wg.Wait()
		close(statusChan)
	}()
	var status Status
	var done bool
	for !done {
		select {
		case <-ctx.Done():
			done = true
		case s, ok := <-statusChan:
			if ok {
				status.Warnings = append(status.Warnings, s.Warnings...)
				status.complete = append(status.complete, s.complete...)
				status.completeSinks = append(status.completeSinks, s.completeSinks...)
			} else {
				done = true
			}
		}
	}
	return status, status.getError(ctx.Err(), g.successThreshold, g.successThresholdSinks)
}

// Recursively process every node in the graph.
//
// # No Status is sent when a request is cancelled by the context
//
// Status will be sent when we stop processing nodes, which can happen if:
//   - a node.Process(...) returns an error, and Status.complete is empty
//   - a node.Process(...) filters an event, and Status.complete contains the
//     filter node's ID
//   - the final node in a pipeline (a sink) finishes, and Status.complete contains
//     the sink node's ID
func (g *graph) doProcess(ctx context.Context, node *linkedNode, e *Event, statusChan chan Status, wg *sync.WaitGroup) {
	defer wg.Done()

	// Process the current Node
	e, err := node.node.Process(ctx, e)
	if err != nil {
		select {
		case <-ctx.Done():
		case statusChan <- Status{Warnings: []error{err}}:
		}
		return
	}

	completeStatus := Status{complete: []NodeID{node.nodeID}}
	if node.node.Type() == NodeTypeSink {
		completeStatus.completeSinks = []NodeID{node.nodeID}
	}

	// If the Event is nil, it has been filtered out and we are done.
	if e == nil {
		select {
		case <-ctx.Done():
		case statusChan <- completeStatus:
		}
		return
	}

	// Process any child nodes.  This is depth-first.
	if len(node.next) != 0 {
		// If the new Event is nil, it has been filtered out and we are done.
		if e == nil {
			statusChan <- Status{}
			return
		}

		for _, child := range node.next {
			wg.Add(1)
			go g.doProcess(ctx, child, e, statusChan, wg)
		}
	} else {
		select {
		case <-ctx.Done():
		case statusChan <- completeStatus:
		}
	}
}

func (g *graph) reopen(ctx context.Context) error {
	var errors *multierror.Error

	g.roots.Range(func(_ PipelineID, pipeline *registeredPipeline) bool {
		err := g.doReopen(ctx, pipeline.rootNode)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
		return true
	})

	return errors.ErrorOrNil()
}

// Recursively reopen every node in the graph.
func (g *graph) doReopen(ctx context.Context, node *linkedNode) error {
	// Process the current Node
	err := node.node.Reopen()
	if err != nil {
		return err
	}

	// Process any child nodes.  This is depth-first.
	for _, child := range node.next {

		err = g.doReopen(ctx, child)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *graph) validate() error {
	var errors *multierror.Error

	g.roots.Range(func(_ PipelineID, pipeline *registeredPipeline) bool {
		err := g.doValidate(nil, pipeline.rootNode)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
		return true
	})

	return errors.ErrorOrNil()
}

func (g *graph) doValidate(parent, node *linkedNode) error {
	isInner := len(node.next) > 0

	switch {
	case len(node.next) == 0 && node.node.Type() != NodeTypeSink:
		return fmt.Errorf("non-sink node has no children")
	case !isInner && parent == nil:
		return fmt.Errorf("sink node at root")
	case !isInner && (parent.node.Type() != NodeTypeFormatter && parent.node.Type() != NodeTypeFormatterFilter):
		return fmt.Errorf("sink node without preceding formatter or formatter filter")
	case !isInner:
		return nil
	}

	// Process any child nodes.  This is depth-first.
	for _, child := range node.next {
		err := g.doValidate(node, child)
		if err != nil {
			return err
		}
	}

	return nil
}
