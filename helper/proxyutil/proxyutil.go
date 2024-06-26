// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package proxyutil

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-sockaddr"
	"github.com/pires/go-proxyproto"
)

// ProxyProtoConfig contains configuration for the PROXY protocol
type ProxyProtoConfig struct {
	sync.RWMutex
	Behavior        string
	AuthorizedAddrs []*sockaddr.SockAddrMarshaler `json:"authorized_addrs"`
}

func (p *ProxyProtoConfig) SetAuthorizedAddrs(addrs interface{}) error {
	aa, err := parseutil.ParseAddrs(addrs)
	if err != nil {
		return err
	}

	p.AuthorizedAddrs = aa
	return nil
}

// WrapInProxyProto wraps the given listener in the PROXY protocol. If behavior
// is "use_if_authorized" or "deny_if_unauthorized" it also configures a
// SourceCheck based on the given ProxyProtoConfig. In an error case it returns
// the original listener and the error.
func WrapInProxyProto(listener net.Listener, config *ProxyProtoConfig) (net.Listener, error) {
	config.Lock()
	defer config.Unlock()

	var newLn *proxyproto.Listener

	switch config.Behavior {
	case "use_always":
		newLn = &proxyproto.Listener{
			Listener:          listener,
			ReadHeaderTimeout: 10 * time.Second,
		}

	case "allow_authorized", "deny_unauthorized":
		newLn = &proxyproto.Listener{
			Listener:          listener,
			ReadHeaderTimeout: 10 * time.Second,
			Policy: func(addr net.Addr) (proxyproto.Policy, error) {
				config.RLock()
				defer config.RUnlock()

				sa, err := sockaddr.NewSockAddr(addr.String())
				if err != nil {
					return proxyproto.REJECT, fmt.Errorf("error parsing remote address: %w", err)
				}

				for _, authorizedAddr := range config.AuthorizedAddrs {
					if authorizedAddr.Contains(sa) {
						return proxyproto.USE, nil
					}
				}

				if config.Behavior == "allow_authorized" {
					return proxyproto.IGNORE, nil
				}

				return proxyproto.REJECT, proxyproto.ErrInvalidUpstream
			},
		}
	default:
		return listener, fmt.Errorf("unknown behavior type for proxy proto config")
	}

	return newLn, nil
}
