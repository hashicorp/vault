// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cache

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

// Start spawns a goroutine that calls ttlcache's Start method.
// It waits to ensure that the goroutine has started before it returns.
// If sleep is true, it sleeps while polling to see if the goroutine has started.
// It only makes sense not to sleep when testing with synctest.Test.
// The purpose of this helper is to prevent cases where Stop is called on a cache
// that hasn't yet been started, which could cause a goroutine leak in some cases.
func Start[K comparable, V any](ctx context.Context, c *ttlcache.Cache[K, V], sleep bool) {
	go c.Start()
	for ctx.Err() == nil {
		if c.IsStarted() {
			return
		}
		if sleep {
			time.Sleep(1 * time.Millisecond)
		}
	}
}
