package inmem

import (
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
)

func TestInmemHA(t *testing.T) {
	logger := logformat.NewVaultLogger(log.Debug)

	inm, err := NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	physical.ExerciseHABackend(t, inm.(physical.HABackend), inm.(physical.HABackend))
}
