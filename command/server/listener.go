package server

import (
	"fmt"
	"net"
)

// ListenerFactory is the factory function to create a listener.
type ListenerFactory func(map[string]string) (net.Listener, error)

// BuiltinListeners is the list of built-in listener types.
var BuiltinListeners = map[string]ListenerFactory{
	"tcp": tcpListenerFactory,
}

// NewListener creates a new listener of the given type with the given
// configuration. The type is looked up in the BuiltinListeners map.
func NewListener(t string, config map[string]string) (net.Listener, error) {
	f, ok := BuiltinListeners[t]
	if !ok {
		return nil, fmt.Errorf("unknown listener type: %s", t)
	}

	return f(config)
}
