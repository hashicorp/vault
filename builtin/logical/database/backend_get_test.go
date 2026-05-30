// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

func newSystemViewWrapper(view logical.SystemView) logical.SystemView {
	return &systemViewWrapper{
		SystemView: view,
	}
}

type systemViewWrapper struct {
	logical.SystemView
	pluginName     string
	builtinFactory func() (interface{}, error)
}

var _ logical.ExtendedSystemView = (*systemViewWrapper)(nil)

func (s *systemViewWrapper) RequestWellKnownRedirect(ctx context.Context, src, dest string) error {
	panic("nope")
}

func (s *systemViewWrapper) DeregisterWellKnownRedirect(ctx context.Context, src string) bool {
	panic("nope")
}

func (s *systemViewWrapper) Auditor() logical.Auditor {
	panic("nope")
}

func (s *systemViewWrapper) ForwardGenericRequest(ctx context.Context, request *logical.Request) (*logical.Response, error) {
	panic("nope")
}

func (s *systemViewWrapper) APILockShouldBlockRequest() (bool, error) {
	panic("nope")
}

func (s *systemViewWrapper) GetPinnedPluginVersion(ctx context.Context, pluginType consts.PluginType, pluginName string) (*pluginutil.PinnedVersion, error) {
	return nil, pluginutil.ErrPinnedVersionNotFound
}

func (s *systemViewWrapper) LookupPluginVersion(ctx context.Context, pluginName string, pluginType consts.PluginType, version string) (*pluginutil.PluginRunner, error) {
	name := s.pluginName
	if name == "" {
		name = mockv5
	}

	factory := s.builtinFactory
	if factory == nil {
		factory = New
	}

	return &pluginutil.PluginRunner{
		Name:           name,
		Type:           consts.PluginTypeDatabase,
		Builtin:        true,
		BuiltinFactory: factory,
	}, nil
}

func getDbBackend(t *testing.T) (*databaseBackend, logical.Storage) {
	t.Helper()
	config := logical.TestBackendConfig()
	config.System = newSystemViewWrapper(config.System)
	config.StorageView = &logical.InmemStorage{}
	// Create and init the backend ourselves instead of using a Factory because
	// the factory function kicks off threads that cause racy tests.
	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	b.schedule = &TestSchedule{}
	b.credRotationQueue = queue.New()
	b.populateQueue(context.Background(), config.StorageView)

	return b, config.StorageView
}

type blockingInitializeDatabase struct {
	initializeDone chan struct{}
}

func newBlockingInitializeDatabase() (interface{}, error) {
	return &blockingInitializeDatabase{initializeDone: make(chan struct{})}, nil
}

func (d *blockingInitializeDatabase) Initialize(context.Context, v5.InitializeRequest) (v5.InitializeResponse, error) {
	<-d.initializeDone
	return v5.InitializeResponse{}, nil
}

func (d *blockingInitializeDatabase) NewUser(context.Context, v5.NewUserRequest) (v5.NewUserResponse, error) {
	return v5.NewUserResponse{}, nil
}

func (d *blockingInitializeDatabase) UpdateUser(context.Context, v5.UpdateUserRequest) (v5.UpdateUserResponse, error) {
	return v5.UpdateUserResponse{}, nil
}

func (d *blockingInitializeDatabase) DeleteUser(context.Context, v5.DeleteUserRequest) (v5.DeleteUserResponse, error) {
	return v5.DeleteUserResponse{}, nil
}

func (d *blockingInitializeDatabase) Type() (string, error) {
	return mockV5Type, nil
}

func (d *blockingInitializeDatabase) Close() error {
	close(d.initializeDone)
	return nil
}

// slowCloseDatabase blocks in Close until closeCh is closed, so a synchronous
// call to closeDatabaseWrapperAfterInitError would stall the test.
type slowCloseDatabase struct {
	closeCh chan struct{}
}

func (d *slowCloseDatabase) Initialize(context.Context, v5.InitializeRequest) (v5.InitializeResponse, error) {
	return v5.InitializeResponse{}, nil
}

func (d *slowCloseDatabase) NewUser(context.Context, v5.NewUserRequest) (v5.NewUserResponse, error) {
	return v5.NewUserResponse{}, nil
}

func (d *slowCloseDatabase) UpdateUser(context.Context, v5.UpdateUserRequest) (v5.UpdateUserResponse, error) {
	return v5.UpdateUserResponse{}, nil
}

func (d *slowCloseDatabase) DeleteUser(context.Context, v5.DeleteUserRequest) (v5.DeleteUserResponse, error) {
	return v5.DeleteUserResponse{}, nil
}
func (d *slowCloseDatabase) Type() (string, error) { return "slow-close", nil }
func (d *slowCloseDatabase) Close() error {
	<-d.closeCh
	return nil
}

// TestCloseDatabaseWrapperAfterInitError_ContextCanceled_IsAsync verifies that
// closeDatabaseWrapperAfterInitError does not block when the error is
// context.Canceled (parent context fired before Vault's own databaseInitTimeout).
// In that path the init goroutine may still be running, so Close must be async.
// context.DeadlineExceeded is symmetric (same code branch), so a single test for
// context.Canceled is sufficient.
func TestCloseDatabaseWrapperAfterInitError_ContextCanceled_IsAsync(t *testing.T) {
	closeCh := make(chan struct{})
	dbw := databaseVersionWrapper{v5: &slowCloseDatabase{closeCh: closeCh}}

	b := &databaseBackend{}

	done := make(chan struct{})
	go func() {
		defer close(done)
		b.closeDatabaseWrapperAfterInitError(dbw, context.Canceled)
	}()

	select {
	case <-done:
		// Good: function returned before Close() completed.
	case <-time.After(100 * time.Millisecond):
		t.Fatal("closeDatabaseWrapperAfterInitError blocked synchronously on Close() for context.Canceled")
	}

	// Unblock the background goroutine so it does not leak.
	close(closeCh)
}

// TestGetConnectionRaceCondition checks that GetConnection always returns the same instance, even when asked
// by multiple goroutines in parallel.
func TestGetConnectionRaceCondition(t *testing.T) {
	ctx := context.Background()
	b, s := getDbBackend(t)
	defer b.Cleanup(ctx)
	configureDBMount(t, s)

	goroutines := 16

	wg := sync.WaitGroup{}
	wg.Add(goroutines)
	dbis := make([]*dbPluginInstance, goroutines)
	errs := make([]error, goroutines)
	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			dbis[i], errs[i] = b.GetConnection(ctx, s, mockv5)
		}(i)
	}
	wg.Wait()
	for i := 0; i < goroutines; i++ {
		if errs[i] != nil {
			t.Fatal(errs[i])
		}
		if dbis[0] != dbis[i] {
			t.Fatal("Error: database instances did not match")
		}
	}
}

// TestGetConnectionInitializeTimeout verifies GetConnection returns an initialize-timeout
// error when plugin initialization blocks longer than the configured timeout
func TestGetConnectionInitializeTimeout(t *testing.T) {
	oldTimeout := databaseInitTimeout
	databaseInitTimeout = 25 * time.Millisecond
	defer func() {
		databaseInitTimeout = oldTimeout
	}()

	config := logical.TestBackendConfig()
	config.System = &systemViewWrapper{
		SystemView:     config.System,
		builtinFactory: newBlockingInitializeDatabase,
	}
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}
	defer b.Cleanup(context.Background())

	entry, err := logical.StorageEntryJSON("config/blocked", &DatabaseConfig{
		AllowedRoles:      []string{"*"},
		PluginName:        mockV5Type,
		VerifyConnection:  true,
		ConnectionDetails: map[string]interface{}{"connection_url": "unused"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := config.StorageView.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	_, err = b.GetConnection(context.Background(), config.StorageView, "blocked")
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !errors.Is(err, errDatabaseInitializeTimeout) {
		t.Fatalf("expected initialize timeout error, got: %v", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("GetConnection took too long to fail: %s", elapsed)
	}
	if conn := b.connections.Get("blocked"); conn != nil {
		t.Fatal("expected timed out connection to not be cached")
	}
}
