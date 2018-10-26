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
	config := map[string]string{
		"path":     fmt.Sprintf("/vault-%d", time.Now().Unix()),
		"etcd_api": "3",
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
