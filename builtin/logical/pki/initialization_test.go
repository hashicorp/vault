// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errInitializationTestStorage = errors.New("temporary storage failure")

// Exercise actual Initialize, routed API requests, and periodic rollback calls.
// None of the recovery tests sets the backend's in-memory storage version.
func TestPKIInitializationRecovery(t *testing.T) {
	for _, test := range []struct {
		name       string
		key        string
		occurrence int
	}{
		{"crl-config", "config/crl", 0},
		{"initial-migration-status", legacyMigrationBundleLogKey, 0},
		{"final-migration-status", legacyMigrationBundleLogKey, 3},
		{"acme-config", storageAcmeConfig, 0},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				storage, ca, material := initializedPKIStorage(t)
				engineStarts := storage.count("list", acmeValidationPrefix)
				var calls int
				storage.setHook(func(_ context.Context, operation, key string) error {
					if operation == "get" && key == test.key {
						calls++
						if test.occurrence == 0 || calls == test.occurrence {
							return errInitializationTestStorage
						}
					}
					return nil
				})
				b, _ := newInitializationTestBackend(t, storage, nil, nil)
				initCtx, cancel := context.WithCancel(context.Background())
				require.ErrorIs(t, b.Initialize(initCtx, &logical.InitializationRequest{Storage: storage}), errInitializationTestStorage)
				cancel() // Reload and plugin RPC initialization contexts can end here.
				assertInitializationUnavailable(t, b, storage)

				storage.setHook(nil)
				reads := storage.count("get", test.key)
				require.NoError(t, initializationRollback(context.Background(), b, storage))
				assert.Equal(t, reads, storage.count("get", test.key), "respect the initial backoff")

				time.Sleep(6 * time.Minute)
				retryCtx := namespace.ContextWithNamespace(context.Background(), &namespace.Namespace{ID: "test", Path: "example/"})
				storage.setHook(func(ctx context.Context, operation, key string) error {
					if operation == "get" && key == legacyMigrationBundleLogKey {
						ns, err := namespace.FromContext(ctx)
						require.NoError(t, err)
						require.Equal(t, "test", ns.ID, "retry uses the current periodic context")
					}
					return nil
				})
				require.NoError(t, initializationRollback(retryCtx, b, storage))
				require.Equal(t, ca, initializationCA(t, b, storage))
				synctest.Wait()
				require.Equal(t, engineStarts+1, storage.count("list", acmeValidationPrefix), "successful recovery starts the ACME engine once")
				require.Equal(t, material, initializationMaterial(t, storage), "recovery must not replace CA, issuer, key, or migration log")
				storage.setHook(nil)
				reads = storage.count("get", legacyMigrationBundleLogKey)
				time.Sleep(6 * time.Minute)
				require.NoError(t, initializationRollback(context.Background(), b, storage))
				require.Equal(t, reads, storage.count("get", legacyMigrationBundleLogKey), "successful initialization is not retried")
			})
		})
	}
}

func TestPKIInitializationRetryLimit(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		storage := newInitializationTestStorage()
		logs := new(initializationTestLog)
		b, _ := newInitializationTestBackend(t, storage, nil, logs)
		storage.failRead("config/crl")
		require.Error(t, b.Initialize(context.Background(), &logical.InitializationRequest{Storage: storage}))
		for i, delay := range []time.Duration{30 * time.Second, time.Minute, 2 * time.Minute} {
			time.Sleep(delay - time.Nanosecond)
			require.NoError(t, initializationRollback(context.Background(), b, storage))
			require.Equal(t, i+1, storage.count("get", "config/crl"), "no retry before the exponential backoff")
			time.Sleep(delay/2 + time.Nanosecond)
			require.ErrorIs(t, initializationRollback(context.Background(), b, storage), errInitializationTestStorage)
			assertInitializationUnavailable(t, b, storage)
		}
		require.Equal(t, 4, storage.count("get", "config/crl"), "one initial attempt and three automatic retries")
		require.Contains(t, logs.String(), "automatic recovery stopped")
		require.Contains(t, logs.String(), "reloading the mount")
		storage.setHook(nil)
		for range 4 {
			time.Sleep(time.Hour)
			require.NoError(t, initializationRollback(context.Background(), b, storage))
		}
		require.Equal(t, 4, storage.count("get", "config/crl"), "exhaustion stops all automatic storage work")
		assertInitializationUnavailable(t, b, storage)
		reloaded, _ := newInitializationTestBackend(t, storage, nil, nil)
		require.NoError(t, reloaded.Initialize(context.Background(), &logical.InitializationRequest{Storage: storage}), "an operator reload gets a fresh initialization budget")
	})
}

func TestPKIInitializationRetryConcurrency(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		healthyStorage, ca, _ := initializedPKIStorage(t)
		healthy, _ := newInitializationTestBackend(t, healthyStorage, nil, nil)
		require.NoError(t, healthy.Initialize(context.Background(), &logical.InitializationRequest{Storage: healthyStorage}))
		var backends []*backend
		var stores []*initializationTestStorage
		release := make(chan struct{})
		for range 3 {
			storage := newInitializationTestStorage()
			b, _ := newInitializationTestBackend(t, storage, nil, nil)
			storage.failRead("config/crl")
			require.Error(t, b.Initialize(context.Background(), &logical.InitializationRequest{Storage: storage}))
			storage.setHook(func(ctx context.Context, operation, key string) error {
				if operation != "get" || key != "config/crl" {
					return nil
				}
				select {
				case <-release:
					return errInitializationTestStorage
				case <-ctx.Done():
					return ctx.Err()
				}
			})
			backends = append(backends, b)
			stores = append(stores, storage)
		}
		time.Sleep(6 * time.Minute)
		results := make(chan error, 2)
		for i := range 2 {
			go func() { results <- initializationRollback(context.Background(), backends[i], stores[i]) }()
		}
		synctest.Wait()
		require.Equal(t, 2, stores[0].count("get", "config/crl"))
		require.Equal(t, 2, stores[1].count("get", "config/crl"))
		require.Equal(t, ca, initializationCA(t, healthy, healthyStorage), "unhealthy mounts do not gate healthy mounts")
		require.NoError(t, initializationRollback(context.Background(), healthy, healthyStorage), "healthy maintenance does not need a recovery permit")
		for range 5 {
			require.NoError(t, initializationRollback(context.Background(), backends[0], stores[0]), "same-mount retries do not overlap")
			require.NoError(t, initializationRollback(context.Background(), backends[2], stores[2]), "a third mount never waits for the global limit")
		}
		require.Equal(t, 2, stores[0].count("get", "config/crl"))
		require.Equal(t, 1, stores[2].count("get", "config/crl"))
		close(release)
		for range 2 {
			require.ErrorIs(t, <-results, errInitializationTestStorage)
		}
		require.ErrorIs(t, initializationRollback(context.Background(), backends[2], stores[2]), errInitializationTestStorage, "failed attempts release global capacity")
		require.Equal(t, 2, stores[2].count("get", "config/crl"))
	})
}

func TestPKIInitializationRetryCancellation(t *testing.T) {
	for _, cleanup := range []bool{false, true} {
		t.Run(map[bool]string{false: "periodic-context", true: "cleanup"}[cleanup], func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				storage := newInitializationTestStorage()
				b, closeBackend := newInitializationTestBackend(t, storage, nil, nil)
				storage.failRead("config/crl")
				require.Error(t, b.Initialize(context.Background(), &logical.InitializationRequest{Storage: storage}))
				storage.setHook(func(ctx context.Context, operation, key string) error {
					if operation == "get" && key == "config/crl" {
						<-ctx.Done()
						return ctx.Err()
					}
					return nil
				})
				time.Sleep(6 * time.Minute)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				done := make(chan error, 1)
				go func() { done <- initializationRollback(ctx, b, storage) }()
				synctest.Wait()
				if cleanup {
					closeBackend()
				} else {
					cancel()
				}
				require.ErrorIs(t, <-done, context.Canceled)
				reads := storage.count("get", "config/crl")
				time.Sleep(6 * time.Minute)
				require.ErrorIs(t, initializationRollback(ctx, b, storage), context.Canceled)
				require.Equal(t, reads, storage.count("get", "config/crl"))
			})
		})
	}
}

func TestPKIInitializationRetryDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		storage := newInitializationTestStorage()
		b, _ := newInitializationTestBackend(t, storage, nil, nil)
		storage.failRead("config/crl")
		require.Error(t, b.Initialize(context.Background(), &logical.InitializationRequest{Storage: storage}))
		storage.setHook(func(ctx context.Context, operation, key string) error {
			if operation == "get" && key == "config/crl" {
				<-ctx.Done()
				return ctx.Err()
			}
			return nil
		})
		time.Sleep(6 * time.Minute)
		start := time.Now()
		require.ErrorIs(t, initializationRollback(context.Background(), b, storage), context.DeadlineExceeded)
		require.Equal(t, 30*time.Second, time.Since(start))
		assertInitializationUnavailable(t, b, storage)
	})
}

func TestPKIInitializationACMEContext(t *testing.T) {
	for _, recover := range []bool{false, true} {
		t.Run(map[bool]string{false: "initial", true: "recovered"}[recover], func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				storage := newInitializationTestStorage()
				b, closeBackend := newInitializationTestBackend(t, storage, nil, nil)
				ctx, cancel := context.WithCancel(namespace.ContextWithNamespace(context.Background(), &namespace.Namespace{ID: "test", Path: "example/"}))
				defer cancel()
				if recover {
					storage.failRead("config/crl")
					require.Error(t, b.Initialize(ctx, &logical.InitializationRequest{Storage: storage}))
					time.Sleep(6 * time.Minute)
				}
				loaded := make(chan context.Context, 1)
				storage.setHook(func(ctx context.Context, operation, key string) error {
					if operation == "list" && key == acmeValidationPrefix {
						loaded <- ctx
						<-ctx.Done()
						return ctx.Err()
					}
					return nil
				})
				if recover {
					require.NoError(t, initializationRollback(ctx, b, storage))
				} else {
					require.NoError(t, b.Initialize(ctx, &logical.InitializationRequest{Storage: storage}))
				}
				workerContext := <-loaded
				require.False(t, b.UseLegacyBundleCaStorage())
				require.Equal(t, 1, storage.count("put", legacyMigrationBundleLogKey), "a retry may perform the first migration when the original call never reached it")
				cancel()
				synctest.Wait()
				require.NoError(t, workerContext.Err(), "initial queue loading outlives the caller and prerequisite deadline")
				ns, err := namespace.FromContext(workerContext)
				require.NoError(t, err)
				require.Equal(t, "test", ns.ID)
				closeBackend()
				synctest.Wait()
				require.ErrorIs(t, workerContext.Err(), context.Canceled, "cleanup cancels asynchronous startup")
				require.Equal(t, 1, storage.count("list", acmeValidationPrefix), "only one ACME engine starts")
			})
		})
	}
}

func TestPKIInitializationOptionalCountTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		storage, ca, _ := initializedPKIStorage(t)
		entry, err := logical.StorageEntryJSON(autoTidyConfigPath, &tidyConfig{MaintainCount: true})
		require.NoError(t, err)
		require.NoError(t, storage.Put(context.Background(), entry))
		b, _ := newInitializationTestBackend(t, storage, nil, nil)
		storage.failRead("config/crl")
		require.Error(t, b.Initialize(context.Background(), &logical.InitializationRequest{Storage: storage}))
		storage.setHook(func(ctx context.Context, operation, key string) error {
			if operation == "list" && key == issuing.PathCerts {
				<-ctx.Done()
				return ctx.Err()
			}
			return nil
		})
		time.Sleep(6 * time.Minute)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		_ = initializationRollback(ctx, b, storage) // Ordinary maintenance can observe caller cancellation.
		require.ErrorIs(t, b.GetCertificateCounter().Error(), context.DeadlineExceeded)
		require.Equal(t, ca, initializationCA(t, b, storage), "optional metrics must not leave a healthy mount unavailable")
		storage.setHook(nil)
		reads := storage.count("get", legacyMigrationBundleLogKey)
		time.Sleep(6 * time.Minute)
		require.NoError(t, initializationRollback(context.Background(), b, storage))
		require.Equal(t, reads, storage.count("get", legacyMigrationBundleLogKey), "a soft count timeout never arms another initializer")
	})
}

func TestPKIInitializationLegacySecondary(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		storage := newInitializationTestStorage()
		system := logical.TestSystemView()
		system.ReplicationStateVal = consts.ReplicationPerformanceSecondary
		b, _ := newInitializationTestBackend(t, storage, system, nil)
		require.NoError(t, b.Initialize(context.Background(), &logical.InitializationRequest{Storage: storage}))
		require.True(t, b.UseLegacyBundleCaStorage())
		reads := storage.count("get", legacyMigrationBundleLogKey)
		for range 4 {
			time.Sleep(6 * time.Minute)
			require.NoError(t, initializationRollback(context.Background(), b, storage))
		}
		require.Equal(t, reads, storage.count("get", legacyMigrationBundleLogKey))
		require.Zero(t, storage.count("put", legacyMigrationBundleLogKey), "a legitimate legacy replica is never migrated by recovery")
		response, err := b.HandleRequest(context.Background(), &logical.Request{Operation: logical.ReadOperation, Path: "ca/pem", Storage: storage})
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, response.Data[logical.HTTPStatusCode])
	})
}

func TestPKIInitializationPreservesForwarding(t *testing.T) {
	for _, test := range []struct {
		name    string
		state   consts.ReplicationState
		local   bool
		forward bool
	}{
		{"performance-standby", consts.ReplicationPerformanceStandby, false, true},
		{"performance-secondary", consts.ReplicationPerformanceSecondary, false, true},
		{"local-secondary", consts.ReplicationPerformanceSecondary, true, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			storage := newInitializationTestStorage()
			system := logical.TestSystemView()
			system.ReplicationStateVal = test.state
			system.LocalMountVal = test.local
			b, _ := newInitializationTestBackend(t, storage, system, nil)
			storage.failRead("config/crl")
			require.Error(t, b.Initialize(context.Background(), &logical.InitializationRequest{Storage: storage}))
			request := &logical.Request{Operation: logical.UpdateOperation, Path: "config/crl", Storage: storage}
			found, exists, err := b.HandleExistenceCheck(context.Background(), request)
			require.NoError(t, err)
			require.False(t, found)
			require.False(t, exists)
			_, err = b.HandleRequest(context.Background(), request)
			if test.forward {
				require.ErrorIs(t, err, logical.ErrReadOnly)
			} else {
				var coded logical.HTTPCodedError
				require.ErrorAs(t, err, &coded)
				require.Equal(t, http.StatusServiceUnavailable, coded.Code())
			}
			if test.state.HasState(consts.ReplicationPerformanceStandby) {
				for _, path := range []string{"issue/example-role", "sign/example-role", "ca/pem"} {
					operation := logical.Operation(logical.UpdateOperation)
					if path == "ca/pem" {
						operation = logical.ReadOperation
					}
					_, err := b.HandleRequest(context.Background(), &logical.Request{Operation: operation, Path: path, Storage: storage})
					require.ErrorIs(t, err, logical.ErrReadOnly, "an unready standby forwards dynamically handled requests too")
				}
			} else {
				assertInitializationUnavailable(t, b, storage)
			}
		})
	}
}

func TestPKIInitializationDoesNotReplayPartialMigration(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		storage := newInitializationTestStorage()
		logs := new(initializationTestLog)
		b, _ := newInitializationTestBackend(t, storage, nil, logs)
		bundle := genCertBundle(t, b, storage)
		entry, err := logical.StorageEntryJSON(legacyCertBundlePath, bundle)
		require.NoError(t, err)
		require.NoError(t, storage.Put(context.Background(), entry))
		storage.setHook(func(_ context.Context, operation, key string) error {
			if operation == "put" && key == storageIssuerConfig {
				return errInitializationTestStorage
			}
			return nil
		})
		require.ErrorIs(t, b.Initialize(context.Background(), &logical.InitializationRequest{Storage: storage}), errInitializationTestStorage)
		issuers, err := storage.Storage.List(context.Background(), issuerPrefix)
		require.NoError(t, err)
		require.Len(t, issuers, 1, "the issuer write succeeded before the default-issuer write failed")
		require.Zero(t, storage.count("put", legacyMigrationBundleLogKey))
		assertInitializationUnavailable(t, b, storage)
		operations := storage.total()
		storage.setHook(nil)
		for range 4 {
			time.Sleep(6 * time.Minute)
			require.NoError(t, initializationRollback(context.Background(), b, storage))
		}
		require.Equal(t, operations, storage.total(), "partial migration is not automatically replayed")
		assertInitializationUnavailable(t, b, storage)
		require.Contains(t, logs.String(), "inspect the error before reloading")
	})
}

func newInitializationTestBackend(t *testing.T, storage logical.Storage, system *logical.StaticSystemView, logs *initializationTestLog) (*backend, func()) {
	t.Helper()
	config := logical.TestBackendConfig()
	config.StorageView = storage
	if system != nil {
		config.System = system
	}
	if logs != nil {
		config.Logger = hclog.New(&hclog.LoggerOptions{Output: logs, Level: hclog.Trace})
	}
	b := Backend(config)
	require.NoError(t, b.Setup(context.Background(), config))
	var once sync.Once
	cleanup := func() { once.Do(func() { b.Cleanup(context.Background()) }) }
	t.Cleanup(cleanup)
	return b, cleanup
}

func initializedPKIStorage(t *testing.T) (*initializationTestStorage, []byte, map[string][]byte) {
	t.Helper()
	storage := newInitializationTestStorage()
	b, cleanup := newInitializationTestBackend(t, storage, nil, nil)
	require.NoError(t, b.Initialize(context.Background(), &logical.InitializationRequest{Storage: storage}))
	response, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation, Path: "root/generate/internal", Storage: storage,
		Data: map[string]interface{}{"common_name": "root.example", "key_type": "ec", "key_bits": 256, "ttl": "24h"},
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.False(t, response.IsError(), "%v", response)
	ca := initializationCA(t, b, storage)
	material := initializationMaterial(t, storage)
	cleanup()
	synctest.Wait()
	return storage, ca, material
}

func initializationCA(t *testing.T, b *backend, storage logical.Storage) []byte {
	t.Helper()
	response, err := b.HandleRequest(context.Background(), &logical.Request{Operation: logical.ReadOperation, Path: "ca/pem", Storage: storage})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, http.StatusOK, response.Data[logical.HTTPStatusCode])
	ca, ok := response.Data[logical.HTTPRawBody].([]byte)
	require.True(t, ok)
	require.NotEmpty(t, ca)
	return ca
}

func initializationMaterial(t *testing.T, storage *initializationTestStorage) map[string][]byte {
	t.Helper()
	keys := []string{storageIssuerConfig, storageKeyConfig, legacyMigrationBundleLogKey}
	for _, prefix := range []string{issuerPrefix, keyPrefix} {
		ids, err := storage.Storage.List(context.Background(), prefix)
		require.NoError(t, err)
		require.Len(t, ids, 1)
		for _, id := range ids {
			keys = append(keys, prefix+id)
		}
	}
	material := make(map[string][]byte)
	for _, key := range keys {
		entry, err := storage.Storage.Get(context.Background(), key)
		require.NoError(t, err)
		require.NotNil(t, entry)
		material[key] = bytes.Clone(entry.Value)
	}
	return material
}

func assertInitializationUnavailable(t *testing.T, b *backend, storage logical.Storage) {
	t.Helper()
	for _, request := range []*logical.Request{
		{Operation: logical.ReadOperation, Path: "ca/pem", Storage: storage},
		{Operation: logical.ListOperation, Path: "issuers", Storage: storage},
	} {
		_, err := b.HandleRequest(context.Background(), request)
		var coded logical.HTTPCodedError
		if assert.ErrorAs(t, err, &coded, "failed initialization must be explicitly unavailable") {
			assert.Equal(t, http.StatusServiceUnavailable, coded.Code())
		}
	}
}

func initializationRollback(ctx context.Context, b *backend, storage logical.Storage) error {
	_, err := b.HandleRequest(ctx, &logical.Request{Operation: logical.RollbackOperation, Storage: storage})
	return err
}

type initializationTestStorage struct {
	logical.Storage
	mu    sync.Mutex
	hook  func(context.Context, string, string) error
	calls [][2]string
}

func newInitializationTestStorage() *initializationTestStorage {
	return &initializationTestStorage{Storage: new(logical.InmemStorage)}
}

func (s *initializationTestStorage) setHook(hook func(context.Context, string, string) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hook = hook
}

func (s *initializationTestStorage) failRead(failedKey string) {
	s.setHook(func(_ context.Context, operation, key string) error {
		if operation == "get" && key == failedKey {
			return errInitializationTestStorage
		}
		return nil
	})
}

func (s *initializationTestStorage) before(ctx context.Context, operation, key string) error {
	s.mu.Lock()
	s.calls = append(s.calls, [2]string{operation, key})
	hook := s.hook
	s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if hook != nil {
		return hook(ctx, operation, key)
	}
	return nil
}

func (s *initializationTestStorage) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	if err := s.before(ctx, "get", key); err != nil {
		return nil, err
	}
	return s.Storage.Get(ctx, key)
}

func (s *initializationTestStorage) Put(ctx context.Context, entry *logical.StorageEntry) error {
	if err := s.before(ctx, "put", entry.Key); err != nil {
		return err
	}
	return s.Storage.Put(ctx, entry)
}

func (s *initializationTestStorage) List(ctx context.Context, prefix string) ([]string, error) {
	if err := s.before(ctx, "list", prefix); err != nil {
		return nil, err
	}
	return s.Storage.List(ctx, prefix)
}

func (s *initializationTestStorage) Delete(ctx context.Context, key string) error {
	if err := s.before(ctx, "delete", key); err != nil {
		return err
	}
	return s.Storage.Delete(ctx, key)
}

func (s *initializationTestStorage) count(operation, key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	var count int
	for _, call := range s.calls {
		if call == [2]string{operation, key} {
			count++
		}
	}
	return count
}

func (s *initializationTestStorage) total() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.calls)
}

type initializationTestLog struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *initializationTestLog) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *initializationTestLog) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}
