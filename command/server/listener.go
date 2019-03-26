package server

import (
	"github.com/hashicorp/errwrap"
	// We must import sha512 so that it registers with the runtime so that
	// certificates that use it can be parsed.
	_ "crypto/sha512"
	"fmt"
	"io"
	"net"

	"github.com/hashicorp/vault/helper/proxyutil"
	"github.com/hashicorp/vault/helper/reload"
	"github.com/mitchellh/cli"
)

// ListenerFactory is the factory function to create a listener.
type ListenerFactory func(map[string]interface{}, io.Writer, cli.Ui) (net.Listener, map[string]string, reload.ReloadFunc, error)

// BuiltinListeners is the list of built-in listener types.
var BuiltinListeners = map[string]ListenerFactory{
	"tcp": tcpListenerFactory,
}

// NewListener creates a new listener of the given type with the given
// configuration. The type is looked up in the BuiltinListeners map.
func NewListener(t string, config map[string]interface{}, logger io.Writer, ui cli.Ui) (net.Listener, map[string]string, reload.ReloadFunc, error) {
	f, ok := BuiltinListeners[t]
	if !ok {
		return nil, nil, nil, fmt.Errorf("unknown listener type: %q", t)
	}

	return f(config, logger, ui)
}

func listenerWrapProxy(ln net.Listener, config map[string]interface{}) (net.Listener, error) {
	behaviorRaw, ok := config["proxy_protocol_behavior"]
	if !ok {
		return ln, nil
	}

	behavior, ok := behaviorRaw.(string)
	if !ok {
		return nil, fmt.Errorf("failed parsing proxy_protocol_behavior value: not a string")
	}

	proxyProtoConfig := &proxyutil.ProxyProtoConfig{
		Behavior: behavior,
	}

	if proxyProtoConfig.Behavior == "allow_authorized" || proxyProtoConfig.Behavior == "deny_unauthorized" {
		authorizedAddrsRaw, ok := config["proxy_protocol_authorized_addrs"]
		if !ok {
			return nil, fmt.Errorf("proxy_protocol_behavior set but no proxy_protocol_authorized_addrs value")
		}

		if err := proxyProtoConfig.SetAuthorizedAddrs(authorizedAddrsRaw); err != nil {
			return nil, errwrap.Wrapf("failed parsing proxy_protocol_authorized_addrs: {{err}}", err)
		}
	}

	newLn, err := proxyutil.WrapInProxyProto(ln, proxyProtoConfig)
	if err != nil {
		return nil, errwrap.Wrapf("failed configuring PROXY protocol wrapper: {{err}}", err)
	}

	return newLn, nil
}
