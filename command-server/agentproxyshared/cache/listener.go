// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/vault/command-server/server"

	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
)

type ListenerBundle struct {
	Listener      net.Listener
	TLSConfig     *tls.Config
	TLSReloadFunc reloadutil.ReloadFunc
}

func StartListener(lnConfig *configutil.Listener) (*ListenerBundle, error) {
	addr := lnConfig.Address

	var ln net.Listener
	var err error
	switch lnConfig.Type {
	case "tcp":
		if addr == "" {
			addr = "127.0.0.1:8200"
		}

		bindProto := "tcp"
		// If they've passed 0.0.0.0, we only want to bind on IPv4
		// rather than golang's dual stack default
		if strings.HasPrefix(addr, "0.0.0.0:") {
			bindProto = "tcp4"
		}

		ln, err = net.Listen(bindProto, addr)
		if err != nil {
			return nil, err
		}
		ln = &server.TCPKeepAliveListener{ln.(*net.TCPListener)}

	case "unix":
		var uConfig *listenerutil.UnixSocketsConfig
		if lnConfig.SocketMode != "" &&
			lnConfig.SocketUser != "" &&
			lnConfig.SocketGroup != "" {
			uConfig = &listenerutil.UnixSocketsConfig{
				Mode:  lnConfig.SocketMode,
				User:  lnConfig.SocketUser,
				Group: lnConfig.SocketGroup,
			}
		}
		ln, err = listenerutil.UnixSocketListener(addr, uConfig)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("invalid listener type: %q", lnConfig.Type)
	}

	props := map[string]string{"addr": ln.Addr().String()}
	tlsConf, reloadFunc, err := listenerutil.TLSConfig(lnConfig, props, nil)
	if err != nil {
		return nil, err
	}
	if tlsConf != nil {
		ln = tls.NewListener(ln, tlsConf)
	}

	cfg := &ListenerBundle{
		Listener:      ln,
		TLSConfig:     tlsConf,
		TLSReloadFunc: reloadFunc,
	}

	return cfg, nil
}
