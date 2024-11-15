// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "time"

// ConnectRequest holds parameters for the broker RPC Connect call to the provider.
type ConnectRequest struct {
	Capability string
	Meta       map[string]string

	Severity string
	Message  string
}

// ConnectResponse is the response to a Connect RPC call.
type ConnectResponse struct {
	Success bool
}

// DisconnectRequest holds parameters for the broker RPC Disconnect call to the provider.
type DisconnectRequest struct {
	NoRetry bool          // Should the client retry
	Backoff time.Duration // Minimum backoff
	Reason  string
}

// DisconnectResponse is the response to a Disconnect RPC call.
type DisconnectResponse struct {
}
