// +build !enterprise

package vault

import (
	"context"
	"crypto/tls"
	"sync"

	cache "github.com/patrickmn/go-cache"
	"golang.org/x/net/http2"
	grpc "google.golang.org/grpc"
)

func perfStandbyRPCServer(*Core, *cache.Cache) *grpc.Server { return nil }

func handleReplicationConn(context.Context, *Core, *sync.WaitGroup, chan struct{}, *http2.Server, *grpc.Server, *cache.Cache, *tls.Conn) {
}
