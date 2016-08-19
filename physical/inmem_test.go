package physical

import (
	"testing"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func TestInmem(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	inm := NewInmem(logger)
	testBackend(t, inm)
	testBackend_ListPrefix(t, inm)
}
