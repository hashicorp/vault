package inmem

import (
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/physical"
)

func TestInmem(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	inm, err := NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	physical.ExerciseBackend(t, inm)
	physical.ExerciseBackend_ListPrefix(t, inm)
}
