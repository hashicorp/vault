package proxyutil

import (
	"fmt"
	"net"
	"sync"

	proxyproto "github.com/armon/go-proxyproto"
	"github.com/hashicorp/errwrap"
	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/helper/strutil"
)

// ProxyProtoConfig contains configuration for the PROXY protocol
type ProxyProtoConfig struct {
	sync.RWMutex
	AllowedAddrs []*sockaddr.SockAddrMarshaler `json:"allowed_addrs"`
}

// WrapInProxyProto wraps the given listener in the PROXY protocol. If behavior
// is "use_if_authorized" or "deny_if_unauthorized" it also configures a
// SourceCheck based on the given ProxyProtoConfig. In an error case it returns
// the original listener and the error.
func WrapInProxyProto(listener net.Listener, config *ProxyProtoConfig, behavior, allowedAddrs string) (net.Listener, error) {
	config.Lock()
	defer config.Unlock()

	switch behavior {
	case "use_always", "use_if_authorized", "deny_if_unauthorized":
	default:
		return listener, fmt.Errorf("unknown behavior type for proxy proto config")
	}

	addrSlice := strutil.ParseArbitraryStringSlice(allowedAddrs, ",")
	if len(addrSlice) == 0 {
		return listener, fmt.Errorf("unable to parse addresses from %q", allowedAddrs)
	}

	config.AllowedAddrs = make([]*sockaddr.SockAddrMarshaler, 0, len(addrSlice))
	for _, addr := range addrSlice {
		sa, err := sockaddr.NewSockAddr(addr)
		if err != nil {
			return listener, errwrap.Wrapf("error parsing allowed address: {{err}}", err)
		}
		config.AllowedAddrs = append(config.AllowedAddrs, &sockaddr.SockAddrMarshaler{
			SockAddr: sa,
		})
	}

	newLn := &proxyproto.Listener{
		Listener: listener,
	}

	if behavior == "use_if_authorized" || behavior == "deny_if_unauthorized" {
		newLn.SourceCheck = func(addr net.Addr) (bool, error) {
			config.RLock()
			defer config.RUnlock()

			sa, err := sockaddr.NewSockAddr(addr.String())
			if err != nil {
				return false, errwrap.Wrapf("error parsing remote address: {{err}}", err)
			}

			for _, allowedAddr := range config.AllowedAddrs {
				if allowedAddr.Contains(sa) {
					return true, nil
				}
			}

			if behavior == "use_if_authorized" {
				return false, nil
			}

			return false, proxyproto.ErrInvalidUpstream
		}
	}

	return newLn, nil
}
