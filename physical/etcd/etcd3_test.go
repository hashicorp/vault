package etcd

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"
)

func TestEtcd3Backend(t *testing.T) {
	addr := os.Getenv("ETCD_ADDR")
	if addr == "" {
		t.Skipf("Skipped. No etcd3 server found")
	}

	logger := logformat.NewVaultLogger(log.LevelTrace)

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
