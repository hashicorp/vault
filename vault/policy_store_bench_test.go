// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// slowGetStorage delays every Get to model the latency of a real storage
// backend, which is what a policy cache miss waits on.
type slowGetStorage struct {
	logical.Storage
	delay time.Duration
}

func (s *slowGetStorage) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	time.Sleep(s.delay)
	return s.Storage.Get(ctx, key)
}

// BenchmarkPolicyStore_ColdLoad measures how long it takes concurrent
// requests to load a set of distinct policies into an empty policy cache.
// Each op purges the cache and loads every policy once.
func BenchmarkPolicyStore_ColdLoad(b *testing.B) {
	const policyCount = 256

	for _, delay := range []time.Duration{0, time.Millisecond} {
		for _, workers := range []int{1, 16, 64} {
			b.Run(fmt.Sprintf("storage_delay=%s/workers=%d", delay, workers), func(b *testing.B) {
				cluster := NewTestCluster(b, &CoreConfig{}, nil)
				ps := cluster.Cores[0].Core.policyStore
				ctx := namespace.RootContext(context.Background())

				names := make([]string, policyCount)
				for i := range names {
					names[i] = fmt.Sprintf("bench-policy-%d", i)
					p, err := ParseACLPolicy(namespace.RootNamespace, fmt.Sprintf(`path "secret/%d/*" { capabilities = ["read"] }`, i))
					if err != nil {
						b.Fatal(err)
					}
					p.Name = names[i]
					if err := ps.SetPolicy(ctx, p); err != nil {
						b.Fatal(err)
					}
				}
				ps.aclView.storage = logical.NewStorageView(&slowGetStorage{Storage: ps.aclView.storage, delay: delay}, "")

				b.ResetTimer()
				for range b.N {
					b.StopTimer()
					ps.tokenPoliciesLRU.Purge()
					var next atomic.Int64
					var wg sync.WaitGroup
					b.StartTimer()

					for range workers {
						wg.Add(1)
						go func() {
							defer wg.Done()
							for {
								n := next.Add(1) - 1
								if n >= policyCount {
									return
								}
								p, err := ps.GetPolicy(ctx, names[n], PolicyTypeToken)
								if err != nil || p == nil {
									b.Errorf("GetPolicy(%s) = %v, %v", names[n], p, err)
									return
								}
							}
						}()
					}
					wg.Wait()
				}
				b.ReportMetric(float64(b.Elapsed().Microseconds())/float64(b.N*policyCount), "us/policy")
			})
		}
	}
}
