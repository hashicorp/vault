package physical

import (
	"fmt"
	"os"
	"testing"
	"time"

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

	b, err := NewBackend("etcd", map[string]string{
		"address": addr,
		"path":    randPath,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

	ha, ok := b.(HABackend)
	if !ok {
		t.Fatalf("etcd does not implement HABackend")
	}
	testHABackend(t, ha, ha)
}
