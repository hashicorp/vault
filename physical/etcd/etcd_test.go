package etcd

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"

	"go.etcd.io/etcd/client"
)

func TestEtcdBackend(t *testing.T) {
	addr := os.Getenv("ETCD_ADDR")
	if addr == "" {
		t.SkipNow()
	}

	cfg := client.Config{
		Endpoints: []string{addr},
		Transport: client.DefaultTransport,
	}

	c, err := client.New(cfg)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	syncErr := c.Sync(ctx)
	cancel()
	if syncErr != nil {
		t.Fatalf("err: %v", EtcdSyncClusterError)
	}

	kAPI := client.NewKeysAPI(c)

	randPath := fmt.Sprintf("/vault-%d", time.Now().Unix())
	defer func() {
		delOpts := &client.DeleteOptions{
			Recursive: true,
		}
		if _, err := kAPI.Delete(context.Background(), randPath, delOpts); err != nil {
			t.Fatalf("err: %v", err)
		}
	}()

	// Generate new etcd backend. The etcd address is read from ETCD_ADDR. No
	// need to provide it explicitly.
	logger := logging.NewVaultLogger(log.Debug)
	config := map[string]string{
		"path": randPath,
	}

	b, err := NewEtcdBackend(config, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	b2, err := NewEtcdBackend(config, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)
	physical.ExerciseHABackend(t, b.(physical.HABackend), b2.(physical.HABackend))
}
