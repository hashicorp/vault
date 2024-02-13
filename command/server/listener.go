// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	_ "crypto/sha512"
	"fmt"
	"io"
	"net"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	"github.com/hashicorp/vault/helper/proxyutil"
	"github.com/hashicorp/vault/internalshared/configutil"
)

// ListenerFactory is the factory function to create a listener.
type ListenerFactory func(*configutil.Listener, io.Writer, cli.Ui) (net.Listener, map[string]string, reloadutil.ReloadFunc, error)

// BuiltinListeners is the list of built-in listener types.
var BuiltinListeners = map[configutil.ListenerType]ListenerFactory{
	"tcp":  tcpListenerFactory,
	"unix": unixListenerFactory,
}

// NewListener creates a new listener of the given type with the given
// configuration. The type is looked up in the BuiltinListeners map.
func NewListener(l *configutil.Listener, logger io.Writer, ui cli.Ui) (net.Listener, map[string]string, reloadutil.ReloadFunc, error) {
	f, ok := BuiltinListeners[l.Type]
	if !ok {
		return nil, nil, nil, fmt.Errorf("unknown listener type: %q", l.Type)
	}

	return f(l, logger, ui)
}

func listenerWrapProxy(ln net.Listener, l *configutil.Listener) (net.Listener, error) {
	behavior := l.ProxyProtocolBehavior
	if behavior == "" {
		return ln, nil
	}

	proxyProtoConfig := &proxyutil.ProxyProtoConfig{
		Behavior:        behavior,
		AuthorizedAddrs: l.ProxyProtocolAuthorizedAddrs,
	}

	newLn, err := proxyutil.WrapInProxyProto(ln, proxyProtoConfig)
	if err != nil {
		return nil, fmt.Errorf("failed configuring PROXY protocol wrapper: %w", err)
	}

	return newLn, nil
}
