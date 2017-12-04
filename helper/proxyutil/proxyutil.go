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
	Behavior        string
	AuthorizedAddrs []*sockaddr.SockAddrMarshaler `json:"authorized_addrs"`
}

func (p *ProxyProtoConfig) SetAuthorizedAddrs(addrs interface{}) error {
	p.AuthorizedAddrs = make([]*sockaddr.SockAddrMarshaler, 0)
	stringAddrs := make([]string, 0)

	switch addrs.(type) {
	case string:
		stringAddrs = strutil.ParseArbitraryStringSlice(addrs.(string), ",")
		if len(stringAddrs) == 0 {
			return fmt.Errorf("unable to parse addresses from %v", addrs)
		}

	case []string:
		stringAddrs = addrs.([]string)

	case []interface{}:
		for _, v := range addrs.([]interface{}) {
			stringAddr, ok := v.(string)
			if !ok {
				return fmt.Errorf("error parsing %q as string")
			}
			stringAddrs = append(stringAddrs, stringAddr)
		}

	default:
		return fmt.Errorf("unknown address input type %T", addrs)
	}

	for _, addr := range stringAddrs {
		sa, err := sockaddr.NewSockAddr(addr)
		if err != nil {
			return errwrap.Wrapf("error parsing authorized address: {{err}}", err)
		}
		p.AuthorizedAddrs = append(p.AuthorizedAddrs, &sockaddr.SockAddrMarshaler{
			SockAddr: sa,
		})
	}

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
			Listener: listener,
		}

	case "allow_authorized", "deny_unauthorized":
		newLn = &proxyproto.Listener{
			Listener: listener,
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
