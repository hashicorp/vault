package etcd

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
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
	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewEtcdBackend(map[string]string{
		"path": randPath,
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)

	ha, ok := b.(physical.HABackend)
	if !ok {
		t.Fatalf("etcd does not implement HABackend")
	}
	physical.ExerciseHABackend(t, ha, ha)
}
