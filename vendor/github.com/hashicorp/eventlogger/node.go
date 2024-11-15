// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventlogger

import (
	"context"
	"fmt"
)

// NodeType defines the possible Node type's in the system.
type NodeType int

const (
	_ NodeType = iota
	NodeTypeFilter
	NodeTypeFormatter
	NodeTypeSink
	NodeTypeFormatterFilter // A node that formats and then filters the events based on the new format.
)

// A Node in a graph
type Node interface {
	// Process does something with the Event: filter, redaction,
	// marshalling, persisting.
	Process(ctx context.Context, e *Event) (*Event, error)
	// Reopen is used to re-read any config stored externally
	// and to close and reopen files, e.g. for log rotation.
	Reopen() error
	// Type describes the type of the node.  This is mostly just used to
	// validate that pipelines are sensibly arranged, e.g. ending with a sink.
	Type() NodeType
}

// A NodeController is used by a Broker to attempt additional control of a given node.
// For instance, when a Node supports being closed via the Closer interface.
type NodeController struct {
	n Node
}

// NewNodeController creates a new NodeController for a given Node. The Node
// should be the original value registered with the broker, or have an Unwrap
// method returning the original Node (see NodeUnwrapper interface).
//
// If the Node implements any of the following methods, the NodeController will
// call them as appropriate/needed:
//
//	Close() error
func NewNodeController(n Node) *NodeController {
	// intentionally not checking the Node for nil.. the caller must ensure it's
	// valid and the docs make that clear.
	return &NodeController{n}
}

// NodeUnwrapper will unwrap a node, returning the original value (see
// NewNodeController docs)
type NodeUnwrapper interface {
	Unwrap() Node
}

// Closer will close without error
type Closer interface {
	Close(ctx context.Context) error
}

// Close the Node if it implements the Closer interface, and if required use the
// NodeUnwrapper interface to unwrap it before closing it.
func (nc *NodeController) Close(ctx context.Context) error {
	n := nc.n
	for {
		switch t := n.(type) {
		case Closer:
			return t.Close(ctx)
		case NodeUnwrapper:
			n = t.Unwrap()
		default:
			return nil
		}
	}
}

type linkedNode struct {
	node   Node
	nodeID NodeID
	next   []*linkedNode
}

// linkNodes is a convenience function that connects Nodes together into a linked list.
func linkNodes(nodes []Node, ids []NodeID) (*linkedNode, error) {
	switch {
	case len(nodes) == 0:
		return nil, fmt.Errorf("no nodes given")
	case len(ids) == 0:
		return nil, fmt.Errorf("no IDs given")
	case len(nodes) != len(ids):
		return nil, fmt.Errorf("number of nodes does not match number of IDs")
	}

	root := &linkedNode{node: nodes[0], nodeID: ids[0]}
	cur := root

	for i, n := range nodes[1:] {
		next := &linkedNode{node: n, nodeID: ids[i+1]}
		cur.next = []*linkedNode{next}
		cur = next
	}

	return root, nil
}

// linkNodesAndSinks is a convenience function that connects
// the inner Nodes together into a linked list.  Then it appends the sinks
// to the end as a set of fan-out leaves.
func linkNodesAndSinks(inner, sinks []Node, nodeIDs, sinkIDs []NodeID) (*linkedNode, error) {
	root, err := linkNodes(inner, nodeIDs)
	if err != nil {
		return nil, err
	}

	// This is inefficient but since it's only used in setup we don't care:
	cur := root
	for cur.next != nil {
		cur = cur.next[0]
	}

	for i, s := range sinks {
		cur.next = append(cur.next, &linkedNode{node: s, nodeID: sinkIDs[i]})
	}

	return root, nil
}

// flatten will attempt to visit every linked node and flatten the overall set of node IDs.
func (l *linkedNode) flatten() map[NodeID]struct{} {
	stack := []*linkedNode{l}
	flattened := make(map[NodeID]struct{})

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Skip already flattened nodes
		if _, ok := flattened[node.nodeID]; ok {
			continue
		}

		flattened[node.nodeID] = struct{}{}

		for _, child := range node.next {
			stack = append(stack, child)
		}
	}

	return flattened
}
