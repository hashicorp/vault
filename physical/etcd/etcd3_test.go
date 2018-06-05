package etcd

import (
	"fmt"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/physical"
)

func TestEtcd3Backend(t *testing.T) {
	addr := os.Getenv("ETCD_ADDR")
	if addr == "" {
		t.Skipf("Skipped. No etcd3 server found")
	}

	logger := logging.NewVaultLogger(log.Debug)

	b, err := NewEtcdBackend(map[string]string{
		"path":     fmt.Sprintf("/vault-%d", time.Now().Unix()),
		"etcd_api": "3",
	}, logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	physical.ExerciseBackend(t, b)
	physical.ExerciseBackend_ListPrefix(t, b)

	ha, ok := b.(physical.HABackend)
	if !ok {
		t.Fatalf("etcd3 does not implement HABackend")
	}
	physical.ExerciseHABackend(t, ha, ha)
}
