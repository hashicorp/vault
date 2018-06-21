package proxyutil

import (
	"fmt"
	"net"
	"sync"
	"time"

	proxyproto "github.com/armon/go-proxyproto"
	"github.com/hashicorp/errwrap"
	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/helper/parseutil"
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
			Listener:           listener,
			ProxyHeaderTimeout: 10 * time.Second,
		}

	case "allow_authorized", "deny_unauthorized":
		newLn = &proxyproto.Listener{
			Listener:           listener,
			ProxyHeaderTimeout: 10 * time.Second,
			SourceCheck: func(addr net.Addr) (bool, error) {
				config.RLock()
				defer config.RUnlock()

				sa, err := sockaddr.NewSockAddr(addr.String())
				if err != nil {
					return false, errwrap.Wrapf("error parsing remote address: {{err}}", err)
				}

				for _, authorizedAddr := range config.AuthorizedAddrs {
					if authorizedAddr.Contains(sa) {
						return true, nil
					}
				}

				if config.Behavior == "allow_authorized" {
					return false, nil
				}

				return false, proxyproto.ErrInvalidUpstream
			},
		}
	default:
		return listener, fmt.Errorf("unknown behavior type for proxy proto config")
	}

	return newLn, nil
}
