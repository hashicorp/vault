package ydb

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/physical"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/coordination"
	coordoptions "github.com/ydb-platform/ydb-go-sdk/v3/coordination/options"
)

const ydbHALockSessionTimeout = 15 * time.Second

var _ physical.Lock = (*ydbHALock)(nil)

type ydbHALock struct {
	backend   *YDBBackend
	key       string
	value     string
	semaphore string

	mu    sync.Mutex
	lease coordination.Lease
}

func (y *YDBBackend) HAEnabled() bool {
	return y.haEnabled
}

func (y *YDBBackend) LockWith(key, value string) (physical.Lock, error) {
	return &ydbHALock{
		backend:   y,
		key:       key,
		value:     value,
		semaphore: ydbHASemaphoreName(key),
	}, nil
}

func (l *ydbHALock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.lease != nil {
		return nil, fmt.Errorf("lock already held")
	}

	ctx, cancel := context.WithTimeout(context.Background(), ydbHALockSessionTimeout)
	defer cancel()

	session, err := l.backend.newCoordinationSession(
		ctx,
		coordoptions.WithDescription("vault lock "+l.key),
	)
	if err != nil {
		return nil, err
	}

	acquireCtx, acquireCancel := contextWithStopCh(stopCh)
	defer acquireCancel()

	lease, err := session.AcquireSemaphore(
		acquireCtx,
		l.semaphore,
		coordination.Exclusive,
		coordoptions.WithEphemeral(true),
		coordoptions.WithAcquireData([]byte(l.value)),
	)
	if err != nil {
		_ = closeCoordinationSession(session)
		if stopCh != nil && errors.Is(err, context.Canceled) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to acquire coordination semaphore: %w", err)
	}

	l.lease = lease

	return lease.Context().Done(), nil
}

func (l *ydbHALock) Unlock() error {
	l.mu.Lock()
	if l.lease == nil {
		l.mu.Unlock()
		return nil
	}

	lease := l.lease
	l.lease = nil
	l.mu.Unlock()

	unlockErr := lease.Release()
	if err := closeCoordinationSession(lease.Session()); unlockErr == nil {
		unlockErr = err
	}

	return unlockErr
}

func (l *ydbHALock) Value() (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ydbHALockSessionTimeout)
	defer cancel()

	session, err := l.backend.newCoordinationSession(ctx)
	if err != nil {
		return false, "", err
	}
	defer closeCoordinationSession(session)

	desc, err := session.DescribeSemaphore(ctx, l.semaphore, coordoptions.WithDescribeOwners(true))
	if err != nil {
		if ydb.IsOperationErrorSchemeError(err) {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to describe coordination semaphore: %w", err)
	}
	if len(desc.Owners) == 0 {
		return false, "", nil
	}

	return true, string(desc.Owners[0].Data), nil
}

func (y *YDBBackend) newCoordinationSession(
	ctx context.Context,
	opts ...coordoptions.SessionOption,
) (coordination.Session, error) {
	if err := y.ensureCoordinationNodeExists(ctx); err != nil {
		return nil, err
	}

	opts = append(opts, coordoptions.WithSessionTimeout(ydbHALockSessionTimeout))
	session, err := y.db.Coordination().Session(ctx, y.coordinationNode, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create coordination session: %w", err)
	}
	return session, nil
}

func contextWithStopCh(stopCh <-chan struct{}) (context.Context, context.CancelFunc) {
	if stopCh == nil {
		return context.Background(), func() {}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-stopCh:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

func closeCoordinationSession(session coordination.Session) error {
	closeCtx, cancel := context.WithTimeout(context.Background(), ydbHALockSessionTimeout)
	defer cancel()
	return session.Close(closeCtx)
}

func (y *YDBBackend) ensureCoordinationNodeExists(ctx context.Context) error {
	err := y.db.Coordination().CreateNode(ctx, y.coordinationNode, coordination.NodeConfig{
		SelfCheckPeriodMillis:    1000,
		SessionGracePeriodMillis: 1000,
		ReadConsistencyMode:      coordination.ConsistencyModeStrict,
		AttachConsistencyMode:    coordination.ConsistencyModeStrict,
		RatelimiterCountersMode:  coordination.RatelimiterCountersModeDetailed,
	})
	if err == nil || ydb.IsOperationErrorAlreadyExistsError(err) {
		return nil
	}
	return fmt.Errorf("failed to ensure coordination node exists: %w", err)
}

func ydbHASemaphoreName(key string) string {
	return "vault-lock-" + base64.RawURLEncoding.EncodeToString([]byte(key))
}
