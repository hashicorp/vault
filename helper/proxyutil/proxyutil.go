package proxyutil

import (
	"sync"

	sockaddr "github.com/hashicorp/go-sockaddr"
)

// ProxyProtoConfig contains configuration for the PROXY protocol
type ProxyProtoConfig struct {
	sync.RWMutex
	AllowedAddrs []sockaddr.SockAddr
}
