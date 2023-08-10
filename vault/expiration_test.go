// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/benchhelpers"
	"github.com/hashicorp/vault/helper/fairshare"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
)

var testImagePull sync.Once

// mockExpiration returns a mock expiration manager
func mockExpiration(t testing.TB) *ExpirationManager {
	c, _, _ := TestCoreUnsealed(benchhelpers.TBtoT(t))

	// Wait until the expiration manager is out of restore mode.
	// This was added to prevent sporadic failures of TestExpiration_unrecoverableErrorMakesIrrevocable.
	timeout := time.Now().Add(time.Second * 10)
	for c.expiration.inRestoreMode() {
		if time.Now().After(timeout) {
			t.Fatal("ExpirationManager is still in restore mode after 10 seconds")
		}
		time.Sleep(50 * time.Millisecond)
	}

	return c.expiration
}

func mockBackendExpiration(t testing.TB, backend physical.Backend) (*Core, *ExpirationManager) {
	c, _, _ := TestCoreUnsealedBackend(benchhelpers.TBtoT(t), backend)
	return c, c.expiration
}

func TestExpiration_Metrics(t *testing.T) {
	var err error

	testCore := TestCore(t)
	testCore.baseLogger = logger
	testCore.logger = logger.Named("core")
	testCoreUnsealed(t, testCore)

	exp := testCore.expiration

	if err := exp.Restore(nil); err != nil {
		t.Fatal(err)
	}

	// Set up a count function to calculate number of leases
	count := 0
	countFunc := func(_ string) {
		count++
	}

	// Scan the storage with the count func set
	if err = logical.ScanView(namespace.RootContext(nil), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Check that there are no leases to begin with
	if count != 0 {
		t.Fatalf("bad: lease count; expected:0 actual:%d", count)
	}

	for i := 0; i < 50; i++ {
		le := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i),
			Path:       "foo/bar/" + fmt.Sprintf("%d", i),
			namespace:  namespace.RootNamespace,
			IssueTime:  time.Now(),
			ExpireTime: time.Now().Add(time.Hour),
		}

		otherNS := &namespace.Namespace{
			ID:   "nsid",
			Path: "foo/bar",
		}

		otherNSle := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i) + "/blah.nsid",
			Path:       "foo/bar/" + fmt.Sprintf("%d", i) + "/blah.nsid",
			namespace:  otherNS,
			IssueTime:  time.Now(),
			ExpireTime: time.Now().Add(time.Hour),
		}

		exp.pendingLock.Lock()
		if err := exp.persistEntry(namespace.RootContext(nil), le); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting entry: %v", err)
		}
		exp.updatePendingInternal(le)

		if err := exp.persistEntry(namespace.RootContext(nil), otherNSle); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting entry: %v", err)
		}
		exp.updatePendingInternal(otherNSle)
		exp.pendingLock.Unlock()
	}

	for i := 50; i < 250; i++ {
		le := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i+1),
			Path:       "foo/bar/" + fmt.Sprintf("%d", i+1),
			namespace:  namespace.RootNamespace,
			IssueTime:  time.Now(),
			ExpireTime: time.Now().Add(2 * time.Hour),
		}

		exp.pendingLock.Lock()
		if err := exp.persistEntry(namespace.RootContext(nil), le); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting entry: %v", err)
		}
		exp.updatePendingInternal(le)
		exp.pendingLock.Unlock()
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	var conf metricsutil.TelemetryConstConfig = metricsutil.TelemetryConstConfig{
		LeaseMetricsEpsilon:         time.Hour,
		NumLeaseMetricsTimeBuckets:  2,
		LeaseMetricsNameSpaceLabels: true,
	}

	flattenedResults, err := exp.leaseAggregationMetrics(context.Background(), conf)
	if err != nil {
		t.Fatal(err)
	}
	if flattenedResults == nil {
		t.Fatal("lease aggregation returns nil metrics")
	}

	labelOneHour := metrics.Label{"expiring", time.Now().Add(time.Hour).Round(time.Hour).String()}
	labelTwoHours := metrics.Label{"expiring", time.Now().Add(2 * time.Hour).Round(time.Hour).String()}
	nsLabel := metrics.Label{"namespace", "root"}
	nsLabelNonRoot := metrics.Label{"namespace", "nsid"}

	foundLabelOne := false
	foundLabelTwo := false
	foundLabelThree := false

	for _, labelVal := range flattenedResults {
		retNsLabel := labelVal.Labels[1]
		retTimeLabel := labelVal.Labels[0]
		if nsLabel == retNsLabel {
			if labelVal.Value == 50 {
				if retTimeLabel == labelOneHour {
					foundLabelOne = true
				}
			}
			if labelVal.Value == 200 {
				if retTimeLabel == labelTwoHours {
					foundLabelTwo = true
				}
			}
		} else if retNsLabel == nsLabelNonRoot {
			if labelVal.Value == 50 {
				if retTimeLabel == labelOneHour {
					foundLabelThree = true
				}
			}
		}
	}

	if !foundLabelOne || !foundLabelTwo || !foundLabelThree {
		t.Errorf("One of the labels is missing. one: %t, two: %t, three: %t", foundLabelOne, foundLabelTwo, foundLabelThree)
	}

	// test the same leases while ignoring namespaces so the 2 different namespaces get aggregated
	conf = metricsutil.TelemetryConstConfig{
		LeaseMetricsEpsilon:         time.Hour,
		NumLeaseMetricsTimeBuckets:  2,
		LeaseMetricsNameSpaceLabels: false,
	}

	flattenedResults, err = exp.leaseAggregationMetrics(context.Background(), conf)
	if err != nil {
		t.Fatal(err)
	}
	if flattenedResults == nil {
		t.Fatal("lease aggregation returns nil metrics")
	}

	foundLabelOne = false
	foundLabelTwo = false

	for _, labelVal := range flattenedResults {
		if len(labelVal.Labels) != 1 {
			t.Errorf("Namespace label is returned when explicitly not requested.")
		}
		retTimeLabel := labelVal.Labels[0]
		if labelVal.Value == 100 {
			if retTimeLabel == labelOneHour {
				foundLabelOne = true
			}
		}
		if labelVal.Value == 200 {
			if retTimeLabel == labelTwoHours {
				foundLabelTwo = true
			}
		}
	}
	if !foundLabelOne || !foundLabelTwo {
		t.Errorf("One of the labels is missing")
	}
}

func TestExpiration_TotalLeaseCount(t *testing.T) {
	// Quotas and internal lease count tracker are coupled, so this is a proxy
	// for testing the total lease count quota
	c, _, _ := TestCoreUnsealed(t)
	exp := c.expiration

	expectedCount := 0
	otherNS := &namespace.Namespace{
		ID:   "nsid",
		Path: "foo/bar",
	}
	for i := 0; i < 50; i++ {
		le := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i),
			Path:       "foo/bar/" + fmt.Sprintf("%d", i),
			namespace:  namespace.RootNamespace,
			IssueTime:  time.Now(),
			ExpireTime: time.Now().Add(time.Hour),
		}

		otherNSle := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i) + "/blah.nsid",
			Path:       "foo/bar/" + fmt.Sprintf("%d", i) + "/blah.nsid",
			namespace:  otherNS,
			IssueTime:  time.Now(),
			ExpireTime: time.Now().Add(time.Hour),
		}

		exp.pendingLock.Lock()
		if err := exp.persistEntry(namespace.RootContext(nil), le); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting irrevocable entry: %v", err)
		}
		exp.updatePendingInternal(le)
		expectedCount++

		if err := exp.persistEntry(namespace.RootContext(nil), otherNSle); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting irrevocable entry: %v", err)
		}
		exp.updatePendingInternal(otherNSle)
		expectedCount++
		exp.pendingLock.Unlock()
	}

	// add some irrevocable leases to each count to ensure they are counted too
	// note: irrevocable leases almost certainly have an expire time set in the
	// past, but for this exercise it should be fine to set it to whatever
	for i := 50; i < 60; i++ {
		le := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i+1),
			Path:       "foo/bar/" + fmt.Sprintf("%d", i+1),
			namespace:  namespace.RootNamespace,
			IssueTime:  time.Now(),
			ExpireTime: time.Now(),
			RevokeErr:  "some err message",
		}

		otherNSle := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i+1) + "/blah.nsid",
			Path:       "foo/bar/" + fmt.Sprintf("%d", i+1) + "/blah.nsid",
			namespace:  otherNS,
			IssueTime:  time.Now(),
			ExpireTime: time.Now(),
			RevokeErr:  "some err message",
		}

		exp.pendingLock.Lock()
		if err := exp.persistEntry(namespace.RootContext(nil), le); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting irrevocable entry: %v", err)
		}
		exp.updatePendingInternal(le)
		expectedCount++

		if err := exp.persistEntry(namespace.RootContext(nil), otherNSle); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting irrevocable entry: %v", err)
		}
		exp.updatePendingInternal(otherNSle)
		expectedCount++
		exp.pendingLock.Unlock()
	}

	exp.pendingLock.RLock()
	count := exp.leaseCount
	exp.pendingLock.RUnlock()

	if count != expectedCount {
		t.Errorf("bad lease count. expected %d, got %d", expectedCount, count)
	}
}

func TestExpiration_TotalLeaseCount_WithRoles(t *testing.T) {
	// Quotas and internal lease count tracker are coupled, so this is a proxy
	// for testing the total lease count quota
	c, _, _ := TestCoreUnsealed(t)
	exp := c.expiration

	expectedCount := 0
	otherNS := &namespace.Namespace{
		ID:   "nsid",
		Path: "foo/bar",
	}
	for i := 0; i < 50; i++ {
		le := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i),
			Path:       "foo/bar/" + fmt.Sprintf("%d", i),
			LoginRole:  "loginRole" + fmt.Sprintf("%d", i),
			namespace:  namespace.RootNamespace,
			IssueTime:  time.Now(),
			ExpireTime: time.Now().Add(time.Hour),
		}

		otherNSle := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i) + "/blah.nsid",
			Path:       "foo/bar/" + fmt.Sprintf("%d", i) + "/blah.nsid",
			LoginRole:  "loginRole" + fmt.Sprintf("%d", i),
			namespace:  otherNS,
			IssueTime:  time.Now(),
			ExpireTime: time.Now().Add(time.Hour),
		}

		exp.pendingLock.Lock()
		if err := exp.persistEntry(namespace.RootContext(nil), le); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting irrevocable entry: %v", err)
		}
		exp.updatePendingInternal(le)
		expectedCount++

		if err := exp.persistEntry(namespace.RootContext(nil), otherNSle); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting irrevocable entry: %v", err)
		}
		exp.updatePendingInternal(otherNSle)
		expectedCount++
		exp.pendingLock.Unlock()
	}

	// add some irrevocable leases to each count to ensure they are counted too
	// note: irrevocable leases almost certainly have an expire time set in the
	// past, but for this exercise it should be fine to set it to whatever
	for i := 50; i < 60; i++ {
		le := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i+1),
			Path:       "foo/bar/" + fmt.Sprintf("%d", i+1),
			LoginRole:  "loginRole" + fmt.Sprintf("%d", i),
			namespace:  namespace.RootNamespace,
			IssueTime:  time.Now(),
			ExpireTime: time.Now(),
			RevokeErr:  "some err message",
		}

		otherNSle := &leaseEntry{
			LeaseID:    "lease" + fmt.Sprintf("%d", i+1) + "/blah.nsid",
			Path:       "foo/bar/" + fmt.Sprintf("%d", i+1) + "/blah.nsid",
			LoginRole:  "loginRole" + fmt.Sprintf("%d", i),
			namespace:  otherNS,
			IssueTime:  time.Now(),
			ExpireTime: time.Now(),
			RevokeErr:  "some err message",
		}

		exp.pendingLock.Lock()
		if err := exp.persistEntry(namespace.RootContext(nil), le); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting irrevocable entry: %v", err)
		}
		exp.updatePendingInternal(le)
		expectedCount++

		if err := exp.persistEntry(namespace.RootContext(nil), otherNSle); err != nil {
			exp.pendingLock.Unlock()
			t.Fatalf("error persisting irrevocable entry: %v", err)
		}
		exp.updatePendingInternal(otherNSle)
		expectedCount++
		exp.pendingLock.Unlock()
	}

	exp.pendingLock.RLock()
	count := exp.leaseCount
	exp.pendingLock.RUnlock()

	if count != expectedCount {
		t.Errorf("bad lease count. expected %d, got %d", expectedCount, count)
	}
}

func TestExpiration_Tidy(t *testing.T) {
	var err error

	// We use this later for tidy testing where we need to check the output
	logOut := new(bytes.Buffer)
	logger := log.New(&log.LoggerOptions{
		Output: logOut,
	})

	testCore := TestCore(t)
	testCore.baseLogger = logger
	testCore.logger = logger.Named("core")
	testCoreUnsealed(t, testCore)

	exp := testCore.expiration

	if err := exp.Restore(nil); err != nil {
		t.Fatal(err)
	}

	// Set up a count function to calculate number of leases
	count := 0
	countFunc := func(leaseID string) {
		count++
	}

	// Scan the storage with the count func set
	if err = logical.ScanView(namespace.RootContext(nil), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Check that there are no leases to begin with
	if count != 0 {
		t.Fatalf("bad: lease count; expected:0 actual:%d", count)
	}

	// Create a lease entry without a client token in it
	le := &leaseEntry{
		LeaseID:   "lease/with/no/client/token",
		Path:      "foo/bar",
		namespace: namespace.RootNamespace,
	}

	// Persist the invalid lease entry
	if err = exp.persistEntry(namespace.RootContext(nil), le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Check that the storage was successful and that the count of leases is
	// now 1
	if count != 1 {
		t.Fatalf("bad: lease count; expected:1 actual:%d", count)
	}

	// Run the tidy operation
	err = exp.Tidy(namespace.RootContext(nil))
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	if err := logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Post the tidy operation, the invalid lease entry should have been gone
	if count != 0 {
		t.Fatalf("bad: lease count; expected:0 actual:%d", count)
	}

	// Set a revoked/invalid token in the lease entry
	le.ClientToken = "invalidtoken"

	// Persist the invalid lease entry
	if err = exp.persistEntry(namespace.RootContext(nil), le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Check that the storage was successful and that the count of leases is
	// now 1
	if count != 1 {
		t.Fatalf("bad: lease count; expected:1 actual:%d", count)
	}

	// Run the tidy operation
	err = exp.Tidy(namespace.RootContext(nil))
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Post the tidy operation, the invalid lease entry should have been gone
	if count != 0 {
		t.Fatalf("bad: lease count; expected:0 actual:%d", count)
	}

	// Attach an invalid token with 2 leases
	if err = exp.persistEntry(namespace.RootContext(nil), le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}

	le.LeaseID = "another/invalid/lease"
	if err = exp.persistEntry(context.Background(), le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}

	// Run the tidy operation
	err = exp.Tidy(namespace.RootContext(nil))
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Post the tidy operation, the invalid lease entry should have been gone
	if count != 0 {
		t.Fatalf("bad: lease count; expected:0 actual:%d", count)
	}

	for i := 0; i < 1000; i++ {
		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        "invalid/lease/" + fmt.Sprintf("%d", i+1),
			ClientToken: "invalidtoken",
		}
		req.SetTokenEntry(&logical.TokenEntry{ID: "invalidtoken", NamespaceID: "root"})
		resp := &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL: 100 * time.Millisecond,
				},
			},
			Data: map[string]interface{}{
				"test_key": "test_value",
			},
		}
		_, err := exp.Register(namespace.RootContext(nil), req, resp, "")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Check that there are 1000 leases now
	if count != 1000 {
		t.Fatalf("bad: lease count; expected:1000 actual:%d", count)
	}

	errCh1 := make(chan error)
	errCh2 := make(chan error)

	// Initiate tidy of the above 1000 invalid leases in quick succession. Only
	// one tidy operation can be in flight at any time. One of these requests
	// should error out.
	go func() {
		errCh1 <- exp.Tidy(namespace.RootContext(nil))
	}()

	go func() {
		errCh2 <- exp.Tidy(namespace.RootContext(nil))
	}()

	var err1, err2 error

	for i := 0; i < 2; i++ {
		select {
		case err1 = <-errCh1:
		case err2 = <-errCh2:
		}
	}

	if err1 != nil || err2 != nil {
		t.Fatalf("got an error: err1: %v; err2: %v", err1, err2)
	}
	if !strings.Contains(logOut.String(), "tidy operation on leases is already in progress") {
		t.Fatalf("expected to see a warning saying operation in progress, output is %s", logOut.String())
	}

	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	le.ClientToken = root.ID

	// Attach a valid token with the leases
	if err = exp.persistEntry(namespace.RootContext(nil), le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}

	// Run the tidy operation
	err = exp.Tidy(namespace.RootContext(nil))
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Post the tidy operation, the valid lease entry should not get affected
	if count != 1 {
		t.Fatalf("bad: lease count; expected:1 actual:%d", count)
	}
}

// To avoid pulling in deps for all users of the package, don't leave these
// uncommented in the public tree
/*
func BenchmarkExpiration_Restore_Etcd(b *testing.B) {
	addr := os.Getenv("PHYSICAL_BACKEND_BENCHMARK_ADDR")
	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())

	logger := logging.NewVaultLogger(log.Trace)
	physicalBackend, err := physEtcd.NewEtcdBackend(map[string]string{
		"address":      addr,
		"path":         randPath,
		"max_parallel": "256",
	}, logger)
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	benchmarkExpirationBackend(b, physicalBackend, 10000) // 10,000 leases
}

func BenchmarkExpiration_Restore_Consul(b *testing.B) {
	addr := os.Getenv("PHYSICAL_BACKEND_BENCHMARK_ADDR")
	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())

	logger := logging.NewVaultLogger(log.Trace)
	physicalBackend, err := physConsul.NewConsulBackend(map[string]string{
		"address":      addr,
		"path":         randPath,
		"max_parallel": "256",
	}, logger)
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	benchmarkExpirationBackend(b, physicalBackend, 10000) // 10,000 leases
}
*/

func BenchmarkExpiration_Restore_InMem(b *testing.B) {
	logger := logging.NewVaultLogger(log.Trace)
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkExpirationBackend(b, inm, 100000) // 100,000 Leases
}

func benchmarkExpirationBackend(b *testing.B, physicalBackend physical.Backend, numLeases int) {
	c, _, _ := TestCoreUnsealedBackend(benchhelpers.TBtoT(b), physicalBackend)
	exp := c.expiration
	noop := &NoopBackend{}
	view := NewBarrierView(c.barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		b.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		b.Fatal(err)
	}

	// Register fake leases
	for i := 0; i < numLeases; i++ {
		pathUUID, err := uuid.GenerateUUID()
		if err != nil {
			b.Fatal(err)
		}

		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        "prod/aws/" + pathUUID,
			ClientToken: "root",
		}
		req.SetTokenEntry(&logical.TokenEntry{ID: "root", NamespaceID: "root"})
		resp := &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL: 400 * time.Second,
				},
			},
			Data: map[string]interface{}{
				"access_key": "xyz",
				"secret_key": "abcd",
			},
		}
		_, err = exp.Register(namespace.RootContext(nil), req, resp, "")
		if err != nil {
			b.Fatalf("err: %v", err)
		}
	}

	// Stop everything
	err = exp.Stop()
	if err != nil {
		b.Fatalf("err: %v", err)
	}
	// Avoid panic due to calling exp.Stop multiple times
	c.expiration = nil

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = exp.Restore(nil)
		// Restore
		if err != nil {
			b.Fatalf("err: %v", err)
		}
	}
	b.StopTimer()
}

func BenchmarkExpiration_Create_Leases(b *testing.B) {
	logger := logging.NewVaultLogger(log.Trace)
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		b.Fatal(err)
	}

	c, _, _ := TestCoreUnsealedBackend(benchhelpers.TBtoT(b), inm)
	exp := c.expiration
	noop := &NoopBackend{}
	view := NewBarrierView(c.barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		b.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		b.Fatal(err)
	}
	req := &logical.Request{
		Operation:   logical.ReadOperation,
		ClientToken: "root",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "root", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: 400 * time.Second,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Path = fmt.Sprintf("prod/aws/%d", i)
		_, err = exp.Register(namespace.RootContext(nil), req, resp, "")
		if err != nil {
			b.Fatalf("err: %v", err)
		}
	}
}

func TestExpiration_Restore(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	exp := c.expiration
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	paths := []string{
		"prod/aws/foo",
		"prod/aws/sub/bar",
		"prod/aws/zip",
	}
	for _, path := range paths {
		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        path,
			ClientToken: "foobar",
		}
		req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
		resp := &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL: time.Second,
				},
			},
			Data: map[string]interface{}{
				"access_key": "xyz",
				"secret_key": "abcd",
			},
		}
		_, err := exp.Register(namespace.RootContext(nil), req, resp, "")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	if exp.leaseCount != len(paths) {
		t.Fatalf("expected %v leases, got %v", len(paths), exp.leaseCount)
	}

	// Stop everything
	err = c.stopExpiration()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if exp.leaseCount != 0 {
		t.Fatalf("expected %v leases, got %v", 0, exp.leaseCount)
	}

	// Restore
	err = exp.Restore(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Would like to test here, but this is a race with the expiration of the leases.

	// Ensure all are reaped
	start := time.Now()
	for time.Now().Sub(start) < time.Second {
		noop.Lock()
		less := len(noop.Requests) < 3
		noop.Unlock()

		if less {
			time.Sleep(5 * time.Millisecond)
			continue
		}
		break
	}
	for _, req := range noop.Requests {
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}
	}

	if exp.leaseCount != 0 {
		t.Fatalf("expected %v leases, got %v", 0, exp.leaseCount)
	}
}

func TestExpiration_Register(t *testing.T) {
	exp := mockExpiration(t)
	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Hour,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	id, err := exp.Register(namespace.RootContext(nil), req, resp, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !strings.HasPrefix(id, req.Path) {
		t.Fatalf("bad: %s", id)
	}

	if len(id) <= len(req.Path) {
		t.Fatalf("bad: %s", id)
	}
}

func TestExpiration_Register_Role(t *testing.T) {
	exp := mockExpiration(t)
	role := "role1"
	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Hour,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	id, err := exp.Register(namespace.RootContext(nil), req, resp, role)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !strings.HasPrefix(id, req.Path) {
		t.Fatalf("bad: %s", id)
	}

	if len(id) <= len(req.Path) {
		t.Fatalf("bad: %s", id)
	}

	le, err := exp.loadEntry(exp.quitContext, id)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if le.LoginRole != role {
		t.Fatalf("Login role incorrect. Expected %s, received %s", role, le.LoginRole)
	}
}

func TestExpiration_Register_BatchToken(t *testing.T) {
	c, _, rootToken := TestCoreUnsealed(t)
	exp := c.expiration
	noop := &NoopBackend{
		RequestHandler: func(ctx context.Context, req *logical.Request) (*logical.Response, error) {
			resp := &logical.Response{Secret: req.Secret}
			resp.Secret.TTL = time.Hour
			return resp, nil
		},
	}
	{
		_, barrier, _ := mockBarrier(t)
		view := NewBarrierView(barrier, "logical/")
		meUUID, err := uuid.GenerateUUID()
		if err != nil {
			t.Fatal(err)
		}
		err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
		if err != nil {
			t.Fatal(err)
		}
	}

	te := &logical.TokenEntry{
		Type:         logical.TokenTypeBatch,
		TTL:          1 * time.Second,
		NamespaceID:  "root",
		CreationTime: time.Now().Unix(),
		Parent:       rootToken,
	}

	err := exp.tokenStore.create(context.Background(), te)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: te.ID,
	}
	req.SetTokenEntry(te)
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Hour,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	leaseID, err := exp.Register(namespace.RootContext(nil), req, resp, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	_, err = exp.Renew(namespace.RootContext(nil), leaseID, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	reqID := 0
	for {
		if time.Now().Sub(start) > 10*time.Second {
			t.Fatal("didn't revoke lease")
		}
		req = nil

		noop.Lock()
		if len(noop.Requests) > reqID {
			req = noop.Requests[reqID]
			reqID++
		}
		noop.Unlock()
		if req == nil || req.Operation == logical.RenewOperation {
			time.Sleep(5 * time.Millisecond)
			continue
		}
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}

		break
	}

	deadline := time.Now().Add(5 * time.Second)
	var idEnts []string
	for time.Now().Before(deadline) {
		idEnts, err = exp.tokenView.List(context.Background(), "")
		if err != nil {
			t.Fatal(err)
		}
		if len(idEnts) == 0 {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("expected no entries in sys/expire/token, got: %v", idEnts)
}

func TestExpiration_RegisterAuth(t *testing.T) {
	exp := mockExpiration(t)

	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
	}

	te := &logical.TokenEntry{
		Path:        "auth/github/login",
		NamespaceID: namespace.RootNamespaceID,
	}
	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	te = &logical.TokenEntry{
		Path:        "auth/github/../login",
		NamespaceID: namespace.RootNamespaceID,
	}
	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExpiration_RegisterAuth_Role(t *testing.T) {
	exp := mockExpiration(t)
	role := "role1"
	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
	}

	te := &logical.TokenEntry{
		Path:        "auth/github/login",
		NamespaceID: namespace.RootNamespaceID,
	}
	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, role)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	te = &logical.TokenEntry{
		Path:        "auth/github/../login",
		NamespaceID: namespace.RootNamespaceID,
	}
	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, role)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExpiration_RegisterAuth_NoLease(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	auth := &logical.Auth{
		ClientToken: root.ID,
		Policies:    []string{"root"},
	}

	te := &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/github/login",
		Policies:    []string{"root"},
		NamespaceID: namespace.RootNamespaceID,
	}
	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should not be able to renew, no expiration
	te = &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/github/login",
		NamespaceID: namespace.RootNamespaceID,
	}
	resp, err := exp.RenewToken(namespace.RootContext(nil), &logical.Request{}, te, 0)
	if err != nil && (err != logical.ErrInvalidRequest || (resp != nil && resp.IsError() && resp.Error().Error() != "lease is not renewable")) {
		t.Fatalf("bad: err:%v resp:%#v", err, resp)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

	// Wait and check token is not invalidated
	time.Sleep(20 * time.Millisecond)

	// Verify token does not get revoked
	out, err := exp.tokenStore.Lookup(namespace.RootContext(nil), root.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("missing token")
	}
}

// Tests both the expiration function and the core function
func TestExpiration_RegisterAuth_NoTTL(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	exp := c.expiration
	ctx := namespace.RootContext(nil)

	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	auth := &logical.Auth{
		ClientToken:   root.ID,
		TokenPolicies: []string{"root"},
	}

	// First on core
	err = c.RegisterAuth(ctx, 0, "auth/github/login", auth, "")
	if err != nil {
		t.Fatal(err)
	}

	auth.TokenPolicies[0] = "default"
	err = c.RegisterAuth(ctx, 0, "auth/github/login", auth, "")
	if err == nil {
		t.Fatal("expected error")
	}

	// Now expiration
	// Should work, root token with zero TTL
	te := &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/github/login",
		Policies:    []string{"root"},
		NamespaceID: namespace.RootNamespaceID,
	}
	err = exp.RegisterAuth(ctx, te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Test non-root token with zero TTL
	te.Policies = []string{"default"}
	err = exp.RegisterAuth(ctx, te, auth, "")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExpiration_Revoke(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Hour,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	id, err := exp.Register(namespace.RootContext(nil), req, resp, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if err := exp.Revoke(namespace.RootContext(nil), id); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = noop.Requests[0]
	if req.Operation != logical.RevokeOperation {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_RevokeOnExpire(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: 20 * time.Millisecond,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	_, err = exp.Register(namespace.RootContext(nil), req, resp, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	start := time.Now()
	for time.Now().Sub(start) < time.Second {
		req = nil

		noop.Lock()
		if len(noop.Requests) > 0 {
			req = noop.Requests[0]
		}
		noop.Unlock()
		if req == nil {
			time.Sleep(5 * time.Millisecond)
			continue
		}
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}

		break
	}
}

func TestExpiration_RevokePrefix(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	paths := []string{
		"prod/aws/foo",
		"prod/aws/sub/bar",
		"prod/aws/zip",
	}
	for _, path := range paths {
		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        path,
			ClientToken: "foobar",
		}
		req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
		resp := &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL: time.Second,
				},
			},
			Data: map[string]interface{}{
				"access_key": "xyz",
				"secret_key": "abcd",
			},
		}
		_, err := exp.Register(namespace.RootContext(nil), req, resp, "")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Should nuke all the keys
	if err := exp.RevokePrefix(namespace.RootContext(nil), "prod/aws/", true); err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(noop.Requests) != 3 {
		t.Fatalf("Bad: %v", noop.Requests)
	}
	for _, req := range noop.Requests {
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}
	}

	expect := []string{
		"foo",
		"sub/bar",
		"zip",
	}
	sort.Strings(noop.Paths)
	sort.Strings(expect)
	if !reflect.DeepEqual(noop.Paths, expect) {
		t.Fatalf("bad: %v", noop.Paths)
	}
}

func TestExpiration_RevokeByToken(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	paths := []string{
		"prod/aws/foo",
		"prod/aws/sub/bar",
		"prod/aws/zip",
	}
	for _, path := range paths {
		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        path,
			ClientToken: "foobarbaz",
		}
		req.SetTokenEntry(&logical.TokenEntry{ID: "foobarbaz", NamespaceID: "root"})
		resp := &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL: time.Second,
				},
			},
			Data: map[string]interface{}{
				"access_key": "xyz",
				"secret_key": "abcd",
			},
		}
		_, err := exp.Register(namespace.RootContext(nil), req, resp, "")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Should nuke all the keys
	te := &logical.TokenEntry{
		ID:          "foobarbaz",
		NamespaceID: namespace.RootNamespaceID,
	}
	if err := exp.RevokeByToken(namespace.RootContext(nil), te); err != nil {
		t.Fatalf("err: %v", err)
	}

	limit := time.Now().Add(3 * time.Second)
	for time.Now().Before(limit) {
		time.Sleep(50 * time.Millisecond)

		noop.Lock()
		currentRequests := len(noop.Requests)
		noop.Unlock()

		if currentRequests == 3 {
			break
		}
	}

	noop.Lock()
	defer noop.Unlock()
	if len(noop.Requests) != 3 {
		t.Errorf("Noop revocation requests less than expected, expected 3, found %d", len(noop.Requests))
	}
	for _, req := range noop.Requests {
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}
	}

	expect := []string{
		"foo",
		"sub/bar",
		"zip",
	}
	sort.Strings(noop.Paths)
	sort.Strings(expect)
	if !reflect.DeepEqual(noop.Paths, expect) {
		t.Fatalf("bad: %v", noop.Paths)
	}
}

func TestExpiration_RevokeByToken_Blocking(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	// Request handle with a timeout context that simulates blocking lease revocation.
	noop.RequestHandler = func(ctx context.Context, req *logical.Request) (*logical.Response, error) {
		ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
		defer cancel()

		select {
		case <-ctx.Done():
			return noop.Response, nil
		}
	}

	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	paths := []string{
		"prod/aws/foo",
		"prod/aws/sub/bar",
		"prod/aws/zip",
	}
	for _, path := range paths {
		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        path,
			ClientToken: "foobarbaz",
		}
		req.SetTokenEntry(&logical.TokenEntry{ID: "foobarbaz", NamespaceID: "root"})
		resp := &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL: 1 * time.Minute,
				},
			},
			Data: map[string]interface{}{
				"access_key": "xyz",
				"secret_key": "abcd",
			},
		}
		_, err := exp.Register(namespace.RootContext(nil), req, resp, "")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Should nuke all the keys
	te := &logical.TokenEntry{
		ID:          "foobarbaz",
		NamespaceID: namespace.RootNamespaceID,
	}
	if err := exp.RevokeByToken(namespace.RootContext(nil), te); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Lock and check that no requests has gone through yet
	noop.Lock()
	if len(noop.Requests) != 0 {
		t.Fatalf("Bad: %v", noop.Requests)
	}
	noop.Unlock()

	// Wait for a bit for timeouts to trigger and pending revocations to go
	// through and then we relock
	time.Sleep(300 * time.Millisecond)

	noop.Lock()
	defer noop.Unlock()

	// Now make sure that all requests have gone through
	if len(noop.Requests) != 3 {
		t.Fatalf("Bad: %v", noop.Requests)
	}
	for _, req := range noop.Requests {
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}
	}

	expect := []string{
		"foo",
		"sub/bar",
		"zip",
	}
	sort.Strings(noop.Paths)
	sort.Strings(expect)
	if !reflect.DeepEqual(noop.Paths, expect) {
		t.Fatalf("bad: %v", noop.Paths)
	}
}

func TestExpiration_RenewToken(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a token
	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
	}

	te := &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/token/login",
		NamespaceID: namespace.RootNamespaceID,
	}
	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Renew the token
	te = &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/token/login",
		NamespaceID: namespace.RootNamespaceID,
	}
	out, err := exp.RenewToken(namespace.RootContext(nil), &logical.Request{}, te, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if auth.ClientToken != out.Auth.ClientToken {
		t.Fatalf("bad: %#v", out)
	}
}

func TestExpiration_RenewToken_period(t *testing.T) {
	exp := mockExpiration(t)
	root := &logical.TokenEntry{
		Policies:     []string{"root"},
		Path:         "auth/token/root",
		DisplayName:  "root",
		CreationTime: time.Now().Unix(),
		Period:       time.Minute,
		NamespaceID:  namespace.RootNamespaceID,
	}
	if err := exp.tokenStore.create(namespace.RootContext(nil), root); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a token
	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
		Period: time.Minute,
	}
	te := &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/token/login",
		NamespaceID: namespace.RootNamespaceID,
	}
	err := exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if exp.leaseCount != 1 {
		t.Fatalf("expected %v leases, got %v", 1, exp.leaseCount)
	}

	// Renew the token
	te = &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/token/login",
		NamespaceID: namespace.RootNamespaceID,
	}
	out, err := exp.RenewToken(namespace.RootContext(nil), &logical.Request{}, te, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if auth.ClientToken != out.Auth.ClientToken {
		t.Fatalf("bad: %#v", out)
	}

	if out.Auth.TTL > time.Minute {
		t.Fatalf("expected TTL to be less than 1 minute, got: %s", out.Auth.TTL)
	}

	if exp.leaseCount != 1 {
		t.Fatalf("expected %v leases, got %v", 1, exp.leaseCount)
	}
}

func TestExpiration_RenewToken_period_backend(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Mount a noop backend
	noop := &NoopBackend{
		Response: &logical.Response{
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL:       10 * time.Second,
					Renewable: true,
				},
				Period: 5 * time.Second,
			},
		},
		DefaultLeaseTTL: 5 * time.Second,
		MaxLeaseTTL:     5 * time.Second,
	}

	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, credentialBarrierPrefix)
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "auth/foo/", &MountEntry{Path: "auth/foo/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	// Register a token
	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL:       10 * time.Second,
			Renewable: true,
		},
		Period: 5 * time.Second,
	}
	te := &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/foo/login",
		NamespaceID: namespace.RootNamespaceID,
	}

	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait 3 seconds
	time.Sleep(3 * time.Second)
	te = &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/foo/login",
		NamespaceID: namespace.RootNamespaceID,
	}
	resp, err := exp.RenewToken(namespace.RootContext(nil), &logical.Request{}, te, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	if resp.Auth.TTL == 0 || resp.Auth.TTL > 5*time.Second {
		t.Fatalf("expected TTL to be greater than zero and less than or equal to period, got: %s", resp.Auth.TTL)
	}

	// Wait another 3 seconds. If period works correctly, this should not fail
	time.Sleep(3 * time.Second)
	resp, err = exp.RenewToken(namespace.RootContext(nil), &logical.Request{}, te, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	if resp.Auth.TTL < 4*time.Second || resp.Auth.TTL > 5*time.Second {
		t.Fatalf("expected TTL to be around period's value, got: %s", resp.Auth.TTL)
	}
}

func TestExpiration_RenewToken_NotRenewable(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a token
	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: false,
		},
	}
	te := &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/foo/login",
		NamespaceID: namespace.RootNamespaceID,
	}
	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to renew the token
	te = &logical.TokenEntry{
		ID:          root.ID,
		Path:        "auth/github/login",
		NamespaceID: namespace.RootNamespaceID,
	}
	resp, err := exp.RenewToken(namespace.RootContext(nil), &logical.Request{}, te, 0)
	if err != nil && (err != logical.ErrInvalidRequest || (resp != nil && resp.IsError() && resp.Error().Error() != "invalid lease ID")) {
		t.Fatalf("bad: err:%v resp:%#v", err, resp)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
}

func TestExpiration_Renew(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       time.Second,
				Renewable: true,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	id, err := exp.Register(namespace.RootContext(nil), req, resp, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	noop.Response = &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Second,
			},
		},
		Data: map[string]interface{}{
			"access_key": "123",
			"secret_key": "abcd",
		},
	}

	out, err := exp.Renew(namespace.RootContext(nil), id, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	noop.Lock()
	defer noop.Unlock()

	if !reflect.DeepEqual(out, noop.Response) {
		t.Fatalf("Bad: %#v", out)
	}

	if len(noop.Requests) != 1 {
		t.Fatalf("Bad: %#v", noop.Requests)
	}
	req = noop.Requests[0]
	if req.Operation != logical.RenewOperation {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_Renew_NotRenewable(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       time.Second,
				Renewable: false,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	id, err := exp.Register(namespace.RootContext(nil), req, resp, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	_, err = exp.Renew(namespace.RootContext(nil), id, 0)
	if err.Error() != "lease is not renewable" {
		t.Fatalf("err: %v", err)
	}

	noop.Lock()
	defer noop.Unlock()

	if len(noop.Requests) != 0 {
		t.Fatalf("Bad: %#v", noop.Requests)
	}
}

func TestExpiration_Renew_RevokeOnExpire(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       20 * time.Millisecond,
				Renewable: true,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	id, err := exp.Register(namespace.RootContext(nil), req, resp, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	noop.Response = &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: 20 * time.Millisecond,
			},
		},
		Data: map[string]interface{}{
			"access_key": "123",
			"secret_key": "abcd",
		},
	}

	_, err = exp.Renew(namespace.RootContext(nil), id, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	start := time.Now()
	for time.Now().Sub(start) < time.Second {
		req = nil

		noop.Lock()
		if len(noop.Requests) >= 2 {
			req = noop.Requests[1]
		}
		noop.Unlock()

		if req == nil {
			time.Sleep(5 * time.Millisecond)
			continue
		}
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}
		break
	}
}

func TestExpiration_Renew_FinalSecond(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       2 * time.Second,
				Renewable: true,
			},
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	ctx := namespace.RootContext(nil)
	id, err := exp.Register(ctx, req, resp, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	le, err := exp.loadEntry(ctx, id)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Give it an auth section to emulate the real world bug
	le.Auth = &logical.Auth{
		LeaseOptions: logical.LeaseOptions{
			Renewable: true,
		},
	}
	exp.persistEntry(ctx, le)

	noop.Response = &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:    2 * time.Second,
				MaxTTL: 2 * time.Second,
			},
		},
		Data: map[string]interface{}{
			"access_key": "123",
			"secret_key": "abcd",
		},
	}

	time.Sleep(1000 * time.Millisecond)
	_, err = exp.Renew(ctx, id, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if _, ok := exp.nonexpiring.Load(id); ok {
		t.Fatalf("expirable lease became nonexpiring")
	}
}

func TestExpiration_Renew_FinalSecond_Lease(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "foobar", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       2 * time.Second,
				Renewable: true,
			},
			LeaseID: "abcde",
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	ctx := namespace.RootContext(nil)
	id, err := exp.Register(ctx, req, resp, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	le, err := exp.loadEntry(ctx, id)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	// Give it an auth section to emulate the real world bug
	le.Auth = &logical.Auth{
		LeaseOptions: logical.LeaseOptions{
			Renewable: true,
		},
	}
	exp.persistEntry(ctx, le)

	time.Sleep(1000 * time.Millisecond)
	_, err = exp.Renew(ctx, id, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if _, ok := exp.nonexpiring.Load(id); ok {
		t.Fatalf("expirable lease became nonexpiring")
	}
}

func TestExpiration_revokeEntry(t *testing.T) {
	exp := mockExpiration(t)

	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "foo/bar/", &MountEntry{Path: "foo/bar/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	le := &leaseEntry{
		LeaseID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
		namespace:  namespace.RootNamespace,
	}

	err = exp.revokeEntry(namespace.RootContext(nil), le)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	noop.Lock()
	defer noop.Unlock()

	req := noop.Requests[0]
	if req.Operation != logical.RevokeOperation {
		t.Fatalf("bad: operation; req: %#v", req)
	}
	if !reflect.DeepEqual(req.Data, le.Data) {
		t.Fatalf("bad: data; req: %#v\n le: %#v\n", req, le)
	}
}

func TestExpiration_revokeEntry_token(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// N.B.: Vault doesn't allow both a secret and auth to be returned, but the
	// reason for both is that auth needs to be included in order to use the
	// token store as it's the only mounted backend, *but* RegisterAuth doesn't
	// actually create the index by token, only Register (for a Secret) does.
	// So without the Secret we don't do anything when removing the index which
	// (at the time of writing) now fails because a bug causing every token
	// expiration to do an extra delete to a non-existent key has been fixed,
	// and this test relies on this nonstandard behavior.
	le := &leaseEntry{
		LeaseID: "foo/bar/1234",
		Auth: &logical.Auth{
			ClientToken: root.ID,
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		ClientToken: root.ID,
		Path:        "foo/bar",
		IssueTime:   time.Now(),
		ExpireTime:  time.Now().Add(time.Minute),
		namespace:   namespace.RootNamespace,
	}

	if err := exp.persistEntry(namespace.RootContext(nil), le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}
	if err := exp.createIndexByToken(namespace.RootContext(nil), le, le.ClientToken); err != nil {
		t.Fatalf("error creating secondary index: %v", err)
	}
	exp.updatePending(le)

	indexEntry, err := exp.indexByToken(namespace.RootContext(nil), le)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if indexEntry == nil {
		t.Fatalf("err: should have found a secondary index entry")
	}

	err = exp.revokeEntry(namespace.RootContext(nil), le)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	limit := time.Now().Add(10 * time.Second)
	for time.Now().Before(limit) {
		indexEntry, err = exp.indexByToken(namespace.RootContext(nil), le)
		if err != nil {
			t.Fatalf("token index lookup error: %v", err)
		}
		if indexEntry == nil {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}

	if indexEntry != nil {
		t.Fatalf("should not have found a secondary index entry after revocation")
	}

	out, err := exp.tokenStore.Lookup(namespace.RootContext(nil), le.ClientToken)
	if err != nil {
		t.Fatalf("error looking up client token after revocation: %v", err)
	}
	if out != nil {
		t.Fatalf("should not have found revoked token in tokenstore: %v", out)
	}
}

func TestExpiration_renewEntry(t *testing.T) {
	exp := mockExpiration(t)

	noop := &NoopBackend{
		Response: &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					Renewable: true,
					TTL:       time.Hour,
				},
			},
			Data: map[string]interface{}{
				"testing": false,
			},
		},
	}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "foo/bar/", &MountEntry{Path: "foo/bar/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	le := &leaseEntry{
		LeaseID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
		namespace:  namespace.RootNamespace,
	}

	resp, err := exp.renewEntry(namespace.RootContext(nil), le, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	noop.Lock()
	defer noop.Unlock()

	if !reflect.DeepEqual(resp, noop.Response) {
		t.Fatalf("bad: %#v", resp)
	}

	req := noop.Requests[0]
	if req.Operation != logical.RenewOperation {
		t.Fatalf("Bad: %v", req)
	}
	if !reflect.DeepEqual(req.Data, le.Data) {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_revokeEntry_rejected_fairsharing(t *testing.T) {
	core, _, _ := TestCoreUnsealed(t)
	exp := core.expiration

	rejected := new(uint32)

	noop := &NoopBackend{
		RequestHandler: func(ctx context.Context, req *logical.Request) (*logical.Response, error) {
			if req.Operation == logical.RevokeOperation {
				if atomic.CompareAndSwapUint32(rejected, 0, 1) {
					t.Logf("denying revocation")
					return nil, errors.New("nope")
				}
				t.Logf("allowing revocation")
			}
			return nil, nil
		},
	}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "foo/bar/", &MountEntry{Path: "foo/bar/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	le := &leaseEntry{
		LeaseID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now().Add(time.Minute),
		namespace:  namespace.RootNamespace,
	}

	err = exp.persistEntry(namespace.RootContext(nil), le)
	if err != nil {
		t.Fatal(err)
	}

	err = exp.LazyRevoke(namespace.RootContext(nil), le.LeaseID)
	if err != nil {
		t.Fatal(err)
	}

	// Give time to let the request be handled
	time.Sleep(1 * time.Second)

	if atomic.LoadUint32(rejected) != 1 {
		t.Fatal("unexpected val for rejected")
	}

	err = exp.Stop()
	if err != nil {
		t.Fatal(err)
	}

	err = core.setupExpiration(expireLeaseStrategyFairsharing)
	if err != nil {
		t.Fatal(err)
	}
	exp = core.expiration

	for {
		if !exp.inRestoreMode() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Now let the revocation actually process
	time.Sleep(1 * time.Second)

	le, err = exp.FetchLeaseTimes(namespace.RootContext(nil), le.LeaseID)
	if err != nil {
		t.Fatal(err)
	}
	if le != nil {
		t.Fatal("lease entry not nil")
	}
}

func TestExpiration_renewAuthEntry(t *testing.T) {
	exp := mockExpiration(t)

	noop := &NoopBackend{
		Response: &logical.Response{
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					Renewable: true,
					TTL:       time.Hour,
				},
			},
		},
	}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "auth/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "auth/foo/", &MountEntry{Path: "auth/foo/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor", namespace: namespace.RootNamespace}, view)
	if err != nil {
		t.Fatal(err)
	}

	le := &leaseEntry{
		LeaseID: "auth/foo/1234",
		Path:    "auth/foo/login",
		Auth: &logical.Auth{
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       time.Minute,
			},
			InternalData: map[string]interface{}{
				"MySecret": "secret",
			},
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now().Add(time.Minute),
		namespace:  namespace.RootNamespace,
	}

	resp, err := exp.renewAuthEntry(namespace.RootContext(nil), &logical.Request{}, le, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	noop.Lock()
	defer noop.Unlock()

	if !reflect.DeepEqual(resp, noop.Response) {
		t.Fatalf("bad: %#v", resp)
	}

	req := noop.Requests[0]
	if req.Operation != logical.RenewOperation {
		t.Fatalf("Bad: %v", req)
	}
	if req.Path != "login" {
		t.Fatalf("Bad: %v", req)
	}
	if req.Auth.InternalData["MySecret"] != "secret" {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_PersistLoadDelete(t *testing.T) {
	exp := mockExpiration(t)
	lastTime := time.Now()
	le := &leaseEntry{
		LeaseID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		IssueTime:       lastTime,
		ExpireTime:      lastTime,
		LastRenewalTime: lastTime,
		namespace:       namespace.RootNamespace,
	}
	if err := exp.persistEntry(namespace.RootContext(nil), le); err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := exp.loadEntry(namespace.RootContext(nil), "foo/bar/1234")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !le.LastRenewalTime.Equal(out.LastRenewalTime) ||
		!le.IssueTime.Equal(out.IssueTime) ||
		!le.ExpireTime.Equal(out.ExpireTime) {
		t.Fatalf("bad: expected:\n%#v\nactual:\n%#v", le, out)
	}
	le.LastRenewalTime = out.LastRenewalTime
	le.IssueTime = out.IssueTime
	le.ExpireTime = out.ExpireTime
	if !reflect.DeepEqual(out, le) {
		t.Fatalf("bad: expected:\n%#v\nactual:\n%#v", le, out)
	}

	err = exp.deleteEntry(namespace.RootContext(nil), le)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err = exp.loadEntry(namespace.RootContext(nil), "foo/bar/1234")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("out: %#v", out)
	}
}

func TestLeaseEntry(t *testing.T) {
	le := &leaseEntry{
		LeaseID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       time.Minute,
				Renewable: true,
			},
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
	}

	enc, err := le.encode()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := decodeLeaseEntry(enc)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(out.Data, le.Data) {
		t.Fatalf("got: %#v, expect %#v", out, le)
	}

	// Test renewability
	le.ExpireTime = time.Time{}
	if r, _ := le.renewable(); r {
		t.Fatal("lease with zero expire time is not renewable")
	}
	le.ExpireTime = time.Now().Add(-1 * time.Hour)
	if r, _ := le.renewable(); r {
		t.Fatal("lease with expire time in the past is not renewable")
	}
	le.ExpireTime = time.Now().Add(1 * time.Hour)
	if r, err := le.renewable(); !r {
		t.Fatalf("lease with future expire time is renewable, err: %v", err)
	}
	le.Secret.LeaseOptions.Renewable = false
	if r, _ := le.renewable(); r {
		t.Fatal("secret is set to not be renewable but returns as renewable")
	}
	le.Secret = nil
	le.Auth = &logical.Auth{
		LeaseOptions: logical.LeaseOptions{
			Renewable: true,
		},
	}
	if r, err := le.renewable(); !r {
		t.Fatalf("auth is renewable but is set to not be, err: %v", err)
	}
	le.Auth.LeaseOptions.Renewable = false
	if r, _ := le.renewable(); r {
		t.Fatal("auth is set to not be renewable but returns as renewable")
	}
}

func TestExpiration_RevokeForce(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	core.logicalBackends["badrenew"] = badRenewFactory
	me := &MountEntry{
		Table:    mountTableType,
		Path:     "badrenew/",
		Type:     "badrenew",
		Accessor: "badrenewaccessor",
	}

	err := core.mount(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "badrenew/creds",
		ClientToken: root,
	}
	err = core.PopulateTokenEntry(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("response was nil")
	}
	if resp.Secret == nil {
		t.Fatalf("response secret was nil, response was %#v", *resp)
	}

	req.Operation = logical.UpdateOperation
	req.Path = "sys/revoke-prefix/badrenew/creds"

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatal("expected error")
	}

	req.Path = "sys/revoke-force/badrenew/creds"
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
}

func TestExpiration_RevokeForceSingle(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	core.logicalBackends["badrenew"] = badRenewFactory
	me := &MountEntry{
		Table:    mountTableType,
		Path:     "badrenew/",
		Type:     "badrenew",
		Accessor: "badrenewaccessor",
	}

	err := core.mount(namespace.RootContext(nil), me)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "badrenew/creds",
		ClientToken: root,
	}
	err = core.PopulateTokenEntry(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("response was nil")
	}
	if resp.Secret == nil {
		t.Fatalf("response secret was nil, response was %#v", *resp)
	}
	leaseID := resp.Secret.LeaseID

	req.Operation = logical.UpdateOperation
	req.Path = "sys/leases/lookup"
	req.Data = map[string]interface{}{"lease_id": leaseID}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.Data["id"].(string) != leaseID {
		t.Fatalf("expected id %q, got %q", leaseID, resp.Data["id"].(string))
	}

	req.Path = "sys/revoke-prefix/" + leaseID

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatal("expected error")
	}

	req.Path = "sys/revoke-force/" + leaseID
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	req.Path = "sys/leases/lookup"
	req.Data = map[string]interface{}{"lease_id": leaseID}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid request") {
		t.Fatalf("bad error: %v", err)
	}
}

func badRenewFactory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	be := &framework.Backend{
		Paths: []*framework.Path{
			{
				Pattern: "creds",
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: func(context.Context, *logical.Request, *framework.FieldData) (*logical.Response, error) {
						resp := &logical.Response{
							Secret: &logical.Secret{
								InternalData: map[string]interface{}{
									"secret_type": "badRenewBackend",
								},
							},
						}
						resp.Secret.TTL = time.Second * 30
						return resp, nil
					},
				},
			},
		},

		Secrets: []*framework.Secret{
			{
				Type: "badRenewBackend",
				Revoke: func(context.Context, *logical.Request, *framework.FieldData) (*logical.Response, error) {
					return nil, fmt.Errorf("always errors")
				},
			},
		},
		BackendType: logical.TypeLogical,
	}

	err := be.Setup(namespace.RootContext(nil), conf)
	if err != nil {
		return nil, err
	}

	return be, nil
}

func sampleToken(t *testing.T, exp *ExpirationManager, path string, expiring bool, policy string) *logical.TokenEntry {
	t.Helper()

	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	auth := &logical.Auth{
		ClientToken: root.ID,
		Policies:    []string{policy},
	}
	if expiring {
		auth.LeaseOptions = logical.LeaseOptions{
			TTL: time.Hour,
		}
	}

	te := &logical.TokenEntry{
		ID:          root.ID,
		Path:        path,
		NamespaceID: namespace.RootNamespaceID,
		Policies:    auth.Policies,
	}

	err = exp.RegisterAuth(namespace.RootContext(nil), te, auth, "")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	return te
}

func findMatchingPath(path string, tokenEntries []*logical.TokenEntry) bool {
	for _, te := range tokenEntries {
		if path == te.Path {
			return true
		}
	}
	return false
}

func findMatchingPolicy(policy string, tokenEntries []*logical.TokenEntry) bool {
	for _, te := range tokenEntries {
		for _, p := range te.Policies {
			if policy == p {
				return true
			}
		}
	}
	return false
}

func TestExpiration_WalkTokens(t *testing.T) {
	exp := mockExpiration(t)

	tokenEntries := []*logical.TokenEntry{
		sampleToken(t, exp, "auth/userpass/login", true, "default"),
		sampleToken(t, exp, "auth/userpass/login", true, "policy23457"),
		sampleToken(t, exp, "auth/token/create", false, "root"),
		sampleToken(t, exp, "auth/github/login", true, "root"),
		sampleToken(t, exp, "auth/github/login", false, "root"),
	}

	waitForRestore(t, exp)

	for true {
		// Count before and after each revocation
		t.Logf("Counting %d tokens.", len(tokenEntries))
		count := 0
		exp.WalkTokens(func(leaseId string, auth *logical.Auth, path string) bool {
			count += 1
			t.Logf("Lease ID %d: %q\n", count, leaseId)
			if !findMatchingPath(path, tokenEntries) {
				t.Errorf("Mismatched Path: %v", path)
			}
			if len(auth.Policies) < 1 || !findMatchingPolicy(auth.Policies[0], tokenEntries) {
				t.Errorf("Mismatched Policies: %v", auth.Policies)
			}
			return true
		})
		if count != len(tokenEntries) {
			t.Errorf("Mismatched number of tokens: %v", count)
		}

		if len(tokenEntries) == 0 {
			break
		}

		// Revoke last token
		toRevoke := len(tokenEntries) - 1
		leaseId, err := exp.CreateOrFetchRevocationLeaseByToken(namespace.RootContext(nil), tokenEntries[toRevoke])
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		t.Logf("revocation lease ID: %q", leaseId)
		err = exp.Revoke(namespace.RootContext(nil), leaseId)
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		tokenEntries = tokenEntries[:len(tokenEntries)-1]

	}
}

func waitForRestore(t *testing.T, exp *ExpirationManager) {
	t.Helper()

	timeout := time.After(200 * time.Millisecond)
	ticker := time.Tick(5 * time.Millisecond)

	for exp.inRestoreMode() {
		select {
		case <-timeout:
			t.Fatalf("Timeout waiting for expiration manager to recover.")
		case <-ticker:
			continue
		}
	}
}

func TestExpiration_CachedPolicyIsShared(t *testing.T) {
	exp := mockExpiration(t)

	tokenEntries := []*logical.TokenEntry{
		sampleToken(t, exp, "auth/userpass/login", true, "policy23457"),
		sampleToken(t, exp, "auth/github/login", true, strings.Join([]string{"policy", "23457"}, "")),
		sampleToken(t, exp, "auth/token/create", true, "policy23457"),
	}

	var policies [][]string

	waitForRestore(t, exp)
	exp.WalkTokens(func(leaseId string, auth *logical.Auth, path string) bool {
		policies = append(policies, auth.Policies)
		return true
	})
	if len(policies) != len(tokenEntries) {
		t.Fatalf("Mismatched number of tokens: %v", len(policies))
	}
	ptrs := make([]*string, len(policies))
	for i := range ptrs {
		ptrs[i] = &((policies[0])[0])
	}
	for i := 1; i < len(ptrs); i++ {
		if ptrs[i-1] != ptrs[i] {
			t.Errorf("Mismatched pointers: %v and %v", ptrs[i-1], ptrs[i])
		}
	}
}

func TestExpiration_FairsharingEnvVar(t *testing.T) {
	testCases := []struct {
		set      string
		expected int
	}{
		{
			set:      "15",
			expected: 15,
		},
		{
			set:      "0",
			expected: numExpirationWorkersTest,
		},
		{
			set:      "10001",
			expected: numExpirationWorkersTest,
		},
	}

	defer os.Unsetenv(fairshareWorkersOverrideVar)
	for _, tc := range testCases {
		os.Setenv(fairshareWorkersOverrideVar, tc.set)
		exp := mockExpiration(t)

		if fairshare.GetNumWorkers(exp.jobManager) != tc.expected {
			t.Errorf("bad worker pool size. expected %d, got %d", tc.expected, fairshare.GetNumWorkers(exp.jobManager))
		}
	}
}

// register one lease ID and return the leaseID
func registerOneLease(t *testing.T, ctx context.Context, exp *ExpirationManager) string {
	t.Helper()

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "irrevocable/lease",
		ClientToken: "sometoken",
	}
	req.SetTokenEntry(&logical.TokenEntry{ID: "sometoken", NamespaceID: "root"})
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: 10 * time.Hour,
			},
		},
	}

	leaseID, err := exp.Register(ctx, req, resp, "")
	if err != nil {
		t.Fatal(err)
	}

	return leaseID
}

func TestExpiration_MarkIrrevocable(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	exp := c.expiration
	ctx := namespace.RootContext(nil)

	leaseID := registerOneLease(t, ctx, exp)
	loadedLE, err := exp.loadEntry(ctx, leaseID)
	if err != nil {
		t.Fatalf("error loading non irrevocable lease: %v", err)
	}

	if loadedLE.isIrrevocable() {
		t.Fatalf("lease is irrevocable and shouldn't be")
	}
	if _, ok := exp.irrevocable.Load(leaseID); ok {
		t.Fatalf("lease included in irrevocable map")
	}
	if _, ok := exp.pending.Load(leaseID); !ok {
		t.Fatalf("lease not included in pending map")
	}

	irrevocableErr := fmt.Errorf("test irrevocable error")

	exp.pendingLock.Lock()
	exp.markLeaseIrrevocable(ctx, loadedLE, irrevocableErr)
	exp.pendingLock.Unlock()

	if !loadedLE.isIrrevocable() {
		t.Fatalf("irrevocable lease is not irrevocable and should be")
	}
	if loadedLE.RevokeErr != irrevocableErr.Error() {
		t.Errorf("irrevocable lease has wrong error message. expected %s, got %s", irrevocableErr.Error(), loadedLE.RevokeErr)
	}
	if _, ok := exp.irrevocable.Load(leaseID); !ok {
		t.Fatalf("irrevocable lease not included in irrevocable map")
	}

	exp.pendingLock.RLock()
	irrevocableLeaseCount := exp.irrevocableLeaseCount
	exp.pendingLock.RUnlock()

	if irrevocableLeaseCount != 1 {
		t.Fatalf("expected 1 irrevocable lease, found %d", irrevocableLeaseCount)
	}
	if _, ok := exp.pending.Load(leaseID); ok {
		t.Fatalf("irrevocable lease included in pending map")
	}
	if _, ok := exp.nonexpiring.Load(leaseID); ok {
		t.Fatalf("irrevocable lease included in nonexpiring map")
	}

	// stop and restore to verify that irrevocable leases are properly loaded from storage
	err = c.stopExpiration()
	if err != nil {
		t.Fatalf("error stopping expiration manager: %v", err)
	}

	err = exp.Restore(nil)
	if err != nil {
		t.Fatalf("error restoring expiration manager: %v", err)
	}

	loadedLE, err = exp.loadEntry(ctx, leaseID)
	if err != nil {
		t.Fatalf("error loading non irrevocable lease after restore: %v", err)
	}
	exp.updatePending(loadedLE)

	if !loadedLE.isIrrevocable() {
		t.Fatalf("irrevocable lease is not irrevocable and should be")
	}
	if loadedLE.RevokeErr != irrevocableErr.Error() {
		t.Errorf("irrevocable lease has wrong error message. expected %s, got %s", irrevocableErr.Error(), loadedLE.RevokeErr)
	}
	if _, ok := exp.irrevocable.Load(leaseID); !ok {
		t.Fatalf("irrevocable lease not included in irrevocable map")
	}
	if _, ok := exp.pending.Load(leaseID); ok {
		t.Fatalf("irrevocable lease included in pending map")
	}
	if _, ok := exp.nonexpiring.Load(leaseID); ok {
		t.Fatalf("irrevocable lease included in nonexpiring map")
	}
}

func TestExpiration_FetchLeaseTimesIrrevocable(t *testing.T) {
	exp := mockExpiration(t)
	ctx := namespace.RootContext(nil)

	leaseID := registerOneLease(t, ctx, exp)
	expectedLeaseTimes, err := exp.FetchLeaseTimes(ctx, leaseID)
	if err != nil {
		t.Fatalf("error getting lease times: %v", err)
	}
	if expectedLeaseTimes == nil {
		t.Fatal("got nil lease")
	}

	le, err := exp.loadEntry(ctx, leaseID)
	if err != nil {
		t.Fatalf("error loading lease: %v", err)
	}
	exp.pendingLock.Lock()
	exp.markLeaseIrrevocable(ctx, le, fmt.Errorf("test irrevocable error"))
	exp.pendingLock.Unlock()

	irrevocableLeaseTimes, err := exp.FetchLeaseTimes(ctx, leaseID)
	if err != nil {
		t.Fatalf("error getting irrevocable lease times: %v", err)
	}
	if irrevocableLeaseTimes == nil {
		t.Fatal("got nil irrevocable lease")
	}

	// strip monotonic clock reading
	expectedLeaseTimes.IssueTime = expectedLeaseTimes.IssueTime.Round(0)
	expectedLeaseTimes.ExpireTime = expectedLeaseTimes.ExpireTime.Round(0)
	expectedLeaseTimes.LastRenewalTime = expectedLeaseTimes.LastRenewalTime.Round(0)

	if !irrevocableLeaseTimes.IssueTime.Equal(expectedLeaseTimes.IssueTime) {
		t.Errorf("bad issue time. expected %v, got %v", expectedLeaseTimes.IssueTime, irrevocableLeaseTimes.IssueTime)
	}
	if !irrevocableLeaseTimes.ExpireTime.Equal(expectedLeaseTimes.ExpireTime) {
		t.Errorf("bad expire time. expected %v, got %v", expectedLeaseTimes.ExpireTime, irrevocableLeaseTimes.ExpireTime)
	}
	if !irrevocableLeaseTimes.LastRenewalTime.Equal(expectedLeaseTimes.LastRenewalTime) {
		t.Errorf("bad last renew time. expected %v, got %v", expectedLeaseTimes.LastRenewalTime, irrevocableLeaseTimes.LastRenewalTime)
	}
}

func TestExpiration_StopClearsIrrevocableCache(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	exp := c.expiration
	ctx := namespace.RootContext(nil)

	leaseID := registerOneLease(t, ctx, exp)
	le, err := exp.loadEntry(ctx, leaseID)
	if err != nil {
		t.Fatalf("error loading non irrevocable lease: %v", err)
	}

	exp.pendingLock.Lock()
	exp.markLeaseIrrevocable(ctx, le, fmt.Errorf("test irrevocable error"))
	exp.pendingLock.Unlock()

	err = c.stopExpiration()
	if err != nil {
		t.Fatalf("error stopping expiration manager: %v", err)
	}

	if _, ok := exp.irrevocable.Load(leaseID); ok {
		t.Error("expiration manager irrevocable cache should be cleared on stop")
	}

	exp.pendingLock.RLock()
	irrevocableLeaseCount := exp.irrevocableLeaseCount
	exp.pendingLock.RUnlock()

	if irrevocableLeaseCount != 0 {
		t.Errorf("expected 0 leases, found %d", irrevocableLeaseCount)
	}
}

func TestExpiration_errorIsUnrecoverable(t *testing.T) {
	testCases := []struct {
		err             error
		isUnrecoverable bool
	}{
		{
			err:             logical.ErrUnrecoverable,
			isUnrecoverable: true,
		},
		{
			err:             logical.ErrUnsupportedOperation,
			isUnrecoverable: true,
		},
		{
			err:             logical.ErrUnsupportedPath,
			isUnrecoverable: true,
		},
		{
			err:             logical.ErrInvalidRequest,
			isUnrecoverable: true,
		},
		{
			err:             logical.ErrPermissionDenied,
			isUnrecoverable: false,
		},
		{
			err:             logical.ErrMultiAuthzPending,
			isUnrecoverable: false,
		},
		{
			err:             fmt.Errorf("some other error"),
			isUnrecoverable: false,
		},
	}

	for _, tc := range testCases {
		out := errIsUnrecoverable(tc.err)
		if out != tc.isUnrecoverable {
			t.Errorf("wrong answer: expected %t, got %t", tc.isUnrecoverable, out)
		}
	}
}

func TestExpiration_unrecoverableErrorMakesIrrevocable(t *testing.T) {
	exp := mockExpiration(t)
	ctx := namespace.RootContext(nil)

	makeJob := func() *revocationJob {
		leaseID := registerOneLease(t, ctx, exp)

		job, err := newRevocationJob(ctx, leaseID, namespace.RootNamespace, exp)
		if err != nil {
			t.Fatalf("err making revocation job: %v", err)
		}

		return job
	}

	testCases := []struct {
		err                 error
		job                 *revocationJob
		shouldBeIrrevocable bool
	}{
		{
			err:                 logical.ErrUnrecoverable,
			job:                 makeJob(),
			shouldBeIrrevocable: true,
		},
		{
			err:                 logical.ErrInvalidRequest,
			job:                 makeJob(),
			shouldBeIrrevocable: true,
		},
		{
			err:                 logical.ErrPermissionDenied,
			job:                 makeJob(),
			shouldBeIrrevocable: false,
		},
		{
			err:                 logical.ErrRateLimitQuotaExceeded,
			job:                 makeJob(),
			shouldBeIrrevocable: false,
		},
		{
			err:                 fmt.Errorf("some random recoverable error"),
			job:                 makeJob(),
			shouldBeIrrevocable: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.err.Error(), func(t *testing.T) {
			tc.job.OnFailure(tc.err)

			le, err := exp.loadEntry(ctx, tc.job.leaseID)
			if err != nil {
				t.Fatalf("could not load leaseID %q: %v", tc.job.leaseID, err)
			}
			if le == nil {
				t.Fatalf("nil lease for leaseID: %q", tc.job.leaseID)
			}

			isIrrevocable := le.isIrrevocable()
			if isIrrevocable != tc.shouldBeIrrevocable {
				t.Errorf("expected irrevocable: %t, got irrevocable: %t", tc.shouldBeIrrevocable, isIrrevocable)
			}
		})
	}
}

func TestExpiration_getIrrevocableLeaseCounts(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	backends := []*backend{
		{
			path: "foo/bar/1/",
			ns:   namespace.RootNamespace,
		},
		{
			path: "foo/bar/2/",
			ns:   namespace.RootNamespace,
		},
		{
			path: "foo/bar/3/",
			ns:   namespace.RootNamespace,
		},
	}
	pathToMount, err := mountNoopBackends(c, backends)
	if err != nil {
		t.Fatal(err)
	}

	exp := c.expiration

	expectedPerMount := 10
	for i := 0; i < expectedPerMount; i++ {
		for _, backend := range backends {
			if _, err := c.AddIrrevocableLease(namespace.RootContext(nil), backend.path); err != nil {
				t.Fatal(err)
			}
		}
	}

	out, err := exp.getIrrevocableLeaseCounts(namespace.RootContext(nil), false)
	if err != nil {
		t.Fatalf("error getting irrevocable lease counts: %v", err)
	}

	exp.pendingLock.RLock()
	irrevocableLeaseCount := exp.irrevocableLeaseCount
	exp.pendingLock.RUnlock()

	if irrevocableLeaseCount != len(backends)*expectedPerMount {
		t.Fatalf("incorrect lease counts. expected %d got %d", len(backends)*expectedPerMount, irrevocableLeaseCount)
	}
	countRaw, ok := out["lease_count"]
	if !ok {
		t.Fatal("no lease count")
	}

	countPerMountRaw, ok := out["counts"]
	if !ok {
		t.Fatal("no count per mount")
	}

	count := countRaw.(int)
	countPerMount := countPerMountRaw.(map[string]int)

	expectedCount := len(backends) * expectedPerMount
	if count != expectedCount {
		t.Errorf("bad count. expected %d, got %d", expectedCount, count)
	}

	if len(countPerMount) != len(backends) {
		t.Fatalf("bad mounts. got %#v, expected %#v", countPerMount, backends)
	}

	for _, backend := range backends {
		mountCount := countPerMount[pathToMount[backend.path]]
		if mountCount != expectedPerMount {
			t.Errorf("bad count for prefix %q. expected %d, got %d", backend.path, expectedPerMount, mountCount)
		}
	}
}

func TestExpiration_listIrrevocableLeases(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	backends := []*backend{
		{
			path: "foo/bar/1/",
			ns:   namespace.RootNamespace,
		},
		{
			path: "foo/bar/2/",
			ns:   namespace.RootNamespace,
		},
		{
			path: "foo/bar/3/",
			ns:   namespace.RootNamespace,
		},
	}
	pathToMount, err := mountNoopBackends(c, backends)
	if err != nil {
		t.Fatal(err)
	}

	exp := c.expiration

	expectedLeases := make([]*basicLeaseTestInfo, 0)
	expectedPerMount := 10
	for i := 0; i < expectedPerMount; i++ {
		for _, backend := range backends {
			le, err := c.AddIrrevocableLease(namespace.RootContext(nil), backend.path)
			if err != nil {
				t.Fatal(err)
			}
			expectedLeases = append(expectedLeases, &basicLeaseTestInfo{
				id:     le.id,
				mount:  pathToMount[backend.path],
				expire: le.expire,
			})
		}
	}

	out, warn, err := exp.listIrrevocableLeases(namespace.RootContext(nil), false, false, MaxIrrevocableLeasesToReturn)
	if err != nil {
		t.Fatalf("error listing irrevocable leases: %v", err)
	}
	if warn != "" {
		t.Errorf("expected no warning, got %q", warn)
	}

	countRaw, ok := out["lease_count"]
	if !ok {
		t.Fatal("no lease count")
	}

	leasesRaw, ok := out["leases"]
	if !ok {
		t.Fatal("no leases")
	}

	count := countRaw.(int)
	leases := leasesRaw.([]*leaseResponse)

	expectedCount := len(backends) * expectedPerMount
	if count != expectedCount {
		t.Errorf("bad count. expected %d, got %d", expectedCount, count)
	}
	if len(leases) != len(expectedLeases) {
		t.Errorf("bad lease results. expected %d, got %d with values %v", len(expectedLeases), len(leases), leases)
	}

	// `leases` is already sorted by lease ID
	sort.Slice(expectedLeases, func(i, j int) bool {
		return expectedLeases[i].id < expectedLeases[j].id
	})
	sort.SliceStable(expectedLeases, func(i, j int) bool {
		return expectedLeases[i].expire.Before(expectedLeases[j].expire)
	})

	for i, lease := range expectedLeases {
		if lease.id != leases[i].LeaseID {
			t.Errorf("bad lease id. expected %q, got %q", lease.id, leases[i].LeaseID)
		}
		if lease.mount != leases[i].MountID {
			t.Errorf("bad mount id. expected %q, got %q", lease.mount, leases[i].MountID)
		}
	}
}

func TestExpiration_listIrrevocableLeases_includeAll(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	exp := c.expiration

	expectedNumLeases := MaxIrrevocableLeasesToReturn + 10
	for i := 0; i < expectedNumLeases; i++ {
		if _, err := c.AddIrrevocableLease(namespace.RootContext(nil), "foo/"); err != nil {
			t.Fatal(err)
		}
	}

	dataRaw, warn, err := exp.listIrrevocableLeases(namespace.RootContext(nil), false, false, MaxIrrevocableLeasesToReturn)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if warn != MaxIrrevocableLeasesWarning {
		t.Errorf("expected warning %q, got %q", MaxIrrevocableLeasesWarning, warn)
	}
	if dataRaw == nil {
		t.Fatal("expected partial data, got nil")
	}

	leaseListLength := len(dataRaw["leases"].([]*leaseResponse))
	if leaseListLength != MaxIrrevocableLeasesToReturn {
		t.Fatalf("expected %d results, got %d", MaxIrrevocableLeasesToReturn, leaseListLength)
	}

	dataRaw, warn, err = exp.listIrrevocableLeases(namespace.RootContext(nil), false, true, 0)
	if err != nil {
		t.Fatalf("got error when using limit=none: %v", err)
	}
	if warn != "" {
		t.Errorf("expected no warning, got %q", warn)
	}
	if dataRaw == nil {
		t.Fatalf("got nil data when using limit=none")
	}

	leaseListLength = len(dataRaw["leases"].([]*leaseResponse))
	if leaseListLength != expectedNumLeases {
		t.Fatalf("expected %d results, got %d", MaxIrrevocableLeasesToReturn, expectedNumLeases)
	}

	numLeasesRaw, ok := dataRaw["lease_count"]
	if !ok {
		t.Fatalf("lease count data not present")
	}
	if numLeasesRaw == nil {
		t.Fatalf("nil lease count")
	}

	numLeases := numLeasesRaw.(int)
	if numLeases != expectedNumLeases {
		t.Errorf("bad lease count. expected %d, got %d", expectedNumLeases, numLeases)
	}
}
