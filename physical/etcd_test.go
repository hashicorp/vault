package physical

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

func TestEtcdBackend(t *testing.T) {
	addr := os.Getenv("ETCD_ADDR")
	if addr == "" {
		t.SkipNow()
	}

	client := etcd.NewClient([]string{addr})
	if !client.SyncCluster() {
		t.Fatalf("err: %v", EtcdSyncClusterError)
	}

	randPath := fmt.Sprintf("/vault-%d", time.Now().Unix())
	defer func() {
		if _, err := client.Delete(randPath, true); err != nil {
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
