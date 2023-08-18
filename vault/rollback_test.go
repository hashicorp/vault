// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/stretchr/testify/require"
)

// mockRollback returns a mock rollback manager
func mockRollback(t *testing.T) (*RollbackManager, *NoopBackend) {
	backend := new(NoopBackend)
	mounts := new(MountTable)
	router := NewRouter()
	core, _, _ := TestCoreUnsealed(t)

	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")

	mounts.Entries = []*MountEntry{
		{
			Path:        "foo",
			NamespaceID: namespace.RootNamespaceID,
			namespace:   namespace.RootNamespace,
		},
	}
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}

	if err := router.Mount(backend, "foo", &MountEntry{UUID: meUUID, Accessor: "noopaccessor", NamespaceID: namespace.RootNamespaceID, namespace: namespace.RootNamespace}, view); err != nil {
		t.Fatalf("err: %s", err)
	}

	mountsFunc := func() []*MountEntry {
		return mounts.Entries
	}

	logger := logging.NewVaultLogger(log.Trace)

	rb := NewRollbackManager(context.Background(), logger, mountsFunc, router, core)
	rb.period = 10 * time.Millisecond
	return rb, backend
}

func TestRollbackManager(t *testing.T) {
	m, backend := mockRollback(t)
	if len(backend.Paths) > 0 {
		t.Fatalf("bad: %#v", backend)
	}

	m.Start()
	time.Sleep(50 * time.Millisecond)
	m.Stop()

	count := len(backend.Paths)
	if count == 0 {
		t.Fatalf("bad: %#v", backend)
	}
	if backend.Paths[0] != "" {
		t.Fatalf("bad: %#v", backend)
	}

	time.Sleep(50 * time.Millisecond)

	if count != len(backend.Paths) {
		t.Fatalf("should stop requests: %#v", backend)
	}
}

func TestRollbackManager_Join(t *testing.T) {
	m, backend := mockRollback(t)
	if len(backend.Paths) > 0 {
		t.Fatalf("bad: %#v", backend)
	}

	m.Start()
	defer m.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	errCh := make(chan error, 3)
	go func() {
		defer wg.Done()
		err := m.Rollback(namespace.RootContext(nil), "foo")
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := m.Rollback(namespace.RootContext(nil), "foo")
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()
		err := m.Rollback(namespace.RootContext(nil), "foo")
		if err != nil {
			errCh <- err
		}
	}()
	wg.Wait()
	close(errCh)
	err := <-errCh
	if err != nil {
		t.Fatalf("Error on rollback:%v", err)
	}
}

// TestRollbackMetrics verifies that the rollback metrics only include the mount
// point in their names when RollbackMetricsIncludeMountPoint is true.
// This test cannot be run in parallel, because we are using the global metrics
// instance
func TestRollbackMetrics(t *testing.T) {
	testCases := []struct {
		name          string
		addMountPoint bool
	}{
		{
			name:          "include mount point",
			addMountPoint: true,
		},
		{
			name:          "exclude mount point",
			addMountPoint: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inMemSink := metrics.NewInmemSink(10000*time.Hour, 10000*time.Hour)
			sink := metricsutil.NewClusterMetricSink("test", inMemSink)
			sink.TelemetryConsts.RollbackMetricsIncludeMountPoint = tc.addMountPoint
			_, err := metrics.NewGlobal(metrics.DefaultConfig("vault"), inMemSink)
			require.NoError(t, err)
			conf := &CoreConfig{
				MetricSink:     sink,
				RollbackPeriod: 50 * time.Millisecond,
				MetricsHelper:  metricsutil.NewMetricsHelper(inMemSink, true),
			}

			core, _, _ := TestCoreUnsealedWithConfig(t, conf)

			samplesWith := func(intervals []*metrics.IntervalMetrics, with func(string) bool) []metrics.SampledValue {
				t.Helper()
				samples := make([]metrics.SampledValue, 0)
				for _, interval := range intervals {
					for name, summary := range interval.Samples {
						if with(name) {
							samples = append(samples, summary)
						}
					}
				}
				return samples
			}

			<-core.rollback.rollbacksDoneCh
			intervals := inMemSink.Data()

			mountPointAttempts := samplesWith(intervals, func(s string) bool {
				return strings.HasPrefix(s, "vault.rollback.attempt.")
			})
			mountPointRoutes := samplesWith(intervals, func(s string) bool {
				return strings.HasPrefix(s, "vault.route.rollback.")
			})

			noMountPointAttempts := samplesWith(intervals, func(s string) bool {
				return s == "vault.rollback.attempt"
			})
			noMountPointRoutes := samplesWith(intervals, func(s string) bool {
				return s == "vault.route.rollback"
			})
			if tc.addMountPoint {
				require.NotEmpty(t, mountPointAttempts)
				require.NotEmpty(t, mountPointRoutes)
				require.Empty(t, noMountPointAttempts)
				require.Empty(t, noMountPointRoutes)
			} else {
				require.Empty(t, mountPointAttempts)
				require.Empty(t, mountPointRoutes)
				require.NotEmpty(t, noMountPointAttempts)
				require.NotEmpty(t, noMountPointRoutes)
			}
		})
	}
}
