// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudsqlconn

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/cloudsqlconn/debug"
	"cloud.google.com/go/cloudsqlconn/instance"
)

// monitoredCache is a wrapper around a connectionInfoCache that tracks the
// number of connections to the associated instance.
type monitoredCache struct {
	openConnsCount *uint64
	cn             instance.ConnName
	resolver       instance.ConnectionNameResolver
	logger         debug.ContextLogger

	// domainNameTicker periodically checks any domain names to see if they
	// changed.
	domainNameTicker *time.Ticker
	closedCh         chan struct{}

	mu        sync.Mutex
	openConns []*instrumentedConn
	closed    bool

	connectionInfoCache
}

func newMonitoredCache(
	ctx context.Context,
	cache connectionInfoCache,
	cn instance.ConnName,
	failoverPeriod time.Duration,
	resolver instance.ConnectionNameResolver,
	logger debug.ContextLogger) *monitoredCache {

	c := &monitoredCache{
		openConnsCount:      new(uint64),
		closedCh:            make(chan struct{}),
		cn:                  cn,
		resolver:            resolver,
		logger:              logger,
		connectionInfoCache: cache,
	}
	if cn.HasDomainName() {
		c.domainNameTicker = time.NewTicker(failoverPeriod)
		go func() {
			for {
				select {
				case <-c.domainNameTicker.C:
					c.purgeClosedConns()
					c.checkDomainName(ctx)
				case <-c.closedCh:
					return
				}
			}
		}()

	}

	return c
}
func (c *monitoredCache) isClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

func (c *monitoredCache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}

	c.closed = true
	close(c.closedCh)

	if c.domainNameTicker != nil {
		c.domainNameTicker.Stop()
	}

	if atomic.LoadUint64(c.openConnsCount) > 0 {
		for _, socket := range c.openConns {
			if !socket.isClosed() {
				_ = socket.Close() // force socket closed, ok to ignore error.
			}
		}
		atomic.StoreUint64(c.openConnsCount, 0)
	}

	return c.connectionInfoCache.Close()
}

func (c *monitoredCache) purgeClosedConns() {
	c.mu.Lock()
	defer c.mu.Unlock()

	var open []*instrumentedConn
	for _, s := range c.openConns {
		if !s.isClosed() {
			open = append(open, s)
		}
	}
	c.openConns = open
}

func (c *monitoredCache) checkDomainName(ctx context.Context) {
	if !c.cn.HasDomainName() {
		return
	}
	newCn, err := c.resolver.Resolve(ctx, c.cn.DomainName())
	if err != nil {
		// The domain name could not be resolved.
		c.logger.Debugf(ctx, "domain name %s for instance %s did not resolve, "+
			"closing all connections: %v",
			c.cn.DomainName(), c.cn.Name(), err)
		c.Close()
	}
	if newCn != c.cn {
		// The instance changed.
		c.logger.Debugf(ctx, "domain name %s changed from %s to %s, "+
			"closing all connections.",
			c.cn.DomainName(), c.cn.Name(), newCn.Name())
		c.Close()
	}

}
