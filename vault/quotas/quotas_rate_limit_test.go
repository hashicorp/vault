// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package quotas

import (
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"go.uber.org/goleak"
)

type clientResult struct {
	atomicNumAllow *atomic.Int32
	atomicNumFail  *atomic.Int32
}

func TestNewRateLimitQuota(t *testing.T) {
	testCases := []struct {
		name      string
		rlq       *RateLimitQuota
		expectErr bool
	}{
		{"valid rate", NewRateLimitQuota("test-rate-limiter", "qa", "/foo/bar", "", "", false, time.Second, 0, 16.7), false},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink())
			require.Equal(t, tc.expectErr, err != nil, err)
			if err == nil {
				require.Nil(t, tc.rlq.close(context.Background()))
			}
		})
	}
}

func TestRateLimitQuota_Close(t *testing.T) {
	rlq := NewRateLimitQuota("test-rate-limiter", "qa", "/foo/bar", "", "", false, time.Second, time.Minute, 16.7)
	require.NoError(t, rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink()))
	require.NoError(t, rlq.close(context.Background()))

	time.Sleep(time.Second) // allow enough time for purgeClientsLoop to receive on closeCh
	require.False(t, rlq.getPurgeBlocked(), "expected blocked client purging to be disabled after explicit close")
}

func TestRateLimitQuota_Allow(t *testing.T) {
	rlq := &RateLimitQuota{
		Name:          "test-rate-limiter",
		Type:          TypeRateLimit,
		NamespacePath: "qa",
		MountPath:     "/foo/bar",
		Rate:          16.7,

		// override values to lower durations for testing purposes
		purgeInterval: 10 * time.Second,
		staleAge:      10 * time.Second,
	}

	require.NoError(t, rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink()))
	defer rlq.close(context.Background())

	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reqFunc := func(addr string, atomicNumAllow, atomicNumFail *atomic.Int32) {
		defer wg.Done()

		for ctx.Err() == nil {
			resp, err := rlq.allow(context.Background(), &Request{ClientAddress: addr})
			if err != nil {
				return
			}

			if resp.Allowed {
				atomicNumAllow.Add(1)
			} else {
				atomicNumFail.Add(1)
			}
			time.Sleep(2 * time.Millisecond)
		}
	}

	results := make(map[string]*clientResult)

	start := time.Now()

	for i := 0; i < 5; i++ {
		wg.Add(1)

		addr := fmt.Sprintf("127.0.0.%d", i)
		cr, ok := results[addr]
		if !ok {
			results[addr] = &clientResult{atomicNumAllow: atomic.NewInt32(0), atomicNumFail: atomic.NewInt32(0)}
			cr = results[addr]
		}

		go reqFunc(addr, cr.atomicNumAllow, cr.atomicNumFail)
	}

	wg.Wait()

	// evaluate the ideal RPS as (ceil(RPS) + (RPS * totalSeconds))
	elapsed := time.Since(start)
	ideal := math.Ceil(rlq.Rate) + (rlq.Rate * float64(elapsed) / float64(time.Second))

	for addr, cr := range results {
		numAllow := cr.atomicNumAllow.Load()
		numFail := cr.atomicNumFail.Load()

		// ensure there were some failed requests for the namespace
		require.NotZerof(t, numFail, "expected some requests to fail; addr: %s, numSuccess: %d, numFail: %d, elapsed: %s", addr, numAllow, numFail, elapsed)

		// ensure that we should never get more requests than allowed for the namespace
		want := int32(ideal + 1)
		require.Falsef(t, numAllow > want, "too many successful requests; addr: %s, want: %d, numSuccess: %d, numFail: %d, elapsed: %s", addr, want, numAllow, numFail, elapsed)
	}
}

func TestRateLimitQuota_Allow_WithBlock(t *testing.T) {
	rlq := &RateLimitQuota{
		Name:          "test-rate-limiter",
		Type:          TypeRateLimit,
		NamespacePath: "qa",
		MountPath:     "/foo/bar",
		Rate:          16.7,
		Interval:      5 * time.Second,
		BlockInterval: 10 * time.Second,

		// override values to lower durations for testing purposes
		purgeInterval: 10 * time.Second,
		staleAge:      10 * time.Second,
	}

	require.NoError(t, rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink()))
	defer rlq.close(context.Background())
	require.True(t, rlq.getPurgeBlocked())

	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reqFunc := func(addr string, atomicNumAllow, atomicNumFail *atomic.Int32) {
		defer wg.Done()

		for ctx.Err() == nil {
			resp, err := rlq.allow(ctx, &Request{ClientAddress: addr})
			if err != nil {
				return
			}

			if resp.Allowed {
				atomicNumAllow.Add(1)
			} else {
				atomicNumFail.Add(1)
			}
			time.Sleep(2 * time.Millisecond)
		}
	}

	results := make(map[string]*clientResult)

	for i := 0; i < 5; i++ {
		wg.Add(1)

		addr := fmt.Sprintf("127.0.0.%d", i)
		cr, ok := results[addr]
		if !ok {
			results[addr] = &clientResult{atomicNumAllow: atomic.NewInt32(0), atomicNumFail: atomic.NewInt32(0)}
			cr = results[addr]
		}

		go reqFunc(addr, cr.atomicNumAllow, cr.atomicNumFail)
	}

	wg.Wait()

	for _, cr := range results {
		numAllow := cr.atomicNumAllow.Load()
		numFail := cr.atomicNumFail.Load()

		// Since blocking is enabled, each client should only have 'rate' successful
		// requests, whereas all subsequent requests fail.
		require.Equal(t, int32(17), numAllow, "Expected 17 got %d allows with %d failures", numAllow, numFail)
		require.NotZero(t, numFail)
	}

	func() {
		timeout := time.After(rlq.purgeInterval * 2)
		ticker := time.Tick(time.Second)
		for {
			select {
			case <-timeout:
				require.Failf(t, "timeout exceeded waiting for blocked clients to be purged", "num blocked: %d", rlq.numBlockedClients())

			case <-ticker:
				if rlq.numBlockedClients() == 0 {
					return
				}
			}
		}
	}()
}

func TestRateLimitQuota_Update(t *testing.T) {
	defer goleak.VerifyNone(t)
	qm, err := NewManager(logging.NewVaultLogger(log.Trace), nil, metricsutil.BlackholeSink())
	require.NoError(t, err)

	quota := NewRateLimitQuota("quota1", "", "", "", "", false, time.Second, 0, 10)
	require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), quota, true))
	require.NoError(t, qm.SetQuota(context.Background(), TypeRateLimit.String(), quota, true))

	require.Nil(t, quota.close(context.Background()))
}
