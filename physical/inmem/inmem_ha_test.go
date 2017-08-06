package inmem

import (
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"
)

func TestInmemHA(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	inm, err := NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	physical.ExerciseHABackend(t, inm.(physical.HABackend), inm.(physical.HABackend))
}
