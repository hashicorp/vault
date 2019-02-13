package cache

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/reload"
	"github.com/mitchellh/cli"
)

func ServerListeners(lnConfigs []*config.Listener, logger io.Writer, ui cli.Ui) ([]net.Listener, error) {
	var listeners []net.Listener
	var listener net.Listener
	var err error
	for _, lnConfig := range lnConfigs {
		switch lnConfig.Type {
		case "unix":
			listener, _, _, err = unixSocketListener(lnConfig.Config, logger, ui)
			if err != nil {
				return nil, err
			}
			listeners = append(listeners, listener)
		case "tcp":
			listener, _, _, err := tcpListener(lnConfig.Config, logger, ui)
			if err != nil {
				return nil, err
			}
			listeners = append(listeners, listener)
		default:
			return nil, fmt.Errorf("unsupported listener type: %q", lnConfig.Type)
		}
	}

	return listeners, nil
}

func unixSocketListener(config map[string]interface{}, _ io.Writer, ui cli.Ui) (net.Listener, map[string]string, reload.ReloadFunc, error) {
	addr, ok := config["address"].(string)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid address: %v", config["address"])
	}

	if addr == "" {
		return nil, nil, nil, fmt.Errorf("address field should point to socket file path")
	}

	// Remove the socket file as it shouldn't exist for the domain socket to
	// work
	err := os.Remove(addr)
	if err != nil && !os.IsNotExist(err) {
		return nil, nil, nil, fmt.Errorf("failed to remove the socket file: %v", err)
	}

	listener, err := net.Listen("unix", addr)
	if err != nil {
		return nil, nil, nil, err
	}

	// Wrap the listener in rmListener so that the Unix domain socket file is
	// removed on close.
	listener = &rmListener{
		Listener: listener,
		Path:     addr,
	}

	props := map[string]string{"addr": addr}

	return server.ListenerWrapTLS(listener, props, config, ui)
}

func tcpListener(config map[string]interface{}, _ io.Writer, ui cli.Ui) (net.Listener, map[string]string, reload.ReloadFunc, error) {
	bindProto := "tcp"
	var addr string
	addrRaw, ok := config["address"]
	if !ok {
		addr = "127.0.0.1:8300"
	} else {
		addr = addrRaw.(string)
	}

	// If they've passed 0.0.0.0, we only want to bind on IPv4
	// rather than golang's dual stack default
	if strings.HasPrefix(addr, "0.0.0.0:") {
		bindProto = "tcp4"
	}

	ln, err := net.Listen(bindProto, addr)
	if err != nil {
		return nil, nil, nil, err
	}

	ln = server.TCPKeepAliveListener{ln.(*net.TCPListener)}

	props := map[string]string{"addr": addr}

	return server.ListenerWrapTLS(ln, props, config, ui)
}

// rmListener is an implementation of net.Listener that forwards most
// calls to the listener but also removes a file as part of the close. We
// use this to cleanup the unix domain socket on close.
type rmListener struct {
	net.Listener
	Path string
}

func (l *rmListener) Close() error {
	// Close the listener itself
	if err := l.Listener.Close(); err != nil {
		return err
	}

	// Remove the file
	return os.Remove(l.Path)
}
