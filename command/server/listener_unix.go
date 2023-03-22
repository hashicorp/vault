// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package server

import (
	"io"
	"net"

	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/mitchellh/cli"
)

func unixListenerFactory(l *configutil.Listener, _ io.Writer, ui cli.Ui) (net.Listener, map[string]string, reloadutil.ReloadFunc, error) {
	addr := l.Address
	if addr == "" {
		addr = "/run/vault.sock"
	}

	var cfg *listenerutil.UnixSocketsConfig
	if l.SocketMode != "" &&
		l.SocketUser != "" &&
		l.SocketGroup != "" {
		cfg = &listenerutil.UnixSocketsConfig{
			Mode:  l.SocketMode,
			User:  l.SocketUser,
			Group: l.SocketGroup,
		}
	}

	ln, err := listenerutil.UnixSocketListener(addr, cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	return ln, map[string]string{}, nil, nil
}
