// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package link

import "context"

// Link offers functionality for linked HCP resources.
type Link interface {
	// Start will expose Link functionality to the control-plane. The SCADAProvider
	// used by Link will need to be started separately.
	Start() error

	// Stop will stop exposing Link functionality to the control-plane.
	Stop() error

	// ReportNodeStatus will get the most recent node status information from
	// the configured node status reporter and push it to HCP.
	//
	// This function only needs to be invoked in situations where it is
	// important that the node status is reported right away. HCP will regularly
	// poll for node status information.
	ReportNodeStatus(ctx context.Context) error
}
