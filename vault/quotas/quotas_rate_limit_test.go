package quotas

import (
	"fmt"
	"sync"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"go.uber.org/atomic"
)

func TestNewClientRateLimiter(t *testing.T) {
	testCases := []struct {
		maxRequests   float64
		burstSize     int
		expectedBurst int
	}{
		{1000, -1, 1000},
		{1000, 5000, 5000},
		{16.1, -1, 17},
		{16.7, -1, 17},
		{16.7, 100, 100},
	}

	for _, tc := range testCases {
		crl := newClientRateLimiter(tc.maxRequests, tc.burstSize)
		b := crl.limiter.Burst()
		if b != tc.expectedBurst {
			t.Fatalf("unexpected burst size; expected: %d, got: %d", tc.expectedBurst, b)
		}
	}
}

func TestNewRateLimitQuota(t *testing.T) {
	rlq := NewRateLimitQuota("test-rate-limiter", "qa", "/foo/bar", 16.7, 50)
	if err := rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink()); err != nil {
		t.Fatal(err)
	}

	if !rlq.purgeEnabled {
		t.Fatal("expected rate limit quota to start purge loop")
	}

	if rlq.purgeInterval != DefaultRateLimitPurgeInterval {
		t.Fatalf("unexpected purgeInterval; expected: %d, got: %d", DefaultRateLimitPurgeInterval, rlq.purgeInterval)
	}
	if rlq.staleAge != DefaultRateLimitStaleAge {
		t.Fatalf("unexpected staleAge; expected: %d, got: %d", DefaultRateLimitStaleAge, rlq.staleAge)
	}
}

func TestRateLimitQuota_Close(t *testing.T) {
	rlq := NewRateLimitQuota("test-rate-limiter", "qa", "/foo/bar", 16.7, 50)

	if err := rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink()); err != nil {
		t.Fatal(err)
	}

	if err := rlq.close(); err != nil {
		t.Fatalf("unexpected error when closing: %v", err)
	}

	time.Sleep(time.Second) // allow enough time for purgeClientsLoop to receive on closeCh

	if rlq.getPurgeEnabled() {
		t.Fatal("expected client purging to be disabled after close")
	}
}

func TestRateLimitQuota_Allow(t *testing.T) {
	rlq := &RateLimitQuota{
		Name:          "test-rate-limiter",
		Type:          TypeRateLimit,
		NamespacePath: "qa",
		MountPath:     "/foo/bar",
		Rate:          16.7,
		Burst:         83,
		purgeEnabled:  true, // to allow manual setting of purgeInterval and staleAge
	}

	if err := rlq.initialize(logging.NewVaultLogger(log.Trace), metricsutil.BlackholeSink()); err != nil {
		t.Fatal(err)
	}

	// override value and manually start purgeClientsLoop for testing purposes
	rlq.purgeInterval = 10 * time.Second
	rlq.staleAge = 10 * time.Second
	go rlq.purgeClientsLoop()

	var wg sync.WaitGroup

	type clientResult struct {
		atomicNumAllow *atomic.Int32
		atomicNumFail  *atomic.Int32
	}

	reqFunc := func(addr string, atomicNumAllow, atomicNumFail *atomic.Int32) {
		defer wg.Done()

		resp, err := rlq.allow(&Request{ClientAddress: addr})
		if err != nil {
			return
		}

		if resp.Allowed {
			atomicNumAllow.Add(1)
		} else {
			atomicNumFail.Add(1)
		}
	}

	results := make(map[string]*clientResult)

	start := time.Now()
	end := start.Add(5 * time.Second)
	for time.Now().Before(end) {

		for i := 0; i < 5; i++ {
			wg.Add(1)

			addr := fmt.Sprintf("127.0.0.%d", i)
			cr, ok := results[addr]
			if !ok {
				results[addr] = &clientResult{atomicNumAllow: atomic.NewInt32(0), atomicNumFail: atomic.NewInt32(0)}
				cr = results[addr]
			}

			go reqFunc(addr, cr.atomicNumAllow, cr.atomicNumFail)

			time.Sleep(2 * time.Millisecond)
		}

	}

	wg.Wait()

	if got, expected := len(results), rlq.NumClients(); got != expected {
		t.Fatalf("unexpected number of tracked client rate limit quotas; got %d, expected; %d", got, expected)
	}

	elapsed := time.Since(start)

	// evaluate the ideal RPS as (burst + (RPS * totalSeconds))
	ideal := float64(rlq.Burst) + (rlq.Rate * float64(elapsed) / float64(time.Second))

	for addr, cr := range results {
		numAllow := cr.atomicNumAllow.Load()
		numFail := cr.atomicNumFail.Load()

		// ensure there were some failed requests for the namespace
		if numFail == 0 {
			t.Fatalf("expected some requests to fail; addr: %s, numSuccess: %d, numFail: %d, elapsed: %d", addr, numAllow, numFail, elapsed)
		}

		// ensure that we should never get more requests than allowed for the namespace
		if want := int32(ideal + 1); numAllow > want {
			t.Fatalf("too many successful requests; addr: %s, want: %d, numSuccess: %d, numFail: %d, elapsed: %d", addr, want, numAllow, numFail, elapsed)
		}
	}

	// allow enough time for the client to be purged
	time.Sleep(rlq.purgeInterval * 2)

	for addr := range results {
		if rlq.HasClient(addr) {
			t.Fatalf("expected stale client to be purged: %s", addr)
		}
	}
}
