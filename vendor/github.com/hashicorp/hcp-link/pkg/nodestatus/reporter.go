// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package nodestatus contains the interface that needs to be implements for a self-managed resource to report
// inforamtion to HashiCorp Cloud Platform (HCP).
package nodestatus

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// Reporter is an interface that needs to be implemented to provide product
// specific status information about a node.
type Reporter interface {
	// GetNodeStatus will return the node's current status. The information will
	// be transmitted to HCP.
	GetNodeStatus(context.Context) (NodeStatus, error)
}

// NodeStatus contains product specific information about a node's status.
type NodeStatus struct {
	// StatusVersion is the version of the status message format. To ensure
	// that the version is not omitted by accident the initial version is 1.
	StatusVersion uint32

	// Status contains product specific status information as a proto message.
	// The status message will be wrapped in a protobuf Any message.
	Status proto.Message
}
