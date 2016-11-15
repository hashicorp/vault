package physical

import (
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func TestCanoeBackend(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	b, err := NewBackend("canoe", logger, map[string]string{
		"config_port":    "1234",
		"raft_port":      "1235",
		"bootstrap_node": "true",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

	ha, ok := b.(HABackend)

	if !ok {
		t.Fatalf("canoe does not implement HABackend")
	}
	testHABackend(t, ha, ha)
}
