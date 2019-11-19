// +build !enterprise

package vault

import (
	"context"
	"errors"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

// badStorage simulates a storage that returns count-1 errors on List before returning success
type badStorage struct {
	logical.Storage
	count int
}

func (b *badStorage) List(ctx context.Context, _ string) ([]string, error) {
	b.count--

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if b.count == 0 {
		return nil, nil
	}

	return nil, errors.New("storage timeout")
}

func TestExpirationManager_collectLeases(t *testing.T) {
	_, barrier, _ := mockBarrier(t)

	m := &ExpirationManager{
		quitContext: context.Background(),
		logger:      log.New(nil),
	}

	// Check that we retry a few errors
	m.idView = NewBarrierView(&badStorage{barrier, maxCollectAttempts}, "")

	_, _, err := m.collectLeases()
	if err != nil {
		t.Fatalf("error despite retries: %v", err)
	}

	// Check that we don't retry beyond the error limit
	m.idView = NewBarrierView(&badStorage{barrier, maxCollectAttempts + 1}, "")
	_, _, err = m.collectLeases()
	if err == nil {
		t.Fatal("no error despite surpassing max retries")
	}

	if err.Error() != `failed to scan for leases: list failed at path "": storage timeout` {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithCancel(m.quitContext)
	cancel()
	m.quitContext = ctx

	// Check that we only try once when there's a context error
	b := &badStorage{barrier, maxCollectAttempts}
	m.idView = NewBarrierView(b, "")
	_, _, err = m.collectLeases()
	if err == nil {
		t.Fatal("no error despite cancelled context")
	}

	if err.Error() != `failed to scan for leases: list failed at path "": context canceled` {
		t.Fatalf("unexpected error: %v", err)
	}

	if b.count != maxCollectAttempts-1 {
		t.Fatalf("unexpected number of calls to badStorage: %d", maxCollectAttempts-b.count)
	}
}
