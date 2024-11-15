// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/go-hclog"
)

type hijackFunc func(net.Conn)

// providerEndpoint is used to implement the Provider.* RPC endpoints
// as part of the provider.
type providerEndpoint struct {
	p      *Provider
	hijack hijackFunc
}

// hijacked is used to check if the connection has been hijacked.
func (pe *providerEndpoint) hijacked() bool {
	return pe.hijack != nil
}

// getHijack returns the hijack function.
func (pe *providerEndpoint) getHijack() hijackFunc {
	return pe.hijack
}

// setHijack is used to take over the yamux stream for Provider.Connect.
func (pe *providerEndpoint) setHijack(cb hijackFunc) {
	pe.hijack = cb
}

// Connect is invoked by the broker to connect to a capability.
func (pe *providerEndpoint) Connect(args *ConnectRequest, resp *ConnectResponse) error {
	pe.p.logger.Debug("connect requested", "capability", args.Capability)

	// Handle potential flash
	if args.Severity != "" && args.Message != "" {
		switch hclog.LevelFromString(args.Severity) {
		case hclog.Trace:
			pe.p.logger.Trace("connect message", "msg", args.Message)
		case hclog.Debug:
			pe.p.logger.Debug("connect message", "msg", args.Message)
		case hclog.Info:
			pe.p.logger.Info("connect message", "msg", args.Message)
		case hclog.Warn:
			pe.p.logger.Warn("connect message", "msg", args.Message)
		}
	}

	// Look for the handler
	pe.p.handlersLock.RLock()
	handler := pe.p.handlers[args.Capability].provider
	pe.p.handlersLock.RUnlock()
	if handler == nil {
		pe.p.logger.Warn("requested capability not available", "capability", args.Capability)
		return fmt.Errorf("invalid capability")
	}

	// Hijack the connection
	pe.setHijack(func(a net.Conn) {
		if err := handler(args.Capability, args.Meta, a); err != nil {
			pe.p.logger.Error("handler errored", "capability", args.Capability, "error", err)
		}
	})
	resp.Success = true
	return nil
}

// Disconnect is invoked by the broker to ask us to backoff.
func (pe *providerEndpoint) Disconnect(args *DisconnectRequest, resp *DisconnectResponse) error {
	if args.Reason == "" {
		args.Reason = "<no reason provided>"
	}
	pe.p.logger.Info("disconnect requested",
		"retry", !args.NoRetry,
		"backoff", args.Backoff,
		"reason", args.Reason)

	// Use the backoff information
	pe.p.backoffLock.Lock()
	pe.p.noRetry = args.NoRetry
	pe.p.backoff = args.Backoff
	pe.p.backoffLock.Unlock()

	// Force the disconnect
	time.AfterFunc(disconnectDelay, func() {
		pe.p.action(actionDisconnect)
	})
	return nil
}
