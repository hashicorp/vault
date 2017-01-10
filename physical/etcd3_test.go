package physical

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func TestEtcd3Backend(t *testing.T) {
	addr := os.Getenv("ETCD_ADDR")
	if addr == "" {
		t.Skipf("Skipped. No etcd3 server found")
	}

	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("etcd", logger, map[string]string{
		"path":     fmt.Sprintf("/vault-%d", time.Now().Unix()),
		"etcd_api": "3",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

	ha, ok := b.(HABackend)
	if !ok {
		t.Fatalf("etcd3 does not implement HABackend")
	}
	testHABackend(t, ha, ha)
}
